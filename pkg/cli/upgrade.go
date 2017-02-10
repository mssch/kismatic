package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/apprenda/kismatic/pkg/data"
	"github.com/apprenda/kismatic/pkg/install"
	"github.com/apprenda/kismatic/pkg/util"
	"github.com/spf13/cobra"
)

type upgradeOpts struct {
	generatedAssetsDir string
	verbose            bool
	outputFormat       string
	skipPreflight      bool
	online             bool
	planFile           string
	restartServices    bool
}

// NewCmdUpgrade returns the upgrade command
func NewCmdUpgrade(out io.Writer) *cobra.Command {
	var opts upgradeOpts
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "upgrade your Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.PersistentFlags().StringVar(&opts.generatedAssetsDir, "generated-assets-dir", "generated", "path to the directory where assets generated during the installation process will be stored")
	cmd.PersistentFlags().BoolVar(&opts.verbose, "verbose", false, "enable verbose logging from the installation")
	cmd.PersistentFlags().StringVarP(&opts.outputFormat, "output", "o", "simple", "installation output format (options \"simple\"|\"raw\")")
	cmd.PersistentFlags().BoolVar(&opts.skipPreflight, "skip-preflight", false, "skip upgrade pre-flight checks")
	cmd.PersistentFlags().BoolVar(&opts.restartServices, "restart-services", false, "force restart cluster services (Use with care)")
	addPlanFileFlag(cmd.PersistentFlags(), &opts.planFile)

	// Subcommands
	cmd.AddCommand(NewCmdUpgradeOffline(out, &opts))
	cmd.AddCommand(NewCmdUpgradeOnline(out, &opts))
	return cmd
}

// NewCmdUpgradeOffline returns the command for running offline upgrades
func NewCmdUpgradeOffline(out io.Writer, opts *upgradeOpts) *cobra.Command {
	cmd := cobra.Command{
		Use:   "offline",
		Short: "perform an offline upgrade of your Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return doUpgrade(out, opts)
		},
	}
	return &cmd
}

// NewCmdUpgradeOnline returns the command for running online upgrades
func NewCmdUpgradeOnline(out io.Writer, opts *upgradeOpts) *cobra.Command {
	cmd := cobra.Command{
		Use:   "online",
		Short: "perform an online upgrade of your Kubernetes cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.online = true
			return doUpgrade(out, opts)
		},
	}
	return &cmd
}

func doUpgrade(out io.Writer, opts *upgradeOpts) error {
	planFile := opts.planFile
	planner := install.FilePlanner{File: planFile}
	executorOpts := install.ExecutorOptions{
		GeneratedAssetsDirectory: opts.generatedAssetsDir,
		RestartServices:          opts.restartServices,
		OutputFormat:             opts.outputFormat,
		Verbose:                  opts.verbose,
	}
	executor, err := install.NewExecutor(out, os.Stderr, executorOpts)
	if err != nil {
		return err
	}
	util.PrintHeader(out, "Computing upgrade plan", '=')

	// Read plan file
	if !planner.PlanExists() {
		util.PrettyPrintErr(out, "Reading plan file")
		return fmt.Errorf("plan file %q does not exist", planFile)
	}
	util.PrettyPrintOk(out, "Reading plan file")
	plan, err := planner.Read()
	if err != nil {
		util.PrettyPrintErr(out, "Reading plan file")
		return fmt.Errorf("error reading plan file %q: %v", planFile, err)
	}

	// Validate SSH connectivity to nodes
	if ok, errs := install.ValidatePlanSSHConnections(plan); !ok {
		util.PrettyPrintErr(out, "Validate SSH connectivity to nodes")
		util.PrintValidationErrors(out, errs)
		return fmt.Errorf("SSH connectivity validation errors found")
	}
	util.PrettyPrintOk(out, "Validate SSH connectivity to nodes")

	// Figure out which nodes to upgrade
	cv, err := install.ListVersions(plan)
	if err != nil {
		return fmt.Errorf("error listing cluster versions: %v", err)
	}
	var toUpgrade []install.ListableNode
	var toSkip []install.ListableNode
	for _, n := range cv.Nodes {
		if install.IsOlderVersion(n.Version) {
			toUpgrade = append(toUpgrade, n)
		} else {
			toSkip = append(toSkip, n)
		}
	}

	// Print the nodes that will be skipped
	if len(toSkip) > 0 {
		util.PrintHeader(out, "Skipping nodes", '=')
		for _, n := range toSkip {
			util.PrettyPrintOk(out, "- %q is at the target version %q", n.Node.IP, n.Version)
		}
		fmt.Fprintln(out)
	}

	// Print message if there's no work to do
	if len(toUpgrade) == 0 {
		fmt.Fprintln(out, "All nodes are at the target version. Skipping node upgrades.")
	} else {
		if err = upgradeNodes(out, *plan, *opts, toUpgrade, executor); err != nil {
			return err
		}
	}

	if plan.ConfigureDockerRegistry() && plan.Cluster.DisconnectedInstallation {
		util.PrintHeader(out, "Upgrade Docker Registry", '=')
		if err := executor.UpgradeDockerRegistry(*plan); err != nil {
			return fmt.Errorf("Failed to upgrade docker registry: %v", err)
		}
	}

	// Upgrade the cluster services
	util.PrintHeader(out, "Upgrade Cluster Services", '=')
	if err := executor.UpgradeClusterServices(*plan); err != nil {
		return fmt.Errorf("Failed to upgrade cluster services: %v", err)
	}

	util.PrintHeader(out, "Smoke Test Cluster", '=')
	if err := executor.RunSmokeTest(plan); err != nil {
		return fmt.Errorf("Smoke test failed: %v", err)
	}

	fmt.Fprintln(out)
	util.PrintColor(out, util.Green, "Upgrade complete\n")
	fmt.Fprintln(out)
	return nil
}

func upgradeNodes(out io.Writer, plan install.Plan, opts upgradeOpts, toUpgrade []install.ListableNode, executor install.Executor) error {
	// Validate that we are able to perform an online upgrade
	if opts.online {
		util.PrintHeader(out, "Validate Online Upgrade", '=')
		var foundErrs bool
		// Use the first master node for running kubectl
		client, _ := plan.GetSSHClient(plan.Master.Nodes[0].Host)
		kubeClient := data.RemoteKubectl{SSHClient: client}
		for _, node := range toUpgrade {
			util.PrettyPrint(out, "Node %q", node.Node.Host)
			errs := install.DetectNodeUpgradeSafety(plan, node.Node, kubeClient)
			if len(errs) != 0 {
				foundErrs = true
				util.PrintError(out)
				fmt.Fprintln(out)
				for _, err := range errs {
					fmt.Println("-", err.Error())
				}
			} else {
				util.PrintOkln(out)
			}
		}
		if foundErrs {
			return errors.New("Unable to perform an online upgrade due to the unsafe conditions detected.")
		}
	}

	// Run upgrade preflight on the nodes that are to be UpgradeNodes
	if !opts.skipPreflight {
		util.PrintHeader(out, "Validating nodes", '=')
		if err := executor.RunUpgradePreFlightCheck(&plan); err != nil {
			return fmt.Errorf("Upgrade preflight check failed: %v", err)
		}
	}

	// Run the upgrade on the nodes that need it
	if err := executor.UpgradeNodes(plan, toUpgrade); err != nil {
		return fmt.Errorf("Failed to upgrade nodes: %v", err)
	}
	return nil
}
