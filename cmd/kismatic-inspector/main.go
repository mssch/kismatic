package main

import (
	"os"

	"github.com/apprenda/kismatic-platform/pkg/inspector/cmd"
)

/*
var local = pflag.Bool("local", false, "run checks locally, and print results to the console")
var node = pflag.String("node", "", "run checks against a node running the Kismatic checker in server mode (e.g. 10.5.6.220:8081)")
var serverPort = pflag.Int("server-port", 8081, "port number used for Kismatic checker in server mode")
var checkedTCPPorts = pflag.String("check-tcp-ports", "", "Comma separated list of TCP ports that should be bound for checking.")
)
*/

func main() {
	cmd := cmd.NewCmdKismaticInspector(os.Stdout)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
