package main

import (
	"os"

	"github.com/apprenda/kismatic-platform/pkg/cli"
	"github.com/apprenda/kismatic-platform/pkg/util"
	"github.com/spf13/cobra/doc"
)

// Set via linker flag
var version string

func main() {

	cmd, err := cli.NewKismaticCommand(version, os.Stdin, os.Stdout)
	doc.GenMarkdownTree(cmd, "./cli-docs")
	if err != nil {
		util.PrintColor(os.Stderr, util.Red, "Error initializing command: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		util.PrintColor(os.Stderr, util.Red, "%v\n", err)
		os.Exit(1)
	}

}
