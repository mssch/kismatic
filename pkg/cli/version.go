package cli

import (
	"io"

	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/spf13/cobra"
)

// NewCmdVersion returns the version command
func NewCmdVersion(version string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display the Kismatic CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			util.PrettyPrintf(out, "Kismatic version: %s", version)
		},
	}
}
