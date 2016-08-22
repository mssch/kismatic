package cli

import (
	"io"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

const planFilename = "kismatic-cluster.yaml"

// NewKismaticCommand creates the kismatic command
func NewKismaticCommand(in io.Reader, out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kismatic",
		Short: "kismatic is the main tool for managing your Kismatic cluster",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Add sub-commands
	cmd.AddCommand(NewCmdInstall(in, out, &install.PlanFile{planFilename}))

	return cmd
}
