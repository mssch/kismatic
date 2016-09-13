package cli

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdValidate creates a new install validate command
func NewCmdValidate(out io.Writer, options *installOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "validate your plan file",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.planFilename}
			return doValidate(out, planner, options.planFilename)
		},
	}

	return cmd
}

func doValidate(out io.Writer, planner install.Planner, planFile string) error {
	util.PrintHeader(out, "Validating")
	// Check if plan file exists
	if !planner.PlanExists() {
		util.PrettyPrintErr(out, "Reading installation plan file [ERROR]")
		util.PrettyPrint(out, "Run \"kismatic install plan\" to generate it\n")
		return fmt.Errorf("plan does not exist")
	}
	plan, err := planner.Read()
	if err != nil {
		util.PrettyPrintErrf(out, "Reading installation plan file %q", planFile)
		return fmt.Errorf("error reading plan file: %v", err)
	}
	util.PrettyPrintOkf(out, "Reading installation plan file %q", planFile)

	// Verify plan file
	ok, errs := install.ValidatePlan(plan)
	if !ok {
		util.PrettyPrintErr(out, "Validating installation plan file")
		for _, err := range errs {
			util.PrintErrorf(out, "- %v\n", err)
		}
		return fmt.Errorf("validation error prevents installation from proceeding")
	}
	util.PrettyPrintOk(out, "Validating installation plan file")

	return nil
}
