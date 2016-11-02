package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

func runViaSSH(cmds []string, hosts []NodeDeets, sshKey string, period time.Duration) error {
	timeout := time.After(period)
	bail := make(chan struct{})
	cmdSuccess := make(chan bool)
	// Create a goroutine per host. Each goroutine runs the commands serially on the host
	// until one of these is true:
	// a) all commands were executed successfully,
	// b) an error occurred when running a command,
	// c) the goroutine got a signal to bail
	for _, host := range hosts {
		go func(node NodeDeets) {
			for _, cmd := range cmds {
				res, err := executeCmd(cmd, node.PublicIP, node.SSHUser, sshKey)
				fmt.Println(res)
				select {
				case cmdSuccess <- err == nil:
					if err != nil {
						return
					}
				case <-bail:
					return
				}
			}
		}(host)
	}

	// The bail channel is closed if we encounter an error, or if the timeout is reached.
	// This will signal all goroutines to return.
	defer close(bail)

	// At most, we will get a total of hosts * cmds status messages in the channel.
	for i := 0; i < len(hosts)*len(cmds); i++ {
		select {
		case ok := <-cmdSuccess:
			if !ok {
				return fmt.Errorf("error running command on node")
			}
		case <-timeout:
			return fmt.Errorf("timed out running commands on nodes")
		}
	}
	return nil
}

func executeCmd(cmd, hostname, user, sshKey string) (string, error) {
	sshCmd := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "-t", "-t", "-i", sshKey, user+"@"+hostname, cmd)
	sshCmd.Stdin = os.Stdin
	sshOut, sshErr := sshCmd.CombinedOutput()
	return hostname + ": " + string(sshOut), sshErr
}

func copyFileToRemote(file string, destFile string, user string, hosts []NodeDeets, period time.Duration) bool {
	results := make(chan string, 10)
	success := true
	timeout := time.After(period)

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(os.Getenv("HOME") + "/.ssh/kismatic-integration-testing.pem"),
		},
	}

	for _, host := range hosts {
		go func(hostname string) {
			results <- scpFile(file, destFile, hostname, config)
		}(host.PublicIP)
	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Printf("%v timed out!", hosts[i])
			return false
		}
	}
	return success
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

// BlockUntilSSHOpen waits until the node with the given IP is accessible via SSH.
func BlockUntilSSHOpen(publicIP, sshUser, sshKey string) {
	for {
		cmd := exec.Command("ssh")
		cmd.Args = append(cmd.Args, "-i", sshKey)
		cmd.Args = append(cmd.Args, "-o", "ConnectTimeout=5")
		cmd.Args = append(cmd.Args, "-o", "BatchMode=yes")
		cmd.Args = append(cmd.Args, "-o", "StrictHostKeyChecking=no")
		cmd.Args = append(cmd.Args, fmt.Sprintf("%s@%s", sshUser, publicIP), "exit") // just call exit if we are able to connect
		if err := cmd.Run(); err == nil {
			// command succeeded
			fmt.Println()
			return
		}
		fmt.Printf("?")
		time.Sleep(3 * time.Second)
	}
}

func scpFile(filePath string, destFilePath string, hostname string, config *ssh.ClientConfig) string {
	ver := exec.Command("scp", "-o", "StrictHostKeyChecking no", "-i", os.Getenv("HOME")+"/.ssh/kismatic-integration-testing.pem", filePath, config.User+"@"+hostname+":"+destFilePath)
	ver.Stdin = os.Stdin

	verbytes, verErr := ver.CombinedOutput()
	if verErr != nil {
		fmt.Printf("Oops: %v", verErr)
	}
	verText := string(verbytes)

	return hostname + ": " + verText
}
