package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/spf13/cobra"
)

type addWorkerOpts struct {
	CAConfigFile             string
	CASigningRequest         string
	CASigningProfile         string
	GeneratedAssetsDirectory string
	RestartServices          bool
	OutputFormat             string
	Verbose                  bool
	SkipPreFlight            bool
	WorkerInternalIP         string
}

// NewCmdAddWorker returns the command for adding workers to the cluster
func NewCmdAddWorker(out io.Writer, installOpts *installOpts) *cobra.Command {
	opts := &addWorkerOpts{}
	cmd := &cobra.Command{
		Use:   "add-worker WORKER_NAME WORKER_IP",
		Short: "add a Worker node to an existing Kismatic cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return cmd.Usage()
			}
			planner := &install.FilePlanner{File: installOpts.planFilename}
			if !planner.PlanExists() {
				return errors.New("add-worker can only be used with an existin plan file")
			}
			execOpts := install.ExecutorOptions{
				CAConfigFile:             opts.CAConfigFile,
				CASigningProfile:         opts.CASigningProfile,
				CASigningRequest:         opts.CASigningRequest,
				GeneratedAssetsDirectory: opts.GeneratedAssetsDirectory,
				RestartServices:          opts.RestartServices,
				OutputFormat:             opts.OutputFormat,
				Verbose:                  opts.Verbose,
				SkipCAGeneration:         true,
			}
			executor, err := install.NewExecutor(out, os.Stderr, execOpts)
			if err != nil {
				return err
			}
			plan, err := planner.Read()
			if err != nil {
				return fmt.Errorf("failed to read plan file: %v", err)
			}
			node := install.Node{
				Host:       args[0],
				IP:         args[1],
				InternalIP: opts.WorkerInternalIP,
			}
			if !opts.SkipPreFlight {
				if err := runPreFlightOnWorker(executor, *plan, node); err != nil {
					return err
				}
			}
			if err := executor.AddWorker(plan, node); err != nil {
				return err
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&opts.CASigningRequest, "ca-csr", "ansible/playbooks/tls/ca-csr.json", "path to the Certificate Authority CSR")
	cmd.Flags().StringVar(&opts.CAConfigFile, "ca-config", "ansible/playbooks/tls/ca-config.json", "path to the Certificate Authority configuration file")
	cmd.Flags().StringVar(&opts.CASigningProfile, "ca-signing-profile", "kubernetes", "name of the profile to be used for signing certificates")
	cmd.Flags().StringVar(&opts.GeneratedAssetsDirectory, "generated-assets-dir", "generated", "path to the directory where assets generated during the installation process are to be stored")
	cmd.Flags().BoolVar(&opts.RestartServices, "restart-services", false, "force restart clusters services (Use with care)")
	cmd.Flags().BoolVar(&opts.Verbose, "verbose", false, "enable verbose logging from the installation")
	cmd.Flags().StringVarP(&opts.OutputFormat, "output", "o", "simple", "installation output format. Supported options: simple|raw")
	cmd.Flags().BoolVar(&opts.SkipPreFlight, "skip-preflight", false, "skip pre-flight checks")
	cmd.Flags().StringVar(&opts.WorkerInternalIP, "worker-internal-ip", "", "the internal IP of the worker, if different than the IP.")
	return cmd
}

func runPreFlightOnWorker(executor install.Executor, plan install.Plan, workerNode install.Node) error {
	// use the original plan, but only run against the new worker
	preFlightPlan := plan
	preFlightPlan.Master.Nodes = []install.Node{}
	preFlightPlan.Etcd.Nodes = []install.Node{}
	preFlightPlan.Worker.Nodes = []install.Node{workerNode}
	return executor.RunPreFlightCheck(&preFlightPlan)
}
