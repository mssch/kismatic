package main // Set via linker flag
import (
	"os"

	"github.com/apprenda/kismatic-platform/pkg/cli"
	"github.com/spf13/cobra/doc"
)

var version string

func main() {

	cmd, _ := cli.NewKismaticCommand(version, os.Stdin, os.Stdout)
	doc.GenMarkdownTree(cmd, "./kismatic-cli-docs")
}
