package cli

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdPlan creates a new install plan command
func NewCmdPlan(in io.Reader, out io.Writer, options *installOpts) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "plan your Kismatic cluster and gerenate a plan file",
		RunE: func(cmd *cobra.Command, args []string) error {
			planner := &install.FilePlanner{File: options.planFilename}
			return doPlan(in, out, planner, options)
		},
	}

	return cmd
}

func doPlan(in io.Reader, out io.Writer, planner install.Planner, options *installOpts) error {
	fmt.Fprintln(out, "Plan your Kismatic cluster:")

	// etcd nodes
	etcdNodes, err := util.PromptForInt(in, out, "Number of etcd nodes", 3)
	if err != nil {
		return fmt.Errorf("Error reading number of etcd nodes: %v", err)
	}
	if etcdNodes <= 0 {
		return fmt.Errorf("The number of etcd nodes must be greater than zero")
	}
	// master nodes
	masterNodes, err := util.PromptForInt(in, out, "Number of master nodes", 2)
	if err != nil {
		return fmt.Errorf("Error reading number of master nodes: %v", err)
	}
	if masterNodes <= 0 {
		return fmt.Errorf("The number of master nodes must be greater than zero")
	}
	// worker nodes
	workerNodes, err := util.PromptForInt(in, out, "Number of worker nodes", 3)
	if err != nil {
		return fmt.Errorf("Error reading number of worker nodes: %v", err)
	}
	if workerNodes <= 0 {
		return fmt.Errorf("The number of worker nodes must be greater than zero")
	}

	fmt.Fprintf(out, "Generating installation plan file with %d etcd nodes, %d master nodes and %d worker nodes\n",
		etcdNodes, masterNodes, workerNodes)

	// Create a plan
	masterNodeGroup := install.MasterNodeGroup{}
	masterNodeGroup.ExpectedCount = masterNodes
	plan := install.Plan{
		Etcd: install.NodeGroup{
			ExpectedCount: etcdNodes,
		},
		Master: masterNodeGroup,
		Worker: install.NodeGroup{
			ExpectedCount: workerNodes,
		},
	}
	// Write out the plan
	err = install.WritePlanTemplate(plan, planner)
	if err != nil {
		return fmt.Errorf("error planning installation: %v", err)
	}
	fmt.Fprintf(out, "Generated installation plan file at %q\n", options.planFilename)
	fmt.Fprintf(out, "Edit the file to further describe your cluster. Once ready, execute the install verify command to proceed.\n")

	return nil
}
