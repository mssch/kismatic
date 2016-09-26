package cmd

import (
	"fmt"
	"io"

	"github.com/apprenda/kismatic-platform/pkg/inspector"
	"github.com/spf13/cobra"
)

func NewCmdServer(out io.Writer) *cobra.Command {
	var port int
	cmd := &cobra.Command{
		Use:   "server",
		Short: "stand up the inspector server for running checks remotely",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer(out, port, cmd.Parent().Name())
		},
	}
	cmd.Flags().IntVar(&port, "port", 8080, "The port number for standing up the Kismatic Inspector server")
	return cmd
}

func runServer(out io.Writer, port int, commandName string) error {
	s := inspector.Server{
		ListenPort: port,
	}
	fmt.Fprintf(out, "Inspector is listening on port %d\n", port)
	fmt.Fprintf(out, "Run %s from another node to run checks remotely: %[1]s client [NODE_IP]:%d\n", commandName, port)
	if err := s.Start(); err != nil {
		return err
	}
	return nil
}
