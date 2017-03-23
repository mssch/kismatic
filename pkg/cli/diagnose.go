package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/util"
	"github.com/spf13/cobra"
)

type diagsOpts struct {
	planFilename string
	verbose      bool
	outputFormat string
}

// NewCmdDiagnostic collects diagnostic data on remote nodes
func NewCmdDiagnostic(out io.Writer) *cobra.Command {
	opts := &diagsOpts{}

	cmd := &cobra.Command{
		Use:   "diagnose",
		Short: "Collects diagnostics about the nodes in the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}

			return doDiagnostics(out, opts)
		},
	}

	// PersistentFlags
	addPlanFileFlag(cmd.PersistentFlags(), &opts.planFilename)
	cmd.Flags().BoolVar(&opts.verbose, "verbose", false, "enable verbose logging from the installation")
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "simple", "installation output format (options \"simple\"|\"raw\")")

	return cmd
}

func doDiagnostics(out io.Writer, opts *diagsOpts) error {
	planFile := opts.planFilename
	planner := install.FilePlanner{File: planFile}

	// Read plan file
	if !planner.PlanExists() {
		util.PrettyPrintErr(out, "Reading plan file")
		return fmt.Errorf("plan file %q does not exist", planFile)
	}
	util.PrettyPrintOk(out, "Reading plan file")
	plan, err := planner.Read()
	if err != nil {
		util.PrettyPrintErr(out, "Reading plan file")
		return fmt.Errorf("error reading plan file %q: %v", planFile, err)
	}

	// Validate SSH connectivity to nodes
	if ok, errs := install.ValidatePlanSSHConnections(plan); !ok {
		util.PrettyPrintErr(out, "Validate SSH connectivity to nodes")
		util.PrintValidationErrors(out, errs)
		return fmt.Errorf("SSH connectivity validation errors found")
	}
	util.PrettyPrintOk(out, "Validate SSH connectivity to nodes")

	// Get versions as only supported nodes are >=1.3
	cv, err := install.ListVersions(plan)
	if err != nil {
		return fmt.Errorf("error listing cluster versions: %v", err)
	}
	var toDiagnose []install.ListableNode
	var toSkip []install.ListableNode
	for _, n := range cv.Nodes {
		if install.IsGreaterOrEqualThanVersion(n.Version, "v1.3.0-alpha") {
			toDiagnose = append(toDiagnose, n)
		} else {
			toSkip = append(toSkip, n)
		}
	}

	// Print the nodes that will be skipped
	if len(toSkip) > 0 {
		util.PrintHeader(out, "Skipping nodes that are not eligible", '=')
		for _, n := range toSkip {
			util.PrettyPrintOk(out, "- %q is at an unsupported version %q", n.Node.Host, n.Version)
		}
		fmt.Fprintln(out)
	}

	// Print message if there's no work to do
	if len(toDiagnose) == 0 {
		fmt.Fprintln(out, "All nodes have an unsupported version")
	} else {
		options := install.ExecutorOptions{
			OutputFormat: opts.outputFormat,
			Verbose:      opts.verbose,
		}
		executor, err := install.NewDiagnosticsExecutor(out, os.Stderr, options)
		if err != nil {
			return err
		}
		return executor.DiagnoseNodes(*plan)
	}

	return nil

}
