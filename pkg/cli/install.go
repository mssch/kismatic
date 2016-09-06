package cli

import (
	"io"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

// NewCmdInstall creates a new install command
func NewCmdInstall(in io.Reader, out io.Writer) *cobra.Command {
	options := &install.CliOpts{}

	cmd := &cobra.Command{
		Use:   "install",
		Short: "install your Kismatic cluster",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	// Subcommands
	cmd.AddCommand(NewCmdPlan(in, out, options))
	cmd.AddCommand(NewCmdValidate(out, options))
	cmd.AddCommand(NewCmdApply(out, options))

	// PersistentFlags
	cmd.PersistentFlags().StringVarP(&options.PlanFilename, "plan-file", "f", "kismatic-cluster.yaml", "path to the installation plan file")

	return cmd
}
