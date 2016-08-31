package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

type installOpts struct {
	planFilename     string
	caCsr            string
	caConfigFile     string
	caSigningProfile string
	certsDestination string
	dryRun           bool
}

// NewCmdInstall creates a new install command
func NewCmdInstall(in io.Reader, out io.Writer) *cobra.Command {
	options := &installOpts{}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install your Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.planFilename}
			pki := &install.LocalPKI{
				CACsr:            options.caCsr,
				CAConfigFile:     options.caConfigFile,
				CASigningProfile: options.caSigningProfile,
				DestinationDir:   options.certsDestination,
				Log:              out,
			}
			executor, err := install.NewAnsibleExecutor(out, os.Stderr, options.certsDestination) // TODO: Do we want to parameterize stderr?
			if err != nil {
				return err
			}
			return doInstall(in, out, planner, executor, pki, options)
		},
	}

	cmd.Flags().StringVarP(&options.planFilename, "plan-file", "f", "kismatic-cluster.yaml", "path to the installation plan file")
	cmd.Flags().StringVar(&options.caCsr, "ca-csr", "ansible/playbooks/tls/ca-csr.json", "path to the Certificate Authority CSR")
	cmd.Flags().StringVar(&options.caConfigFile, "ca-config", "ansible/playbooks/tls/ca-config.json", "path to the Certificate Authority configuration file")
	cmd.Flags().StringVar(&options.caSigningProfile, "ca-signing-profile", "kubernetes", "name of the profile to be used for signing certificates")
	cmd.Flags().StringVar(&options.certsDestination, "generated-certs-dir", "generated-certs", "path to the directory where generated cluster certificates will be stored")
	cmd.Flags().BoolVar(&options.dryRun, "dry-run", false, "run planning and validation phases, but don't perform installation")

	return cmd
}

func doInstall(in io.Reader, out io.Writer, planner install.Planner, executor install.Executor, pki install.PKI, options *installOpts) error {
	if !planner.PlanExists() {
		// Plan file not found, planning phase
		fmt.Fprintln(out, "Plan your Kismatic cluster:")

		etcdNodes, err := promptForInt(in, out, "Number of etcd nodes", 3)
		if err != nil {
			return fmt.Errorf("Error reading number of etcd nodes: %v", err)
		}
		if etcdNodes <= 0 {
			return fmt.Errorf("The number of etcd nodes must be greater than zero")
		}

		masterNodes, err := promptForInt(in, out, "Number of master nodes", 2)
		if err != nil {
			return fmt.Errorf("Error reading number of master nodes: %v", err)
		}
		if masterNodes <= 0 {
			return fmt.Errorf("The number of master nodes must be greater than zero")
		}

		workerNodes, err := promptForInt(in, out, "Number of worker nodes", 3)
		if err != nil {
			return fmt.Errorf("Error reading number of worker nodes: %v", err)
		}
		if workerNodes <= 0 {
			return fmt.Errorf("The number of worker nodes must be greater than zero")
		}

		fmt.Fprintf(out, "Generating installation plan file with %d etcd nodes, %d master nodes and %d worker nodes\n",
			etcdNodes, masterNodes, workerNodes)

		m := install.MasterNodeGroup{}
		m.ExpectedCount = masterNodes

		p := install.Plan{
			Etcd: install.NodeGroup{
				ExpectedCount: etcdNodes,
			},
			Master: m,
			Worker: install.NodeGroup{
				ExpectedCount: workerNodes,
			},
		}
		err = install.WritePlanTemplate(p, planner)
		if err != nil {
			return fmt.Errorf("error planning installation: %v", err)
		}
		fmt.Fprintf(out, "Generated installation plan file at %q\n", options.planFilename)
		fmt.Fprintf(out, "Edit the file to further describe your cluster. Once ready, execute the install command to proceed.\n")
		return nil
	}

	// Plan file exists, validate and install
	p, err := planner.Read()
	if err != nil {
		fmt.Fprintf(out, "Reading installation plan file %q [ERROR]\n", options.planFilename)
		return fmt.Errorf("error reading plan file: %v", err)
	}
	fmt.Fprintf(out, "Reading installation plan file %q [OK]\n", options.planFilename)
	fmt.Fprintln(out, "")

	ok, errs := install.ValidatePlan(p)
	if !ok {
		fmt.Fprint(out, "Validating installation plan file [ERROR]\n")
		for _, err := range errs {
			fmt.Fprintf(out, "- %v\n", err)
		}
		fmt.Fprintln(out, "")
		return fmt.Errorf("validation error prevents installation from proceeding")

	}
	fmt.Fprint(out, "Validating installation plan file [OK]\n\n")

	// Generate certs for the cluster before performing installation
	fmt.Fprintln(out, "Generating cluster certificates")
	err = pki.GenerateClusterCerts(p)
	if err != nil {
		return fmt.Errorf("error generating certificates for the cluster: %v", err)
	}
	fmt.Fprintf(out, "Generated cluster certificates at %q [OK]\n\n", options.certsDestination)

	if options.dryRun {
		return nil
	}

	err = executor.Install(p)
	if err != nil {
		return fmt.Errorf("error installing: %v", err)
	}

	// Generate kubeconfig
	fmt.Fprint(out, "Generating kubecofnig\n")
	fmt.Fprint(out, "===========================================================================\n")
	err = install.GenerateKubeconfig(p, "admin", options.certsDestination, out)
	if err != nil {
		fmt.Fprint(out, "Kubeconfig generation error, you may need to setup kubectl manually [ERROR]\n")
	} else {
		fmt.Fprint(out, "===========================================================================\n")
		fmt.Fprint(out, "Generated kubecofnig, to use run \"kubectl config --kubeconfig $PATH_TO_YOUR_KUBECONFIG\" with the content above [OK]\n")
	}

	return nil
}

func promptForInt(in io.Reader, out io.Writer, prompt string, defaultValue int) (int, error) {
	fmt.Fprintf(out, "=> %s [%d]: ", prompt, defaultValue)
	s := bufio.NewScanner(in)
	// Scan the first token
	s.Scan()
	if s.Err() != nil {
		return defaultValue, fmt.Errorf("error reading number: %v", s.Err())
	}
	ans := s.Text()
	if ans == "" {
		return defaultValue, nil
	}
	// Convert input into integer
	i, err := strconv.Atoi(ans)
	if err != nil {
		return defaultValue, fmt.Errorf("%q is not a number", ans)
	}
	return i, nil
}
