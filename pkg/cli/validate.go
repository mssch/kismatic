package cli

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

// NewCmdValidate creates a new install validate command
func NewCmdValidate(in io.Reader, out io.Writer, options *installOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate your plan file",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.planFilename}
			return doValidate(in, out, planner, options)
		},
	}

	return cmd
}

func doValidate(in io.Reader, out io.Writer, planner install.Planner, options *installOpts) error {
	// Check if plan file exists
	plan, err := planner.Read()
	if err != nil {
		fmt.Fprintf(out, "Reading installation plan file %q [ERROR]\n", options.planFilename)
		return fmt.Errorf("error reading plan file: %v", err)
	}
	fmt.Fprintf(out, "Reading installation plan file %q [OK]\n", options.planFilename)

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
