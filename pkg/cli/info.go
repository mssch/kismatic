package cli

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/util"
	"github.com/spf13/cobra"
)

type infoOpts struct {
	planFilename string
	outputFormat string
}

// NewCmdInfo returns the info command
func NewCmdInfo(out io.Writer) *cobra.Command {
	opts := &infoOpts{}
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Display info about nodes in the cluster",
		Long: `will list the external IP addresses of nodes that make up the cluster, along with their current versions & roles.

This will retreived by connecting to each node via ssh`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return list(out, opts)
		},
	}
	cmd.Flags().StringVarP(&opts.planFilename, "plan-file", "f", "kismatic-cluster.yaml", "path to the installation plan file")
	cmd.Flags().StringVarP(&opts.outputFormat, "output", "o", "simple", `output format (options "simple"|"json")`)
	return cmd
}

func list(out io.Writer, opts *infoOpts) error {
	// Check if plan file exists
	planner := &install.FilePlanner{File: opts.planFilename}
	if !planner.PlanExists() {
		return fmt.Errorf("plan does not exist")
	}
	plan, err := planner.Read()
	if err != nil {
		return fmt.Errorf("error reading plan file: %v", err)
	}

	// Validate SSH connections
	if ok, errs := install.ValidatePlanSSHConnections(plan); !ok {
		util.PrintValidationErrors(out, errs)
		return fmt.Errorf("error getting info from cluster nodes")
	}

	lv, err := install.ListVersions(plan)
	if err != nil {
		return fmt.Errorf("error getting version: %v", err)
	}

	if opts.outputFormat == "json" {
		b, err := json.MarshalIndent(lv, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshalling struct: %v", err)
		}
		fmt.Fprintf(out, string(b))
		return nil
	}

	fmt.Fprintf(out, "Cluster: ")
	if lv.IsTransitioning {
		fmt.Fprintf(out, "Transitioning from v%v to v%v\n", lv.EarliestVersion, lv.LatestVersion)
	} else {
		fmt.Fprintf(out, "v%v\n", lv.LatestVersion)
	}
	fmt.Fprintln(out)
	fmt.Fprintf(out, "Nodes:\n")
	for _, listNode := range lv.Nodes {
		fmt.Fprintf(out, "  - %v: v%v %v\n", listNode.Node.IP, listNode.Version, listNode.Roles)
	}

	return err
}
