package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
)

func validateOutputType(outputType string) error {
	if outputType != "json" && outputType != "table" {
		return fmt.Errorf("output type %q not supported", outputType)
	}
	return nil
}

func printResults(out io.Writer, results []rule.RuleResult, outputType string) error {
	switch outputType {
	case "json":
		return printResultsAsJSON(out, results)
	case "table":
		return printResultsAsTable(out, results)
	default:
		return fmt.Errorf("output type %q not supported", outputType)
	}
}

func printResultsAsJSON(out io.Writer, results []rule.RuleResult) error {
	err := json.NewEncoder(out).Encode(results)
	if err != nil {
		return fmt.Errorf("error marshaling results as JSON: %v", err)
	}
	return nil
}

func printResultsAsTable(out io.Writer, results []rule.RuleResult) error {
	w := tabwriter.NewWriter(out, 1, 8, 4, '\t', 0)
	fmt.Fprintf(w, "CHECK\tSUCCESS\tMSG\n")
	for _, r := range results {
		fmt.Fprintf(w, "%s\t%t\t%v\n", r.Name, r.Success, r.Error)
	}
	w.Flush()
	return nil
}
