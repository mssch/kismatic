package cli

import (
	"errors"
	"fmt"
	"io"
	"os/exec"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/util"
	"github.com/spf13/cobra"
)

const seedRegistryLong = `
Seed a registry with the container images required by KET during the installation
or upgrade of your Kubernetes cluster.

Before using this command, you must set the location of the registry in the plan file.
Furthermore, the docker command line client must be installed on this machine to
be able to push the images to the registry.
`

type seedRegistryOptions struct {
	listOnly           bool
	verbose            bool
	outputFormat       string
	generatedAssetsDir string
	planFile           string
}

func NewCmdSeedRegistry(out io.Writer) *cobra.Command {
	var options seedRegistryOptions
	cmd := &cobra.Command{
		Use:   "seed-registry",
		Short: "seed a registry with the container images required by KET",
		Long:  seedRegistryLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return cmd.Usage()
			}
			return doSeedRegistry(out, options)
		},
	}
	cmd.Flags().BoolVar(&options.listOnly, "list-only", false, "when true, the images will only be listed instead of pushed to the registry")
	cmd.Flags().BoolVar(&options.verbose, "verbose", false, "enable verbose logging")
	cmd.Flags().StringVarP(&options.outputFormat, "output", "o", "simple", `output format (options simple|raw)`)
	cmd.Flags().StringVar(&options.generatedAssetsDir, "generated-assets-dir", "generated", "path to the directory where assets generated during the installation process will be stored")
	addPlanFileFlag(cmd.Flags(), &options.planFile)
	return cmd
}

func doSeedRegistry(out io.Writer, options seedRegistryOptions) error {
	util.PrintHeader(out, "Seed Container Image Registry", '=')
	planner := install.FilePlanner{File: options.planFile}
	if !planner.PlanExists() {
		util.PrettyPrintErr(out, "Reading installation plan file [ERROR]")
		fmt.Fprintln(out, "Run \"kismatic install plan\" to generate it")
		return fmt.Errorf("plan does not exist")
	}
	plan, err := planner.Read()
	if err != nil {
		util.PrettyPrintErr(out, "Reading installation plan file %q", options.planFile)
		return fmt.Errorf("error reading plan file: %v", err)
	}
	util.PrettyPrintOk(out, "Reading installation plan file %q", options.planFile)

	// Validate our prereqs:
	// 1. The private registry is set in the plan file
	// 2. The docker CLI is available on the local machine
	errs := []error{}
	if plan.DockerRegistry.Address == "" {
		errs = append(errs, errors.New("The private registry's address must be set in the plan file."))
	}
	if plan.DockerRegistry.Port == 0 {
		errs = append(errs, errors.New("The private registry's port must be set in the plan file."))
	}
	if plan.DockerRegistry.Port < 1 || plan.DockerRegistry.Port > 65535 {
		errs = append(errs, fmt.Errorf("The docker registry port '%d' provided in the plan file is not valid.", plan.DockerRegistry.Port))
	}
	if _, err := exec.LookPath("docker"); err != nil {
		errs = append(errs, errors.New("Did not find docker installed on this node. the docker CLI must be available for seeding the registry."))
	}

	if len(errs) > 0 {
		util.PrettyPrintErr(out, "Validating plan file")
		util.PrintValidationErrors(out, errs)
		return fmt.Errorf("Cannot seed registry due to plan file validation errors.")
	}
	util.PrettyPrintOk(out, "Validating plan file")

	execOpts := install.ExecutorOptions{
		OutputFormat:             options.outputFormat,
		Verbose:                  options.verbose,
		GeneratedAssetsDirectory: options.generatedAssetsDir,
	}
	ae, err := install.NewExecutor(out, out, execOpts)
	if err != nil {
		return err
	}
	return ae.SeedRegistry(*plan)
}
