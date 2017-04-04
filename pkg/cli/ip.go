package cli

import (
	"errors"
	"fmt"
	"io"

	"github.com/apprenda/kismatic/pkg/install"
	"github.com/spf13/cobra"
)

type ipOpts struct {
	planFilename string
}

// NewCmdIP prints the cluster's IP
func NewCmdIP(out io.Writer) *cobra.Command {
	opts := &ipOpts{}

	cmd := &cobra.Command{
		Use:   "ip",
		Short: "retrieve the IP address of the cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("Unexpected args: %v", args)
			}
			planner := &install.FilePlanner{File: opts.planFilename}
			return doIP(out, planner, opts)
		},
	}

	// PersistentFlags
	cmd.PersistentFlags().StringVarP(&opts.planFilename, "plan-file", "f", "kismatic-cluster.yaml", "path to the installation plan file")

	return cmd
}

func doIP(out io.Writer, planner install.Planner, opts *ipOpts) error {
	// Check if plan file exists
	if !planner.PlanExists() {
		return planFileNotFoundErr{filename: opts.planFilename}
	}
	plan, err := planner.Read()
	if err != nil {
		return fmt.Errorf("error reading plan file: %v", err)
	}
	address, err := getClusterAddress(*plan)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, address)
	return nil
}

func getClusterAddress(plan install.Plan) (string, error) {
	if plan.Master.LoadBalancedFQDN == "" {
		return "", errors.New("Master load balanced FQDN is not set in the plan file")
	}
	return plan.Master.LoadBalancedFQDN, nil
}
