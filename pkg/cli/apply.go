package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

// NewCmdApply creates a cluter using the plan file
func NewCmdApply(in io.Reader, out io.Writer, options *installOpts) *cobra.Command {

	cmd := &cobra.Command{
		Use:   "apply",
		Short: "apply your plan file to create a Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.planFilename}
			executor, err := install.NewAnsibleExecutor(out, os.Stderr, options.certsDestination) // TODO: Do we want to parameterize stderr?
			if err != nil {
				return err
			}
			pki := &install.LocalPKI{
				CACsr:            options.caCsr,
				CAConfigFile:     options.caConfigFile,
				CASigningProfile: options.caSigningProfile,
				DestinationDir:   options.certsDestination,
				Log:              out,
			}
			return doApply(in, out, planner, executor, pki, options)
		},
	}

	return cmd
}

func doApply(in io.Reader, out io.Writer, planner install.Planner, executor install.Executor, pki install.PKI, options *installOpts) error {
	// Check if plan file exists
	err := doValidate(in, out, planner, options)
	if err != nil {
		return fmt.Errorf("error validating plan: %v", err)
	}
	plan, err := planner.Read()
	// Generate certs for the cluster before performing installation
	fmt.Fprintln(out, "Generating cluster certificates")
	err = pki.GenerateClusterCerts(plan)
	if err != nil {
		return fmt.Errorf("error generating certificates for the cluster: %v", err)
	}
	fmt.Fprintf(out, "Generated cluster certificates at %q [OK]\n\n", options.certsDestination)

	// Execute playbooks
	err = executor.Install(plan)
	if err != nil {
		return fmt.Errorf("error installing: %v", err)
	}

	// Generate kubeconfig
	err = install.GenerateKubeconfig(plan, options.certsDestination)
	if err != nil {
		fmt.Fprint(out, "Kubeconfig generation error, you may need to setup kubectl manually [ERROR]\n", err)
	} else {
		fmt.Fprint(out, "Generated \"config\", to use \"cp config ~/.kube/config\" [OK]")
	}

	return nil
}
