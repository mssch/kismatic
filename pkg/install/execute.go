package install

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// The Executor will carry out the installation plan
type Executor interface {
	Install(p *Plan) error
}

type ansibleExecutor struct {
	out        io.Writer
	errOut     io.Writer
	pythonPath string
	ansibleBin string
}

// NewAnsibleExecutor returns an ansible based installation executor.
func NewAnsibleExecutor(out io.Writer, errOut io.Writer) (Executor, error) {
	ppath, err := getPythonPath()
	if err != nil {
		return nil, err
	}
	return &ansibleExecutor{
		out:        out,
		errOut:     errOut,
		pythonPath: ppath,
		ansibleBin: "./ansible/bin", // TODO: What's the best way to handle this?
	}, nil
}

func (e *ansibleExecutor) Install(p *Plan) error {
	inv := &bytes.Buffer{}
	writeHostGroup(inv, "etcd", p.Etcd.Nodes)
	writeHostGroup(inv, "master", p.Master.Nodes)
	writeHostGroup(inv, "worker", p.Worker.Nodes)

	err := ioutil.WriteFile("inventory.ini", inv.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing ansible inventory file: %v", err)
	}

	// run ansible
	err = e.runAnsiblePlaybook("inventory.ini", "./ansible/playbooks/kubernetes.yaml")
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}

	return nil
}

func (e *ansibleExecutor) runAnsiblePlaybook(inventoryFile, playbookFile string) error {
	cmd := exec.Command(fmt.Sprintf("%s/ansible-playbook", e.ansibleBin), "-i", inventoryFile, "-s", playbookFile)
	cmd.Stdout = e.out
	cmd.Stderr = e.errOut
	os.Setenv("PYTHONPATH", e.pythonPath)

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error running playbook: %v", err)
	}

	return nil
}

func writeHostGroup(inv io.Writer, groupName string, nodes []Node) {
	fmt.Fprintf(inv, "[%s]\n", groupName)
	for _, n := range nodes {
		internalIP := n.IP
		if n.InternalIP != "" {
			internalIP = n.InternalIP
		}
		fmt.Fprintf(inv, "%s ansible_host=%s internal_ipv4=%s\n", n.Host, n.IP, internalIP)
	}
}

func getPythonPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working dir: %v", err)
	}
	return fmt.Sprintf("%s/ansible/lib/python2.7/site-packages:%[1]s/ansible/lib64/python2.7/site-packages", wd), nil
}
