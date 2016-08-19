package cli

import (
	"bufio"
	"fmt"
	"io"
	"strconv"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

func NewCmdInstall(in io.Reader, out io.Writer, plan install.PlanReaderWriter) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "install your Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doInstall(in, out, plan)
		},
	}

	return cmd
}

func doInstall(in io.Reader, out io.Writer, plan install.PlanReaderWriter) error {
	if !plan.Exists() {
		// Plan file not found, planning phase
		fmt.Fprintln(out, "Plan your Kismatic cluster:")

		etcdNodes, err := promptForInt(in, out, "Number of etcd nodes", 3)
		if err != nil {
			return fmt.Errorf("Error reading number of etcd nodes: %v", err)
		}
		if etcdNodes <= 0 {
			return fmt.Errorf("The number of etcd nodes must be greater than zero")
		}

		masterNodes, err := promptForInt(in, out, "Number of master nodes", 2)
		if err != nil {
			return fmt.Errorf("Error reading number of master nodes: %v", err)
		}
		if masterNodes <= 0 {
			return fmt.Errorf("The number of master nodes must be greater than zero")
		}

		workerNodes, err := promptForInt(in, out, "Number of worker nodes", 3)
		if err != nil {
			return fmt.Errorf("Error reading number of worker nodes: %v", err)
		}
		if workerNodes <= 0 {
			return fmt.Errorf("The number of worker nodes must be greater than zero")
		}

		fmt.Fprintf(out, "Generating installation plan file with %d etcd nodes, %d master nodes and %d worker nodes\n",
			etcdNodes, masterNodes, workerNodes)

		p := install.Plan{
			EtcdNodeCount:   etcdNodes,
			MasterNodeCount: masterNodes,
			WorkerNodeCount: workerNodes,
		}
		err = install.WritePlanTemplate(p, plan)
		if err != nil {
			return fmt.Errorf("error planning installation: %v", err)
		}
		fmt.Fprintf(out, "Generated installation plan file at %q\n", planFilename)
		return nil
	}

	// Plan file exists, validate and install
	fmt.Fprintln(out, "Found plan file.")
	p, err := plan.Read()
	if err != nil {
		return fmt.Errorf("error reading plan file: %v", err)
	}

	err = install.ValidatePlan(p)
	if err != nil {
		return fmt.Errorf("errors validating installation plan file: %v", err)
	}

	return nil
}

func promptForInt(in io.Reader, out io.Writer, prompt string, defaultValue int) (int, error) {
	fmt.Fprintf(out, "=> %s [%d]: ", prompt, defaultValue)
	s := bufio.NewScanner(in)
	// Scan the first token
	s.Scan()
	if s.Err() != nil {
		return defaultValue, fmt.Errorf("error reading number: %v", s.Err())
	}
	ans := s.Text()
	if ans == "" {
		return defaultValue, nil
	}
	// Convert input into integer
	i, err := strconv.Atoi(ans)
	if err != nil {
		return defaultValue, fmt.Errorf("%q is not a number", ans)
	}
	return i, nil
}
