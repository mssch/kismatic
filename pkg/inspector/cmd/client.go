package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCmdClient returns the "client" command
func NewCmdClient(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "client",
		Short: "run the inspector against a remote inspector server",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	return cmd
}
