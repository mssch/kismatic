package main

import (
	"fmt"
	"log"
	"os"

	"github.com/apprenda/kismatic-platform/integration"
)

func main() {
	if err := os.Setenv("LEAVE_ARTIFACTS", "true"); err != nil {
		log.Fatal("Error setting environment variable", err)
	}
	os.Setenv("BAIL_BEFORE_ANSIBLE", "true")

	kisPath, err := integration.ExtractKismaticToTemp()
	if err != nil {
		log.Fatalln("Error unpacking installer", err)
	}
	os.Chdir(kisPath)

	aws, ok := integration.AWSClientFromEnvironment()
	if !ok {
		log.Fatal("Required AWS environment variables not defined")
	}

	nodes, err := aws.ProvisionNodes(integration.NodeCount{Etcd: 1, Master: 1, Worker: 1}, integration.Ubuntu1604LTS)
	if err != nil {
		aws.TerminateNodes(nodes)
		log.Fatal("Failed to provision nodes")
	}

	fmt.Println("Your cluster is ready.")
	fmt.Println(nodes)
}
