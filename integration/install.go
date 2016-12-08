package integration

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	homedir "github.com/mitchellh/go-homedir"
	. "github.com/onsi/ginkgo"
)

func leaveIt() bool {
	return os.Getenv("LEAVE_ARTIFACTS") != ""
}
func bailBeforeAnsible() bool {
	return os.Getenv("BAIL_BEFORE_ANSIBLE") != ""
}

func GetSSHKeyFile() (string, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, ".ssh", "kismatic-integration-testing.pem"), nil
}

type installOptions struct {
	allowPackageInstallation    bool
	autoConfigureDockerRegistry bool
	dockerRegistryIP            string
	dockerRegistryPort          int
	dockerRegistryCAPath        string
}

func installKismaticMini(node NodeDeets, sshKey string) error {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Building a plan to set up an overlay network cluster on this hardware")
	sshUser := node.SSHUser
	plan := PlanAWS{
		Etcd:                     []NodeDeets{node},
		Master:                   []NodeDeets{node},
		Worker:                   []NodeDeets{node},
		Ingress:                  []NodeDeets{node},
		MasterNodeFQDN:           node.Hostname,
		MasterNodeShortName:      node.Hostname,
		SSHKeyFile:               sshKey,
		SSHUser:                  sshUser,
		AllowPackageInstallation: true,
	}

	By("Writing plan file out to disk")
	f, err := os.Create("kismatic-testing.yaml")
	FailIfError(err, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = template.Execute(w, &plan)
	FailIfError(err, "Error filling in plan template")
	w.Flush()

	By("Validing our plan")
	cmd := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	FailIfError(err, "Error validating plan")

	By("Punch it Chewie!")
	cmd = exec.Command("./kismatic", "install", "apply", "-f", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func installKismatic(nodes provisionedNodes, installOpts installOptions, sshKey string) error {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Building a plan to set up an overlay network cluster on this hardware")
	sshUser := nodes.master[0].SSHUser

	masterDNS := nodes.master[0].Hostname
	if nodes.dnsRecord != nil && nodes.dnsRecord.Name != "" {
		masterDNS = nodes.dnsRecord.Name
	}
	plan := PlanAWS{
		AllowPackageInstallation: installOpts.allowPackageInstallation,
		Etcd:                nodes.etcd,
		Master:              nodes.master,
		Worker:              nodes.worker,
		Ingress:             nodes.ingress,
		MasterNodeFQDN:      masterDNS,
		MasterNodeShortName: masterDNS,
		SSHKeyFile:          sshKey,
		SSHUser:             sshUser,
		AutoConfiguredDockerRegistry: installOpts.autoConfigureDockerRegistry,
		DockerRegistryCAPath:         installOpts.dockerRegistryCAPath,
		DockerRegistryIP:             installOpts.dockerRegistryIP,
		DockerRegistryPort:           installOpts.dockerRegistryPort,
	}

	f, err := os.Create("kismatic-testing.yaml")
	FailIfError(err, "Error creating plan")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = template.Execute(w, &plan)
	FailIfError(err, "Error filling in plan template")
	w.Flush()

	By("Punch it Chewie!")
	cmd := exec.Command("./kismatic", "install", "apply", "-f", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()

}

func verifyMasterNodeFailure(nodes provisionedNodes, provisioner infrastructureProvisioner, sshKey string) error {
	By("Removing a Kubernetes master node")
	if err := provisioner.TerminateNode(nodes.master[0]); err != nil {
		return fmt.Errorf("Could not remove node: %v", err)
	}

	By("Rerunning Kuberang")
	if err := runViaSSH([]string{"sudo kuberang"}, []NodeDeets{nodes.master[1]}, sshKey, 5*time.Minute); err != nil {
		return fmt.Errorf("Failed to run kuberang: %v", err)
	}

	return nil
}

func verifyIngressNodes(nodes provisionedNodes, sshKey string) error {
	By("Adding a service and an ingress resource")
	addIngressResource(nodes.master[0], sshKey)

	By("Verifying the service is accessible via the ingress point(s)")
	for _, ingNode := range nodes.ingress {
		if err := verifyIngressPoint(ingNode, sshKey); err != nil {
			return err
		}
	}

	return nil
}

func verifyIngressNode(node NodeDeets, sshKey string) error {
	By("Adding a service and an ingress resource")
	addIngressResource(node, sshKey)

	By("Verifying the service is accessible via the ingress point(s)")
	return verifyIngressPoint(node, sshKey)
}

func addIngressResource(node NodeDeets, sshKey string) {
	err := copyFileToRemote("test-resources/ingress.yaml", "/tmp/ingress.yaml", node, sshKey, 1*time.Minute)
	FailIfError(err, "Error copying ingress test file")

	err = runViaSSH([]string{"sudo kubectl apply -f /tmp/ingress.yaml"}, []NodeDeets{node}, sshKey, 1*time.Minute)
	FailIfError(err, "Error creating ingress resource")
}

func verifyIngressPoint(node NodeDeets, sshKey string) error {
	client := http.Client{
		Timeout: 1000 * time.Millisecond,
	}
	url := "http://" + node.PublicIP + "/echo"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("Could not create request for ingress via %s, %v", url, err)
	}
	// Set the host header since this is not a real domain, curl $IP/echo -H 'Host: kismaticintegration.com'
	req.Host = "kismaticintegration.com"
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Could not reach ingress via %s, %v", url, err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("Ingress status code is not 200, got %d vi %s", resp.StatusCode, url)
	}

	return nil
}

func installKismaticWithABadNode() {
	By("Building a template")
	template, err := template.New("planAWSOverlay").Parse(planAWSOverlay)
	FailIfError(err, "Couldn't parse template")

	By("Faking infrastructure")
	fakeNode := NodeDeets{
		id:       "FakeId",
		PublicIP: "10.0.0.0",
		Hostname: "FakeHostname",
	}

	By("Building a plan to set up an overlay network cluster on this hardware")
	sshKey, err := GetSSHKeyFile()
	FailIfError(err, "Error getting SSH Key file")
	plan := PlanAWS{
		Etcd:                []NodeDeets{fakeNode},
		Master:              []NodeDeets{fakeNode},
		Worker:              []NodeDeets{fakeNode},
		Ingress:             []NodeDeets{fakeNode},
		MasterNodeFQDN:      "yep.nope",
		MasterNodeShortName: "yep",
		SSHUser:             "Billy Rubin",
		SSHKeyFile:          sshKey,
	}
	By("Writing plan file out to disk")
	f, err := os.Create("kismatic-testing.yaml")
	FailIfError(err, "Error waiting for nodes")
	defer f.Close()
	w := bufio.NewWriter(f)
	err = template.Execute(w, &plan)
	FailIfError(err, "Error filling in plan template")
	w.Flush()
	f.Close()

	By("Validing our plan")
	cmd := exec.Command("./kismatic", "install", "validate", "-f", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		Fail("Validation succeeeded even though it shouldn't have")
	}

	By("Well, try it anyway")
	cmd = exec.Command("./kismatic", "install", "apply", "-f", f.Name())
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		Fail("Application succeeeded even though it shouldn't have")
	}
}

func completesInTime(dothis func(), howLong time.Duration) bool {
	c1 := make(chan string, 1)
	go func() {
		dothis()
		c1 <- "completed"
	}()

	select {
	case <-c1:
		return true
	case <-time.After(howLong):
		return false
	}
}

func FailIfError(err error, message ...string) {
	if err != nil {
		log.Printf(message[0]+": %v\n%v", err, message[1:])
		Fail(message[0])
	}
}

func FailIfSuccess(err error, message ...string) {
	if err == nil {
		Fail("Expected failure")
	}
}

func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
