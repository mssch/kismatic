package cmd

import (
	"errors"
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector"
	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
	"github.com/spf13/cobra"
)

type clientOpts struct {
	outputType string
	nodeRole   string
	rulesFile  string
	targetNode string
}

// NewCmdClient returns the "client" command
func NewCmdClient(out io.Writer) *cobra.Command {
	opts := clientOpts{}
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run the inspector against a remote inspector server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if opts.targetNode == "" {
				return errors.New("target node (--target) cannot be empty")
			}
			if err := validateOutputType(opts.outputType); err != nil {
				return err
			}
			c, err := inspector.NewClient(opts.targetNode, opts.nodeRole)
			if err != nil {
				return fmt.Errorf("error creating inspector client: %v", err)
			}
			var rules []rule.Rule
			if opts.rulesFile != "" {
				rules, err = rule.ReadFromFile(opts.rulesFile)
				if err != nil {
					return fmt.Errorf("error reading rules from %q: %v", opts.rulesFile, err)
				}
				if ok := validateRules(out, rules); !ok {
					return fmt.Errorf("rules read from %q did not pass validation", opts.rulesFile)
				}
			} else {
				rules = rule.DefaultRules()
			}
			results, err := c.ExecuteRules(rules)
			if err != nil {
				return fmt.Errorf("error running inspector against remote node: %v", err)
			}
			if err := printResults(out, results, opts.outputType); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVarP(&opts.outputType, "output", "o", "table", "set the result output type. Options are 'json', 'table'")
	cmd.Flags().StringVar(&opts.nodeRole, "node-role", "", "the node's role in the cluster. Options are 'etcd', 'master', 'worker'")
	cmd.Flags().StringVarP(&opts.rulesFile, "file", "f", "", "the path to an inspector rules file. If blank, the inspector uses the default rules")
	cmd.Flags().StringVar(&opts.targetNode, "target", "", "the node ip:port that is running the inspector in server mode")
	return cmd
}
