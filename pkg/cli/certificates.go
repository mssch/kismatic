package cli

import (
	"io"

	"github.com/spf13/cobra"
)

// NewCmdCertificates creates a new certificates command
func NewCmdCertificates(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "certificates",
		Short: "Manage cluster certificates",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewCmdGenerate(out))

	return cmd
}
