package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/spf13/cobra"
)

type applyCmd struct {
	out                io.Writer
	planner            install.Planner
	executor           install.Executor
	planFile           string
	skipCAGeneration   bool
	generatedAssetsDir string
}

type applyOpts struct {
	caCSR              string
	caConfigFile       string
	caSigningProfile   string
	generatedAssetsDir string
	skipCAGeneration   bool
	restartServices    bool
	verbose            bool
	outputFormat       string
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
			executorOpts := install.ExecutorOptions{
				CAConfigFile:             applyOpts.caConfigFile,
				CASigningRequest:         applyOpts.caCSR,
				CASigningProfile:         applyOpts.caSigningProfile,
				SkipCAGeneration:         applyOpts.skipCAGeneration,
				GeneratedAssetsDirectory: applyOpts.generatedAssetsDir,
				RestartServices:          applyOpts.restartServices,
				OutputFormat:             applyOpts.outputFormat,
				Verbose:                  applyOpts.verbose,
			}
			// TODO: Do we want to parameterize stderr?
			executor, err := install.NewExecutor(out, os.Stderr, executorOpts)
			if err != nil {
				return err
			}

			applyCmd := &applyCmd{
				out,
				planner,
				executor,
				installOpts.planFilename,
				skipCAGeneration,
				applyOpts.generatedAssetsDir,
			}
			return applyCmd.run()
		},
	}

	// Flags
	cmd.Flags().StringVar(&applyOpts.caCSR, "ca-csr", "ansible/playbooks/tls/ca-csr.json", "path to the Certificate Authority CSR")
	cmd.Flags().StringVar(&applyOpts.caConfigFile, "ca-config", "ansible/playbooks/tls/ca-config.json", "path to the Certificate Authority configuration file")
	cmd.Flags().StringVar(&applyOpts.caSigningProfile, "ca-signing-profile", "kubernetes", "name of the profile to be used for signing certificates")
	cmd.Flags().StringVar(&applyOpts.generatedAssetsDir, "generated-assets-dir", "generated", "path to the directory where assets generated during the installation process are to be stored")
	cmd.Flags().BoolVar(&applyOpts.skipCAGeneration, "skip-ca-generation", false, "skip CA generation and use an existing file")
	cmd.Flags().BoolVar(&applyOpts.restartServices, "restart-services", false, "force restart clusters services (Use with care)")
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

	// Run pre-flight check
	err = c.executor.RunPreflightCheck(plan)
	if err != nil {
		return fmt.Errorf("error during pre-flight checks: %v", err)
	}

	// Perform the installation
	err = c.executor.Install(plan)
	if err != nil {
		return fmt.Errorf("error installing: %v", err)
	}
	util.PrintColor(c.out, util.Green, "\nThe cluster was installed successfully\n")

	// Generate kubeconfig
	util.PrintHeader(c.out, "Generating Kubeconfig File", '=')
	err = install.GenerateKubeconfig(plan, c.generatedAssetsDir)
	if err != nil {
		util.PrettyPrintWarn(c.out, "Error generating kubeconfig file: %v\n", err)
	} else {
		util.PrettyPrintOk(c.out, "Generated kubeconfig file in the %q directory", c.generatedAssetsDir)
		fmt.Fprintf(c.out, "\n")
		msg := "To use the generated kubeconfig file with kubectl:" +
			"\n  * use \"kubectl --kubeconfig %s/kubeconfig\"" +
			"\n  * or copy the config file \"cp %[1]s/kubeconfig ~/.kube/config\"\n"
		fmt.Fprintf(c.out, msg, c.generatedAssetsDir)
	}

	fmt.Fprintf(c.out, "\n")
	return nil
}
