package install

import (
	"fmt"
	"io"
	"path/filepath"
	"strconv"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
)

// The Executor will carry out the installation plan
type Executor interface {
	Install(p *Plan) error
}

type ansibleExecutor struct {
	runner          ansible.Runner
	tlsDirectory    string
	restartServices bool
	ansibleStdout   io.Reader
	out             io.Writer
	explainer       Explainer
}

// NewExecutor returns an executor for performing installations according to the installation plan.
func NewExecutor(out io.Writer, errOut io.Writer, tlsDirectory string, restartServices, verbose bool, outputFormat string) (Executor, error) {
	// TODO: Is there a better way to handle this path to the ansible install dir?
	ansibleDir := "ansible"

	// configure ansible output
	var outFormat ansible.OutputFormat
	var explainer Explainer
	switch outputFormat {
	case "raw":
		outFormat = ansible.RawFormat
		explainer = &RawExplainer{out}
	case "simple":
		outFormat = ansible.JSONLinesFormat
		explainer = &AnsibleEventExplainer{ansible.EventStream, out, verbose}
	default:
		return nil, fmt.Errorf("Output format %q is not supported", outputFormat)
	}

	// Make ansible write to pipe, so that we can read on our end.
	r, w := io.Pipe()
	runner, err := ansible.NewRunner(w, errOut, ansibleDir, outFormat, verbose)
	if err != nil {
		return nil, fmt.Errorf("error creating ansible runner: %v", err)
	}

	td, err := filepath.Abs(tlsDirectory)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path from %q: %v", tlsDirectory, err)
	}

	return &ansibleExecutor{
		runner:          runner,
		tlsDirectory:    td,
		restartServices: restartServices,
		ansibleStdout:   r,
		out:             out,
		explainer:       explainer,
	}, nil
}

// Install the cluster according to the installation plan
func (ae *ansibleExecutor) Install(p *Plan) error {
	// Build the ansible inventory
	etcdNodes := []ansible.Node{}
	for _, n := range p.Etcd.Nodes {
		etcdNodes = append(etcdNodes, installNodeToAnsibleNode(&n, &p.Cluster.SSH))
	}
	masterNodes := []ansible.Node{}
	for _, n := range p.Master.Nodes {
		masterNodes = append(masterNodes, installNodeToAnsibleNode(&n, &p.Cluster.SSH))
	}
	workerNodes := []ansible.Node{}
	for _, n := range p.Worker.Nodes {
		workerNodes = append(workerNodes, installNodeToAnsibleNode(&n, &p.Cluster.SSH))
	}
	inventory := ansible.Inventory{
		{
			Name:  "etcd",
			Nodes: etcdNodes,
		},
		{
			Name:  "master",
			Nodes: masterNodes,
		},
		{
			Name:  "worker",
			Nodes: workerNodes,
		},
	}

	dnsIP, err := getDNSServiceIP(p)
	if err != nil {
		return fmt.Errorf("error getting DNS service IP: %v", err)
	}

	ev := ansible.ExtraVars{
		"kubernetes_cluster_name":   p.Cluster.Name,
		"kubernetes_admin_password": p.Cluster.AdminPassword,
		"tls_directory":             ae.tlsDirectory,
		"calico_network_type":       p.Cluster.Networking.Type,
		"kubernetes_services_cidr":  p.Cluster.Networking.ServiceCIDRBlock,
		"kubernetes_pods_cidr":      p.Cluster.Networking.PodCIDRBlock,
		"kubernetes_dns_service_ip": dnsIP,
	}

	if p.Cluster.LocalRepository != "" {
		ev["local_repoository_path"] = p.Cluster.LocalRepository
	}

	if ae.restartServices {
		services := []string{"etcd", "apiserver", "controller", "scheduler", "proxy", "kubelet", "calico_node", "docker"}
		for _, s := range services {
			ev[fmt.Sprintf("force_%s_restart", s)] = strconv.FormatBool(true)
		}
	}

	// Start explainer for handling ansible's stdout stream
	go ae.explainer.Explain(ae.ansibleStdout)

	// Run the installation playbook
	err = ae.runner.RunPlaybook(inventory, "kubernetes.yaml", ev)
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}
	return nil
}

// Converts plan node to ansible node
func installNodeToAnsibleNode(n *Node, s *SSHConfig) ansible.Node {
	return ansible.Node{
		Host:          n.Host,
		PublicIP:      n.IP,
		InternalIP:    n.InternalIP,
		SSHPrivateKey: s.Key,
		SSHUser:       s.User,
		SSHPort:       s.Port,
	}
}
