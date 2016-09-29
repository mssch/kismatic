package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector/check"
	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
	"github.com/spf13/cobra"
)

// NewCmdLocal returns the "local" command
func NewCmdLocal(out io.Writer) *cobra.Command {
	var outputType string
	var nodeRole string
	cmd := &cobra.Command{
		Use:   "local",
		Short: "run the inspector checks locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLocal(out, outputType, nodeRole)
		},
	}
	cmd.Flags().StringVarP(&outputType, "output", "o", "table", "set the result output type. Options are 'json', 'table'")
	cmd.Flags().StringVar(&nodeRole, "node-role", "", "the node's role in the cluster. Options are 'etcd', 'master', 'worker'")
	return cmd
}

func runLocal(out io.Writer, outputType, nodeRole string) error {
	if nodeRole == "" {
		return fmt.Errorf("node role is required")
	}
	if nodeRole != "etcd" && nodeRole != "master" && nodeRole != "worker" {
		return fmt.Errorf("%s is not a valid node role", nodeRole)
	}
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
	distro, err := check.DetectDistro()
	if err != nil {
		return fmt.Errorf("error running checks locally: %v", err)
	}
	_, err = check.NewPackageManager(distro)
	if err != nil {
		return err
	}
	pkgMgr, err := check.NewPackageManager(distro)
	if err != nil {
		return err
	}
	m := rule.DefaultRules()
	e := rule.Engine{
		RuleCheckMapper: rule.DefaultCheckMapper{
			PackageManager: pkgMgr,
		},
	}
	labels := []string{string(distro), nodeRole}
	results, err := e.ExecuteRules(m, labels)
	if err != nil {
		return fmt.Errorf("error running local rules: %v", err)
	}
	printResults(out, results)
	for _, r := range results {
		if !r.Success {
			return errors.New("inspector found checks that failed")
		}
	}
	return nil
}
