package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/spf13/cobra"
)

// NewCmdApply creates a cluter using the plan file
func NewCmdApply(out io.Writer, options *install.CliOpts) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply your plan file to create a Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.PlanFilename}
			executor, err := install.NewAnsibleExecutor(out, os.Stderr, options.CertsDestination) // TODO: Do we want to parameterize stderr?
			if err != nil {
				return err
			}
			pki := &install.LocalPKI{
				CACsr:            options.CaCsr,
				CAConfigFile:     options.CaConfigFile,
				CASigningProfile: options.CaSigningProfile,
				DestinationDir:   options.CertsDestination,
				Log:              out,
			}
			return doApply(out, planner, executor, pki, options)
		},
	}

	// Flags
	cmd.Flags().StringVar(&options.CaCsr, "ca-csr", "ansible/playbooks/tls/ca-csr.json", "path to the Certificate Authority CSR")
	cmd.Flags().StringVar(&options.CaConfigFile, "ca-config", "ansible/playbooks/tls/ca-config.json", "path to the Certificate Authority configuration file")
	cmd.Flags().StringVar(&options.CaSigningProfile, "ca-signing-profile", "kubernetes", "name of the profile to be used for signing certificates")
	cmd.Flags().StringVar(&options.CertsDestination, "generated-certs-dir", "generated-certs", "path to the directory where generated cluster certificates will be stored")
	cmd.Flags().BoolVar(&options.SkipCAGeneration, "skip-ca-generation", false, "skip CA generation and use an existing file")
	cmd.Flags().BoolVar(&options.RestartEtcdService, "restart-etcd-service", false, "restart etcd service")
	cmd.Flags().BoolVar(&options.RestartKubernetesService, "restart-kubernetes-service", false, "restart kubernetes service")
	cmd.Flags().BoolVar(&options.RestartCalicoService, "restart-calico-service", false, "restart calico service")
	cmd.Flags().BoolVar(&options.RestartDockerService, "restart-docker-service", false, "restart docker service")

	return cmd
}

func doApply(out io.Writer, planner install.Planner, executor install.Executor, pki install.PKI, options *install.CliOpts) error {
	// Check if plan file exists
	err := doValidate(out, planner, options)
	if err != nil {
		return fmt.Errorf("error validating plan: %v", err)
	}
	plan, err := planner.Read()
	// Generate certs for the cluster before performing installation
	var ca *tls.CA
	if !options.SkipCAGeneration {
		fmt.Fprintln(out, "Generating cluster CA")
		ca, err = pki.GenerateClusterCA(plan)
	} else {
		fmt.Fprintln(out, "Reading cluster CA")
		ca, err = pki.ReadClusterCA(plan)
	}
	if err != nil {
		return fmt.Errorf("error generating CA for the cluster: %v", err)
	}
	fmt.Fprintln(out, "Generating cluster certificates")
	err = pki.GenerateClusterCerts(plan, ca, []string{"admin"})
	if err != nil {
		return fmt.Errorf("error generating certificates for the cluster: %v", err)
	}

	fmt.Fprintf(out, "Generated cluster certificates at %q [OK]\n\n", options.CertsDestination)

	// Execute playbooks
	av, err := executor.GetVars(plan, options)
	if err != nil {
		return fmt.Errorf("error creating variables: %v", err)
	}
	err = executor.Install(plan, av)
	if err != nil {
		return fmt.Errorf("error installing: %v", err)
	}

	// Generate kubeconfig
	err = install.GenerateKubeconfig(plan, options.CertsDestination)
	if err != nil {
		fmt.Fprint(out, "Kubeconfig generation error, you may need to setup kubectl manually [ERROR]\n", err)
	} else {
		fmt.Fprint(out, "Generated \"config\", to use \"cp config ~/.kube/config\" [OK]")
	}

	return nil
}
