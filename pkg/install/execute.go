package install

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

func ExecutePlan(p *Plan) error {

	// produce ansible inventory and variables files
	inv := &bytes.Buffer{}
	writeHostGroup(inv, "etcd", p.Etcd.Nodes)
	writeHostGroup(inv, "master", p.Master.Nodes)
	writeHostGroup(inv, "worker", p.Worker.Nodes)

	err := ioutil.WriteFile("inventory.ini", inv.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing ansible inventory file: %v", err)
	}

	// run ansible
	err = runAnsiblePlaybook("inventory.ini", "./ansible/playbooks/kubernetes.yaml")
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}

	return nil
}

func getPythonPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working dir: %v", err)
	}
	return fmt.Sprintf("%s/ansible/lib/python2.7/site-packages:%[1]s/ansible/lib64/python2.7/site-packages", wd), nil
}

func runAnsiblePlaybook(inventoryFile, playbookFile string) error {
	cmd := exec.Command("./ansible/bin/ansible-playbook", "-i", inventoryFile, "-s", playbookFile, "--extra-vars", "@./ansible/playbooks/runtime_vars.yaml")
	fmt.Println(cmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ppath, err := getPythonPath()
	if err != nil {
		return fmt.Errorf("error building python path: %v", err)
	}
	os.Setenv("PYTHONPATH", ppath)

	err = cmd.Run()

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
