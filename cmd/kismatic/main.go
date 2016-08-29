package main

import (
	"fmt"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/cli"
)

// Set via linker flag
var version string

func main() {

	cmd, err := cli.NewKismaticCommand(version, os.Stdin, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing command: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running command: %v\n", err)
		os.Exit(1)
	}

}
