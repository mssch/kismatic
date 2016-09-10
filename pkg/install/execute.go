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
}

// NewExecutor returns an executor for performing installations according to the installation plan.
func NewExecutor(out io.Writer, errOut io.Writer, tlsDirectory string, restartServices bool) (Executor, error) {
	// TODO: Is there a better way to handle this path to the ansible install dir?
	ansibleDir := "ansible"

	// Send runner output to the pipe
	r, w := io.Pipe()
	runner, err := ansible.NewRunner(w, errOut, ansibleDir)
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

	// Setup event handler for the install
	events := ansible.EventStream(ae.ansibleStdout)
	go func(es <-chan ansible.Event) {
		for e := range es {
			switch event := e.(type) {
			default:
				fmt.Fprintf(ae.out, "Unhandled event: %T\n", event)
			case *ansible.PlaybookStartEvent:
				fmt.Fprintf(ae.out, "Running playbook %s\n", event.Name)
			case *ansible.PlayStartEvent:
				fmt.Fprintf(ae.out, "- %s\n", event.Name)
			case *ansible.RunnerItemRetryEvent:
				fmt.Fprintf(ae.out, "[RETRYING] %s\n", event.Host)
			case *ansible.RunnerUnreachableEvent:
				fmt.Fprintf(ae.out, "[UNREACHABLE] %s\n", event.Host)
			case *ansible.RunnerFailedEvent:
				fmt.Fprintf(ae.out, "[ERROR] %s\n", event.Host)
				fmt.Fprintf(ae.out, "|- stdout: %s\n", event.Result.Stdout)
				fmt.Fprintf(ae.out, "|- stderr: %s\n", event.Result.Stderr)
			// Do nothing with the following events
			case *ansible.TaskStartEvent:
				continue
			case *ansible.HandlerTaskStartEvent:
				continue
			case *ansible.RunnerItemOKEvent:
				continue
			case *ansible.RunnerSkippedEvent:
				continue
			case *ansible.RunnerOKEvent:
				continue
			}
		}
	}(events)

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
