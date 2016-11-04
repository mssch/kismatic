package main

import (
	"os"

	"github.com/apprenda/kismatic/pkg/cli"
	"github.com/apprenda/kismatic/pkg/util"
)

// Set via linker flag
var version string

func main() {

	cmd, err := cli.NewKismaticCommand(version, os.Stdin, os.Stdout)
	if err != nil {
		util.PrintColor(os.Stderr, util.Red, "Error initializing command: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		util.PrintColor(os.Stderr, util.Red, "%v\n", err)
		os.Exit(1)
	}

}
