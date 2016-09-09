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

// NewCmdApply creates a cluter using the plan file
func NewCmdApply(out io.Writer, options *install.CliOpts) *cobra.Command {
	skipCAGeneration := false
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply your plan file to create a Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.PlanFilename}
			// TODO: Do we want to parameterize stderr?
			executor, err := install.NewExecutor(out, os.Stderr, options.CertsDestination, options.RestartServices)
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
			applyCmd := &applyCmd{
				out,
				planner,
				executor,
				pki,
				options.PlanFilename,
				skipCAGeneration,
				options.CertsDestination,
			}
			return applyCmd.run()
		},
	}

	// Flags
	cmd.Flags().StringVar(&options.CaCsr, "ca-csr", "ansible/playbooks/tls/ca-csr.json", "path to the Certificate Authority CSR")
	cmd.Flags().StringVar(&options.CaConfigFile, "ca-config", "ansible/playbooks/tls/ca-config.json", "path to the Certificate Authority configuration file")
	cmd.Flags().StringVar(&options.CaSigningProfile, "ca-signing-profile", "kubernetes", "name of the profile to be used for signing certificates")
	cmd.Flags().StringVar(&options.CertsDestination, "generated-certs-dir", "generated-certs", "path to the directory where generated cluster certificates will be stored")
	cmd.Flags().BoolVar(&skipCAGeneration, "skip-ca-generation", false, "skip CA generation and use an existing file")
	cmd.Flags().BoolVar(&options.RestartServices, "restart-services", false, "force restart clusters services (Use with care)")

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
