package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

const long string = `The Kismatic Inspector verifies the infrastructure that has
been provisioned for installing a Kubernetes cluster.
`

// NewCmdKismaticInspector builds the kismatic-inspector command
func NewCmdKismaticInspector(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kismatic-inspector",
		Short: "kismatic-inspector verifies infrastructure to be used for installing Kismatic",
		Long:  long,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		SilenceUsage: true,
	}
	cmd.AddCommand(NewCmdClient(out))
	cmd.AddCommand(NewCmdServer(out))
	cmd.AddCommand(NewCmdLocal(out))
	cmd.AddCommand(NewCmdRules(out))
	return cmd
}
