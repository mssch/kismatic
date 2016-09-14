package main

import (
	"os"

	"github.com/apprenda/kismatic-platform/pkg/cli"
	"github.com/apprenda/kismatic-platform/pkg/util"
)

// Set via linker flag
var version string

func main() {

	cmd, err := cli.NewKismaticCommand(version, os.Stdin, os.Stdout)
	if err != nil {
		util.PrintErrorf(os.Stderr, "Error initializing command: %v", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		util.PrintErrorf(os.Stderr, "Error running command: %v", err)
		os.Exit(1)
	}

}
