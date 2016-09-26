package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/apprenda/kismatic-platform/pkg/inspector"
)

type resultPrinter func(out io.Writer, r []inspector.CheckResult) error

func printResultsAsJSON(out io.Writer, r []inspector.CheckResult) error {
	err := json.NewEncoder(out).Encode(r)
	if err != nil {
		return fmt.Errorf("error marshaling results as JSON: %v", err)
	}
	return nil
}

func printResultsAsTable(out io.Writer, r []inspector.CheckResult) error {
	w := tabwriter.NewWriter(out, 1, 8, 4, '\t', 0)
	fmt.Fprintf(w, "CHECK\tSUCCESS\tMSG\n")
	for _, cr := range r {
		fmt.Fprintf(w, "%s\t%t\t%v\n", cr.Name, cr.Success, cr.Error)
	}
	w.Flush()
	return nil
}
