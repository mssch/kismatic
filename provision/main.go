package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/apprenda/kismatic-platform/integration"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "provision",
	Short: "Provision is a tool for making Kubernetes capable infrastructure",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
	SilenceUsage: true,
}

func init() {
	RootCmd.AddCommand(AWSCmd())
}

type AWSOpts struct {
	EtcdNodeCount   uint16
	MasterNodeCount uint16
	WorkerNodeCount uint16
	Ubuntu          bool
	LeaveArtifacts  bool
	RunKismatic     bool
	Plan            bool
}

// NewCmdAddWorker returns the command for adding workers to the cluster
func AWSCmd() *cobra.Command {
	opts := AWSOpts{}
	cmd := &cobra.Command{
		Use:   "aws",
		Short: "Creates a new cluster in AWS",
		Run: func(cmd *cobra.Command, args []string) {
			makeInfra(opts)
		},
	}

	cmd.Flags().Uint16VarP(&opts.EtcdNodeCount, "etcdNodeCount", "e", 1, "Count of etcd nodes to produce.")
	cmd.Flags().Uint16VarP(&opts.MasterNodeCount, "masterdNodeCount", "m", 1, "Count of master nodes to produce.")
	cmd.Flags().Uint16VarP(&opts.WorkerNodeCount, "workerNodeCount", "w", 1, "Count of worker nodes to produce.")
	cmd.Flags().BoolVarP(&opts.Ubuntu, "useUbuntu", "u", false, "If present, will install Ubuntu 16.04 rather than RHEL")
	//cmd.Flags().BoolVarP(&opts.LeaveArtifacts, "leaveArtifacts", "l", false, "If present, will leave behind all the artifacts produced while deploying kismatic. Use for debugging")
	cmd.Flags().BoolVarP(&opts.RunKismatic, "runKismatic", "k", false, "If present, a basic Kismatic install will occur immediately after nodes are provisioned.")
	cmd.Flags().BoolVarP(&opts.Plan, "plan", "p", false, "If present, generate a plan file in this directory referencing the newly created nodes")

	return cmd
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func makeInfra(opts AWSOpts) {
	os.Setenv("LEAVE_ARTIFACTS", "true")
	if opts.LeaveArtifacts {
		os.Setenv("LEAVE_INSTALL", "true")
	} else {
		os.Setenv("LEAVE_INSTALL", "")
	}
	if opts.RunKismatic {
		os.Setenv("BAIL_BEFORE_ANSIBLE", "")
		moveKismaticToTemp()
	} else {
		os.Setenv("BAIL_BEFORE_ANSIBLE", "true")
	}

	curDir, _ := os.Getwd()

	ami := integration.UbuntuEast
	if !opts.Ubuntu {
		ami = integration.CentosEast
	}

	cluster := integration.InstallBigKismatic(integration.NodeCount{Etcd: opts.EtcdNodeCount, Worker: opts.WorkerNodeCount, Master: opts.MasterNodeCount}, ami, false, false, false)

	if opts.Plan {
		os.Chdir(curDir)
		integration.MakePlanFile(&cluster)
	}

	fmt.Println("Your cluster is ready.\n")
	integration.PrintNodes(&cluster)
}

func moveKismaticToTemp() {
	kisPath := integration.CopyKismaticToTemp()

	fmt.Println("Unpacking kismatic to", kisPath)
	c := exec.Command("tar", "-zxf", "out/kismatic.tar.gz", "-C", kisPath)
	tarOut, tarErr := c.CombinedOutput()
	if tarErr != nil {
		log.Fatal("Error unpacking installer", string(tarOut), tarErr)
	}
	os.Chdir(kisPath)
}
