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

	cluster := integration.InstallBigKismatic(integration.NodeCount{Etcd: 1, Master: 1, Worker: 1}, integration.UbuntuEast)

	fmt.Println("Your cluster is ready.\n")
	integration.PrintNodes(&cluster)
}
