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
	var nodeRoles string
	cmd := &cobra.Command{
		Use:   "server",
		Short: "stand up the inspector server for running checks remotely",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(out, cmd.Parent().Name(), port, nodeRoles)
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "the port number for standing up the Inspector server")
	cmd.Flags().StringVar(&nodeRoles, "node-roles", "", "comma-separated list of the node's roles. Valid roles are 'etcd', 'master', 'worker'")
	return cmd
}

func runServer(out io.Writer, commandName string, port int, nodeRoles string) error {
	if nodeRoles == "" {
		return fmt.Errorf("--node-roles is required")
	}
	roles, err := getNodeRoles(nodeRoles)
	if err != nil {
		return err
	}
	s, err := inspector.NewServer(roles, port)
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
