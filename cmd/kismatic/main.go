package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

// Set via linker flag
var version string

func getPythonPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting working dir: %v", err)
	}
	return fmt.Sprintf("%s/ansible/lib/python2.7/site-packages:%[1]s/ansible/lib64/python2.7/site-packages", wd), nil
}

func main() {
	fmt.Println("Hello from Kismatic. Version:", version)

	// Create command with vendored ansible
	cmd := exec.Command("./ansible/bin/ansible", "localhost", "-a", "uname -a")
	ppath, err := getPythonPath()
	if err != nil {
		log.Fatalf("Failed to get python path: %v", err)
	}
	os.Setenv("PYTHONPATH", ppath)
	//cmd.Env = []string{"PYTHONPATH=" + ppath}

	// Run command with ansible and print output
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s\n", string(out))
	if err != nil {
		log.Fatalf("Error running ansible: %v", err)
	}
	fmt.Println("Ran ansible command")
}
