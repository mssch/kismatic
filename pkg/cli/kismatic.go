package cli

import (
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

const planFilename = "kismatic-cluster.yaml"

// NewKismaticCommand creates the kismatic command
func NewKismaticCommand(version string, in io.Reader, out io.Writer) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "kismatic",
		Short: "kismatic is the main tool for managing your Kismatic cluster",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewCmdVersion(version, out))

	// Add Install sub-command
	planner := &install.FilePlanner{File: planFilename}
	executor, err := install.NewAnsibleExecutor(out, os.Stderr) // TODO: Do we want to parameterize stderr?
	if err != nil {
		return nil, err
	}
	cmd.AddCommand(NewCmdInstall(in, out, planner, executor))

	return cmd, nil
}
