package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/spf13/cobra"
)

type applyCmd struct {
	out              io.Writer
	planner          install.Planner
	executor         install.Executor
	pki              install.PKI
	planFile         string
	skipCAGeneration bool
	certsDestination string
}

type applyOpts struct {
	caCSR            string
	caConfigFile     string
	caSigningProfile string
	certsDestination string
	skipCAGeneration bool
	restartServices  bool
	modifyHostsFile  bool
	verbose          bool
	outputFormat     string
}

// NewCmdApply creates a cluter using the plan file
func NewCmdApply(out io.Writer, installOpts *installOpts) *cobra.Command {
	applyOpts := applyOpts{}
	skipCAGeneration := false
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply your plan file to create a Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: installOpts.planFilename}
			// TODO: Do we want to parameterize stderr?
			executor, err := install.NewExecutor(out, os.Stderr, applyOpts.certsDestination, applyOpts.restartServices, applyOpts.modifyHostsFile, applyOpts.verbose, applyOpts.outputFormat)
			if err != nil {
				return err
			}
			pki := &install.LocalPKI{
				CACsr:            applyOpts.caCSR,
				CAConfigFile:     applyOpts.caConfigFile,
				CASigningProfile: applyOpts.caSigningProfile,
				DestinationDir:   applyOpts.certsDestination,
				Log:              out,
			}
			applyCmd := &applyCmd{
				out,
				planner,
				executor,
				pki,
				installOpts.planFilename,
				skipCAGeneration,
				applyOpts.certsDestination,
			}
			return applyCmd.run()
		},
	}

	// Flags
	cmd.Flags().StringVar(&applyOpts.caCSR, "ca-csr", "ansible/playbooks/tls/ca-csr.json", "path to the Certificate Authority CSR")
	cmd.Flags().StringVar(&applyOpts.caConfigFile, "ca-config", "ansible/playbooks/tls/ca-config.json", "path to the Certificate Authority configuration file")
	cmd.Flags().StringVar(&applyOpts.caSigningProfile, "ca-signing-profile", "kubernetes", "name of the profile to be used for signing certificates")
	cmd.Flags().StringVar(&applyOpts.certsDestination, "generated-certs-dir", "generated-certs", "path to the directory where generated cluster certificates will be stored")
	cmd.Flags().BoolVar(&applyOpts.skipCAGeneration, "skip-ca-generation", false, "skip CA generation and use an existing file")
	cmd.Flags().BoolVar(&applyOpts.restartServices, "restart-services", false, "force restart clusters services (Use with care)")
	cmd.Flags().BoolVar(&applyOpts.modifyHostsFile, "modify-hosts-file", false, "modify hosts files on all target nodes, only required if DNS is not available")
	cmd.Flags().BoolVar(&applyOpts.verbose, "verbose", false, "enable verbose logging from the installation")
	cmd.Flags().StringVarP(&applyOpts.outputFormat, "output", "o", "simple", "installation output format. Supported options: simple|raw")

	return cmd
}

func (c *applyCmd) run() error {
	// Check if plan file exists
	err := doValidate(c.out, c.planner, c.planFile)
	if err != nil {
		return fmt.Errorf("error validating plan: %v", err)
	}
	plan, err := c.planner.Read()

	// Generate or read cluster Certificate Authority
	var ca *tls.CA
	if !c.skipCAGeneration {
		fmt.Fprintln(c.out, "Generating cluster Certificate Authority")
		ca, err = c.pki.GenerateClusterCA(plan)
		if err != nil {
			return fmt.Errorf("error generating CA for the cluster: %v", err)
		}
	} else {
		fmt.Fprintln(c.out, "Skipping Certificate Authority generation.")
		ca, err = c.pki.ReadClusterCA(plan)
		if err != nil {
			return fmt.Errorf("error reading cluster CA: %v", err)
		}
	}

	// Generate node and user certificates
	fmt.Fprintln(c.out, "Generating cluster certificates")
	err = c.pki.GenerateClusterCerts(plan, ca, []string{"admin"})
	if err != nil {
		return fmt.Errorf("error generating certificates for the cluster: %v", err)
	}
	fmt.Fprintf(c.out, "Generated cluster certificates at %q [OK]\n\n", c.certsDestination)

	// Perform the installation
	err = c.executor.Install(plan)
	if err != nil {
		return fmt.Errorf("error installing: %v", err)
	}

	// Generate kubeconfig
	err = install.GenerateKubeconfig(plan, c.certsDestination)
	if err != nil {
		fmt.Fprint(c.out, "Kubeconfig generation error, you may need to setup kubectl manually [ERROR]\n", err)
	} else {
		fmt.Fprint(c.out, "Generated \"config\", to use \"cp config ~/.kube/config\" [OK]")
	}

	return nil
}
