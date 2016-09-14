package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/apprenda/kismatic-platform/pkg/preflight"
	"github.com/spf13/pflag"
)

var local = pflag.Bool("local", false, "run checks locally, and print results to the console")
var node = pflag.String("node", "", "run checks against a node running the Kismatic checker in server mode (e.g. 10.5.6.220:8081)")
var serverPort = pflag.Int("server-port", 8081, "port number used for Kismatic checker in server mode")
var checkedTCPPorts = pflag.String("check-tcp-ports", "", "Comma separated list of TCP ports that should be bound for checking.")
var output = pflag.StringP("output", "o", "table", "set the result output type. Options are 'json', 'table'")

func main() {
	pflag.Parse()

	// TODO: Extract this to a config file?
	cr := &preflight.CheckRequest{
		BinaryDependencies:  []string{"iptables", "iptables-save", "iptables-restore", "ip", "nsenter", "mount", "umount"},
		PackageDependencies: []string{"glibc"},
	}

	s := preflight.Server{ListenPort: *serverPort}

	var printResults resultPrinter
	switch *output {
	case "table":
		printResults = printResultsAsTable
	case "json":
		printResults = printResultsAsJSON
	default:
		fmt.Fprintf(os.Stderr, "%q is not supported as an output option\n", *output)
		os.Exit(1)
	}

	if *local {
		results := s.RunChecks(cr)
		err := printResults(os.Stdout, results)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error printing results: %v\n", err)
			os.Exit(1)
		}
		// TODO: Exit with 1 if checks failed
		os.Exit(0)
	}

	if *node != "" {
		c := preflight.Client{TargetNode: *node}

		if *checkedTCPPorts != "" {
			// Add TCP Port checks to CheckRequest
			ports := strings.Split(*checkedTCPPorts, ",")
			for _, port := range ports {
				p, err := strconv.Atoi(port)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Invalid port %q\n", port)
					os.Exit(1)
				}
				cr.TCPPorts = append(cr.TCPPorts, p)
			}
		}

		results, err := c.RunChecks(cr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running checks against %q: %v\n", *node, err)
			os.Exit(1)
		}

		err = printResults(os.Stdout, results)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error printing results: %v\n", err)
		}

		for _, r := range results {
			if !r.Success {
				os.Exit(1)
			}
		}

		os.Exit(0)
	}

	fmt.Println("Listening on port", *serverPort)
	fmt.Printf("Run %s from another node to run checks remotely: %[1]s --node [NodeIP]:%d\n", os.Args[0], *serverPort)
	log.Fatal(s.Start())
}

type resultPrinter func(out io.Writer, r []preflight.CheckResult) error

func printResultsAsJSON(out io.Writer, r []preflight.CheckResult) error {
	err := json.NewEncoder(out).Encode(r)
	if err != nil {
		return fmt.Errorf("error marshaling results as JSON: %v", err)
	}
	return nil
}

func printResultsAsTable(out io.Writer, r []preflight.CheckResult) error {
	w := tabwriter.NewWriter(out, 1, 8, 4, '\t', 0)
	fmt.Fprintf(w, "CHECK\tSUCCESS\tMSG\n")
	for _, cr := range r {
		fmt.Fprintf(w, "%s\t%t\t%v\n", cr.Name, cr.Success, cr.Error)
	}
	w.Flush()
	return nil
}
