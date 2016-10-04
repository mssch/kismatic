package cmd

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector/check"
	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
	"github.com/spf13/cobra"
)

type localOpts struct {
	outputType string
	nodeRole   string
	rulesFile  string
}

// NewCmdLocal returns the "local" command
func NewCmdLocal(out io.Writer) *cobra.Command {
	opts := localOpts{}
	cmd := &cobra.Command{
		Use:   "local",
		Short: "run the inspector checks locally",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLocal(out, opts)
		},
	}
	cmd.Flags().StringVarP(&opts.outputType, "output", "o", "table", "set the result output type. Options are 'json', 'table'")
	cmd.Flags().StringVar(&opts.nodeRole, "node-role", "", "the node's role in the cluster. Options are 'etcd', 'master', 'worker'")
	cmd.Flags().StringVarP(&opts.rulesFile, "file", "f", "", "the path to an inspector rules file. If blank, the inspector uses the default rules")
	return cmd
}

func runLocal(out io.Writer, opts localOpts) error {
	nodeRole := opts.nodeRole
	if nodeRole == "" {
		return fmt.Errorf("node role is required")
	}
	if nodeRole != "etcd" && nodeRole != "master" && nodeRole != "worker" {
		return fmt.Errorf("%s is not a valid node role", nodeRole)
	}
	if err := validateOutputType(opts.outputType); err != nil {
		return err
	}
	// Gather rules
	var rules []rule.Rule
	var err error
	if opts.rulesFile != "" {
		rules, err = rule.ReadFromFile(opts.rulesFile)
		if err != nil {
			return err
		}
		if ok := validateRules(out, rules); !ok {
			return fmt.Errorf("rules read from %q did not pass validation", opts.rulesFile)
		}
	} else {
		rules = rule.DefaultRules()
	}
	// Set up engine dependencies
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
	// Create rule engine
	e := rule.Engine{
		RuleCheckMapper: rule.DefaultCheckMapper{
			PackageManager: pkgMgr,
		},
	}
	labels := []string{string(distro), nodeRole}
	results, err := e.ExecuteRules(rules, labels)
	if err != nil {
		return fmt.Errorf("error running local rules: %v", err)
	}
	if err := printResults(out, results, opts.outputType); err != nil {
		return fmt.Errorf("error printing results: %v", err)
	}
	return nil
}
