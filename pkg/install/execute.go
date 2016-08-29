package install

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
	pki        PKI
}

type ansibleVars struct {
	TLSDirectory           string `json:"tls_directory"`
	KubernetesServicesCIDR string `json:"kubernetes_services_cidr"`
	KubernetesPodsCIDR     string `json:"kubernetes_pods_cidr"`
}

func (av *ansibleVars) CommandLineVars() (string, error) {
	b, err := json.Marshal(av)
	if err != nil {
		return "", fmt.Errorf("error marshaling ansible vars")
	}
	return string(b), nil
}

// NewAnsibleExecutor returns an ansible based installation executor.
func NewAnsibleExecutor(out io.Writer, errOut io.Writer, pki PKI) (Executor, error) {
	ppath, err := getPythonPath()
	if err != nil {
		return nil, err
	}
	return &ansibleExecutor{
		out:        out,
		errOut:     errOut,
		pythonPath: ppath,
		ansibleBin: "./ansible/bin", // TODO: What's the best way to handle this?
		pki:        pki,
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

	// Generate certs for the cluster
	err = e.pki.GenerateClusterCerts(p)
	if err != nil {
		return fmt.Errorf("error generating certificates for the cluster: %v", err)
	}

	tlsDir, err := filepath.Abs(e.pki.Location())
	if err != nil {
		return fmt.Errorf("error getting absolute path from cert location: %v", err)
	}
	vars := ansibleVars{
		TLSDirectory:           tlsDir,
		KubernetesServicesCIDR: p.Cluster.Networking.ServiceCIDRBlock,
		KubernetesPodsCIDR:     p.Cluster.Networking.PodCIDRBlock,
	}

	// run ansible
	err = e.runAnsiblePlaybook("inventory.ini", "./ansible/playbooks/kubernetes.yaml", vars)
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}

	return nil
}

func (e *ansibleExecutor) runAnsiblePlaybook(inventoryFile, playbookFile string, vars ansibleVars) error {
	extraVars, err := vars.CommandLineVars()
	if err != nil {
		return fmt.Errorf("error getting vars: %v", err)
	}

	cmd := exec.Command(fmt.Sprintf("%s/ansible-playbook", e.ansibleBin), "-i", inventoryFile, "-s", playbookFile, "--extra-vars", extraVars)
	cmd.Stdout = e.out
	cmd.Stderr = e.errOut
	os.Setenv("PYTHONPATH", e.pythonPath)

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

func getPythonPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working dir: %v", err)
	}
	return fmt.Sprintf("%s/ansible/lib/python2.7/site-packages:%[1]s/ansible/lib64/python2.7/site-packages", wd), nil
}
