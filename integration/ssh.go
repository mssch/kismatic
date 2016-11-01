package integration

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

// Run the given command(s) as the given user on all hosts via SSH within the given period
func RunViaSSH(cmds []string, user string, hosts []NodeDeets, period time.Duration) bool {
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
			for _, cmd := range cmds {
				results <- executeCmd(cmd, hostname, config)
			}
		}(host.PublicIP)
	}

	for i := 0; i < len(hosts)*len(cmds); i++ {
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

func CopyFileToRemote(file string, destFile string, user string, hosts []NodeDeets, period time.Duration) bool {
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

func executeCmd(cmd, hostname string, config *ssh.ClientConfig) string {
	ver := exec.Command("ssh", "-o", "StrictHostKeyChecking no", "-t", "-t", "-i", os.Getenv("HOME")+"/.ssh/kismatic-integration-testing.pem", config.User+"@"+hostname, cmd)
	ver.Stdin = os.Stdin

	verbytes, verErr := ver.CombinedOutput()
	if verErr != nil {
		fmt.Printf("Oops: %v", verErr)
	}
	verText := string(verbytes)

	return hostname + ": " + verText
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
