package cli

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

// NewCmdValidate creates a new install validate command
func NewCmdValidate(out io.Writer, options *install.CliOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate your plan file",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.PlanFilename}
			return doValidate(out, planner, options)
		},
	}

	return cmd
}

func doValidate(out io.Writer, planner install.Planner, options *install.CliOpts) error {
	// Check if plan file exists
	if !planner.PlanExists() {
		fmt.Fprintf(out, "Reading installation plan file [ERROR]\n")
		fmt.Fprintf(out, "Run \"kismatic install plan\" to generate it\n")
		return fmt.Errorf("plan does not exist")
	}
	plan, err := planner.Read()
	if err != nil {
		fmt.Fprintf(out, "Reading installation plan file %q [ERROR]\n", options.PlanFilename)
		return fmt.Errorf("error reading plan file: %v", err)
	}
	fmt.Fprintf(out, "Reading installation plan file %q [OK]\n", options.PlanFilename)

	// Verify plan file
	ok, errs := install.ValidatePlan(plan)
	if !ok {
		fmt.Fprint(out, "Validating installation plan file [ERROR]\n")
		for _, err := range errs {
			fmt.Fprintf(out, "- %v\n", err)
		}
		return fmt.Errorf("validation error prevents installation from proceeding")
	}
	fmt.Fprint(out, "Validating installation plan file [OK]\n")

	return nil
}
