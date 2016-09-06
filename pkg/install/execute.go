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
	"strconv"
)

// The Executor will carry out the installation plan
type Executor interface {
	Install(p *Plan, av *AnsibleVars) error
	GetVars(p *Plan, options *CliOpts) (*AnsibleVars, error)
}

type ansibleExecutor struct {
	out        io.Writer
	errOut     io.Writer
	pythonPath string
	ansibleDir string
	certsDir   string
}

// CommandLineVars returns "--extra-vars" as json string
func (av *AnsibleVars) CommandLineVars() (string, error) {
	b, err := json.Marshal(av)
	if err != nil {
		return "", fmt.Errorf("error marshaling ansible vars")
	}
	return string(b), nil
}

// NewAnsibleExecutor returns an ansible based installation executor.
func NewAnsibleExecutor(out io.Writer, errOut io.Writer, certsDir string) (Executor, error) {
	ppath, err := getPythonPath()
	if err != nil {
		return nil, err
	}
	return &ansibleExecutor{
		out:        out,
		errOut:     errOut,
		pythonPath: ppath,
		ansibleDir: "ansible", // TODO: What's the best way to handle this?
		certsDir:   certsDir,
	}, nil
}

func (e *ansibleExecutor) GetVars(p *Plan, options *CliOpts) (*AnsibleVars, error) {
	tlsDir, err := filepath.Abs(e.certsDir)
	if err != nil {
		return &AnsibleVars{}, fmt.Errorf("error getting absolute path from cert location: %v", err)
	}

	dnsServiceIP, err := getDNSServiceIP(p)
	if err != nil {
		return &AnsibleVars{}, fmt.Errorf("error getting DNS servie IP address: %v", err)
	}

	vars := AnsibleVars{
		ClusterName:                   p.Cluster.Name,
		AdminPassword:                 p.Cluster.AdminPassword,
		TLSDirectory:                  tlsDir,
		KubernetesServicesCIDR:        p.Cluster.Networking.ServiceCIDRBlock,
		KubernetesPodsCIDR:            p.Cluster.Networking.PodCIDRBlock,
		KubernetesDNSServiceIP:        dnsServiceIP,
		CalicoNetworkType:             p.Cluster.Networking.Type,
		LocalRepository:               p.Cluster.LocalRepository,
		ForceEtcdRestart:              strconv.FormatBool(options.RestartEtcdService),
		ForceApiserverRestart:         strconv.FormatBool(options.RestartKubernetesService),
		ForceControllerManagerRestart: strconv.FormatBool(options.RestartKubernetesService),
		ForceSchedulerRestart:         strconv.FormatBool(options.RestartKubernetesService),
		ForceProxyRestart:             strconv.FormatBool(options.RestartKubernetesService),
		ForceKubeletRestart:           strconv.FormatBool(options.RestartKubernetesService),
		ForceCalicoRestart:            strconv.FormatBool(options.RestartCalicoService),
		ForceDockerRestart:            strconv.FormatBool(options.RestartDockerService),
	}

	return &vars, nil
}

func (e *ansibleExecutor) Install(p *Plan, av *AnsibleVars) error {
	// build inventory
	inventory := buildNodeInventory(p)
	inventoryFile := filepath.Join(e.ansibleDir, "inventory.ini")
	err := ioutil.WriteFile(inventoryFile, inventory, 0644)
	if err != nil {
		return fmt.Errorf("error writing ansible inventory file: %v", err)
	}

	// run ansible
	playbook := filepath.Join(e.ansibleDir, "playbooks", "kubernetes.yaml")
	err = e.runAnsiblePlaybook(inventoryFile, playbook, *av)
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}

	return nil
}

func (e *ansibleExecutor) runAnsiblePlaybook(inventoryFile, playbookFile string, vars AnsibleVars) error {
	extraVars, err := vars.CommandLineVars()
	if err != nil {
		return fmt.Errorf("error getting vars: %v", err)
	}

	cmd := exec.Command(filepath.Join(e.ansibleDir, "bin", "ansible-playbook"), "-i", inventoryFile, "-s", playbookFile, "--extra-vars", extraVars)
	cmd.Stdout = e.out
	cmd.Stderr = e.errOut
	os.Setenv("PYTHONPATH", e.pythonPath)
	// Ansible config
	os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "False")

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running playbook: %v", err)
	}

	return nil
}

func buildNodeInventory(p *Plan) []byte {
	inv := &bytes.Buffer{}
	writeHostGroup(inv, "etcd", p.Etcd.Nodes, p.Cluster.SSH)
	writeHostGroup(inv, "master", p.Master.Nodes, p.Cluster.SSH)
	writeHostGroup(inv, "worker", p.Worker.Nodes, p.Cluster.SSH)

	return inv.Bytes()
}

func writeHostGroup(inv io.Writer, groupName string, nodes []Node, ssh SSHConfig) {
	fmt.Fprintf(inv, "[%s]\n", groupName)
	for _, n := range nodes {
		internalIP := n.IP
		if n.InternalIP != "" {
			internalIP = n.InternalIP
		}
		fmt.Fprintf(inv, "%s ansible_host=%s internal_ipv4=%s ansible_ssh_private_key_file=%s ansible_port=%d ansible_user=%s\n", n.Host, n.IP, internalIP, ssh.Key, ssh.Port, ssh.User)
	}
}

func getPythonPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working dir: %v", err)
	}
	lib := filepath.Join(wd, "ansible", "lib", "python2.7", "site-packages")
	lib64 := filepath.Join(wd, "ansible", "lib64", "python2.7", "site-packages")
	return fmt.Sprintf("%s:%s", lib, lib64), nil
}
