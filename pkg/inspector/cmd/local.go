package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector"
	"github.com/apprenda/kismatic-platform/pkg/inspector/check"
	"github.com/spf13/cobra"
)

func NewCmdLocal(out io.Writer) *cobra.Command {
	var outputType string
	cmd := &cobra.Command{
		Use:   "local",
		Short: "run the inspector checks locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLocal(out, outputType)
		},
	}
	cmd.Flags().StringVarP(&outputType, "output", "o", "table", "set the result output type. Options are 'json', 'table'")
	return cmd
}

func runLocal(out io.Writer, outputType string) error {
	// Get the printer
	var printResults resultPrinter
	switch outputType {
	case "json":
		printResults = printResultsAsJSON
	case "table":
		printResults = printResultsAsTable
	default:
		return fmt.Errorf("output type %q not supported", outputType)
	}
	d, err := check.DetectDistro()
	if err != nil {
		return fmt.Errorf("error running checks locally: %v", err)
	}
	_, err = check.NewPackageManager(d)
	if err != nil {
		return err
	}
	m := inspector.DefaultRules()
	e := inspector.Engine{}
	labels := []string{"centos", "worker"}
	results, err := e.ExecuteRules(m, labels)
	if err != nil {
		return fmt.Errorf("Error running local rules: %v", err)
	}
	printResults(out, results)
	for _, r := range results {
		if !r.Success {
			return errors.New("inspector found checks that failed")
		}
	}
	return nil
}
