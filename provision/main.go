package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/apprenda/kismatic-platform/integration"
)

func main() {
	if err := os.Setenv("LEAVE_ARTIFACTS", "true"); err != nil {
		log.Fatal("Error setting environment variable", err)
	}
	os.Setenv("BAIL_BEFORE_ANSIBLE", "")

	kisPath := integration.CopyKismaticToTemp()

	fmt.Println("Unpacking kismatic to", kisPath)
	c := exec.Command("tar", "-zxf", "out/kismatic.tar.gz", "-C", kisPath)
	tarOut, tarErr := c.CombinedOutput()
	if tarErr != nil {
		log.Fatal("Error unpacking installer", string(tarOut), tarErr)
	}
	os.Chdir(kisPath)

	//defer os.RemoveAll(kisPath)

	cluster := integration.InstallKismatic(integration.AMIUbuntu1604USEAST, "ubuntu")

	fmt.Println("Your cluster is ready.\n")
	fmt.Println("Etcd Nodes:")
	for _, etcd := range cluster.Etcd {
		fmt.Println(etcd.Publicip)
	}
	fmt.Println("Master Nodes:")
	for _, master := range cluster.Master {
		fmt.Println(master.Publicip)
	}
	fmt.Println("Worker Nodes:")
	for _, worker := range cluster.Worker {
		fmt.Println(worker.Publicip)
	}
}
