package cmd

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector"
	"github.com/spf13/cobra"
)

type clientOpts struct {
	outputType string
	nodeRoles  string
	rulesFile  string
	targetNode string
}

// NewCmdClient returns the "client" command
func NewCmdClient(out io.Writer) *cobra.Command {
	opts := clientOpts{}
	cmd := &cobra.Command{
		Use:   "client TARGET_NODE_IP",
		Short: "run the inspector against a remote inspector server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Usage()
			}
			opts.targetNode = args[0]
			return runClient(out, opts)
		},
	}
	cmd.Flags().StringVarP(&opts.outputType, "output", "o", "table", "set the result output type. Options are 'json', 'table'")
	cmd.Flags().StringVar(&opts.nodeRoles, "node-roles", "", "comma-separated list of the node's roles. Valid roles are 'etcd', 'master', 'worker'")
	cmd.Flags().StringVarP(&opts.rulesFile, "file", "f", "", "the path to an inspector rules file. If blank, the inspector uses the default rules")
	cmd.Flags().StringVar(&opts.targetNode, "target", "", "the node ip:port that is running the inspector in server mode")
	return cmd
}

func runClient(out io.Writer, opts clientOpts) error {
	if err := validateOutputType(opts.outputType); err != nil {
		return err
	}
	if opts.nodeRoles == "" {
		return fmt.Errorf("--node-roles is required")
	}
	roles, err := getNodeRoles(opts.nodeRoles)
	if err != nil {
		return err
	}
	c, err := inspector.NewClient(opts.targetNode, roles)
	if err != nil {
		return fmt.Errorf("error creating inspector client: %v", err)
	}
	rules, err := getRulesFromFileOrDefault(out, opts.rulesFile)
	if err != nil {
		return err
	}
	results, err := c.ExecuteRules(rules)
	if err != nil {
		return fmt.Errorf("error running inspector against remote node: %v", err)
	}
	if err := printResults(out, results, opts.outputType); err != nil {
		return err
	}
	return nil
}
