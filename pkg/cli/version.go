package cli

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

// NewCmdVersion returns the version command
func NewCmdVersion(version string, out io.Writer) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "display the Kismatic CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Fprintf(out, "Kismatic version: %s\n", version)
		},
	}
}
