package cmd

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector"
	"github.com/spf13/cobra"
)

// NewCmdServer returns the "server" command
func NewCmdServer(out io.Writer) *cobra.Command {
	var port int
	var nodeRole string
	cmd := &cobra.Command{
		Use:   "server",
		Short: "stand up the inspector server for running checks remotely",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(out, cmd.Parent().Name(), port, nodeRole)
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "The port number for standing up the Inspector server")
	cmd.Flags().StringVar(&nodeRole, "node-role", "", "The node's role in the cluster. Options are 'etcd', 'master', 'worker'")
	return cmd
}

func runServer(out io.Writer, commandName string, port int, nodeRole string) error {
	if nodeRole == "" {
		return fmt.Errorf("node role is required")
	}
	if nodeRole != "etcd" && nodeRole != "master" && nodeRole != "worker" {
		return fmt.Errorf("%s is not a valid node role", nodeRole)
	}
	s, err := inspector.NewServer(nodeRole, port)
	if err != nil {
		return fmt.Errorf("error starting up inspector server: %v", err)
	}
	fmt.Fprintf(out, "Inspector is listening on port %d\n", port)
	fmt.Fprintf(out, "Run %s from another node to run checks remotely: %[1]s client [NODE_IP]:%d\n", commandName, port)
	if err := s.Start(); err != nil {
		return err
	}
	return nil
}
