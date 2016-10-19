package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/util"
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
			workerHost := args[0]
			workerIP := args[1]
			return doAddWorker(out, installOpts.planFilename, opts, workerHost, workerIP)
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

func doAddWorker(out io.Writer, planFile string, opts *addWorkerOpts, workerHost string, workerIP string) error {
	planner := &install.FilePlanner{File: planFile}
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
		Host:       workerHost,
		IP:         workerIP,
		InternalIP: opts.WorkerInternalIP,
	}
	if _, errs := install.ValidateNode(&node); errs != nil {
		printValidationErrors(out, errs)
		return errors.New("information provided about the new worker node is invalid")
	}
	if _, errs := install.ValidatePlan(plan); errs != nil {
		printValidationErrors(out, errs)
		return errors.New("the plan file failed validation")
	}
	if err := ensureNodeIsNew(*plan, node); err != nil {
		return err
	}
	if !opts.SkipPreFlight {
		util.PrintHeader(out, "Running Pre-Flight Checks On New Worker", '=')
		if err := runPreFlightOnWorker(executor, *plan, node); err != nil {
			return err
		}
	}
	updatedPlan, err := executor.AddWorker(plan, node)
	if err != nil {
		return err
	}
	if err := planner.Write(updatedPlan); err != nil {
		return fmt.Errorf("error updating plan file to inlcude new worker node: %v", err)
	}
	return nil
}

// returns an error if the plan contains a worker that is "equivalent"
// to the new worker that is being added
func ensureNodeIsNew(plan install.Plan, newWorker install.Node) error {
	for _, n := range plan.Worker.Nodes {
		if n.Host == newWorker.Host {
			return fmt.Errorf("according to the plan file, the host name of the new node is already being used by another worker node")
		}
		if n.IP == newWorker.IP {
			return fmt.Errorf("according to the plan file, the IP of the new node is already being used by another worker node")
		}
		if newWorker.InternalIP != "" && n.InternalIP == newWorker.InternalIP {
			return fmt.Errorf("according to the plan file, the internal IP of the new node is already being used by another worker node")
		}
	}
	return nil
}

func runPreFlightOnWorker(executor install.Executor, plan install.Plan, workerNode install.Node) error {
	// use the original plan, but only run against the new worker
	preFlightPlan := plan
	preFlightPlan.Master.Nodes = []install.Node{}
	preFlightPlan.Etcd.Nodes = []install.Node{}
	preFlightPlan.Worker.Nodes = []install.Node{workerNode}
	return executor.RunPreFlightCheck(&preFlightPlan)
}
