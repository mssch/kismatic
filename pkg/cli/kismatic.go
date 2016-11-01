package cli

import (
	"io"

	"github.com/spf13/cobra"
)

// NewKismaticCommand creates the kismatic command
func NewKismaticCommand(version string, in io.Reader, out io.Writer) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "kismatic",
		Short: "kismatic is the main tool for managing your Kismatic cluster",
		Long: `kismatic is the main tool for managing your Kismatic cluster
More documentation is availble at https://github.com/apprenda/kismatic-platform`,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(NewCmdVersion(version, out))
	cmd.AddCommand(NewCmdInstall(in, out))

	return cmd, nil
}
