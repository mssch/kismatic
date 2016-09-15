package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
	"github.com/apprenda/kismatic-platform/pkg/install/explain"
	"github.com/apprenda/kismatic-platform/pkg/tls"
	"github.com/apprenda/kismatic-platform/pkg/util"
)

// ExecutorOptions are used to configure the executor
type ExecutorOptions struct {
	// CASigningRequest in JSON format expected by cfSSL
	CASigningRequest string
	// CAConfigFile is the Certificate Authority configuration file
	// in the JSON format expected by cfSSL
	CAConfigFile string
	// CASigningProfile is the signing profile to be used when signing
	// certificates. The profile must be defined in the CAConfigFile
	CASigningProfile string
	// SkipCAGeneration determines whether the Certificate Authority should
	// be generated. If false, an existing CA file must exist.
	SkipCAGeneration bool
	// GeneratedAssetsDirectory is the location where generated assets
	// are to be stored
	GeneratedAssetsDirectory string
	// RestartServices determines whether the cluster services should be
	// restarted during the installation.
	RestartServices bool
	// OutputFormat sets the format of the executor
	OutputFormat string
	// Verbose output from the executor
	Verbose bool
	// RunsDirectory is where information about installation runs is kept
	RunsDirectory string
}

// The Executor will carry out the installation plan
type Executor interface {
	Install(p *Plan) error
	RunPreflightCheck(*Plan) error
}

type ansibleExecutor struct {
	options       ExecutorOptions
	runner        ansible.Runner
	ansibleStdout io.Reader
	out           io.Writer
	outputFormat  ansible.OutputFormat
}

// NewExecutor returns an executor for performing installations according to the installation plan.
func NewExecutor(out io.Writer, errOut io.Writer, options ExecutorOptions) (Executor, error) {
	// TODO: Is there a better way to handle this path to the ansible install dir?
	ansibleDir := "ansible"

	// TODO: Validate options here
	if options.RunsDirectory == "" {
		options.RunsDirectory = "./runs"
	}

	// configure ansible output
	var outFormat ansible.OutputFormat
	switch options.OutputFormat {
	case "raw":
		outFormat = ansible.RawFormat
	case "simple":
		outFormat = ansible.JSONLinesFormat
	default:
		return nil, fmt.Errorf("Output format %q is not supported", options.OutputFormat)
	}

	// Make ansible write to pipe, so that we can read on our end.
	r, w := io.Pipe()
	runner, err := ansible.NewRunner(w, errOut, ansibleDir)
	if err != nil {
		return nil, fmt.Errorf("error creating ansible runner: %v", err)
	}

	return &ansibleExecutor{
		options:       options,
		runner:        runner,
		ansibleStdout: r,
		out:           out,
		outputFormat:  outFormat,
	}, nil
}

// Install the cluster according to the installation plan
func (ae *ansibleExecutor) Install(p *Plan) error {
	start := time.Now()
	runDirectory := filepath.Join(ae.options.RunsDirectory, start.Format("20060102030405"))
	if err := os.MkdirAll(runDirectory, 0777); err != nil {
		return fmt.Errorf("error creating working directory for installation: %v", err)
	}

	// Save the plan file that was used for this execution
	fp := FilePlanner{
		File: filepath.Join(runDirectory, "kismatic-cluster.yaml"),
	}
	if err := fp.Write(p); err != nil {
		return fmt.Errorf("error recording plan file to %s: %v", fp.File, err)
	}

	// Generate cluster TLS assets
	keysDir := filepath.Join(ae.options.GeneratedAssetsDirectory, "keys")
	if err := os.MkdirAll(keysDir, 0777); err != nil {
		return fmt.Errorf("error creating directory %s for storing TLS assets: %v", keysDir, err)
	}
	pki := LocalPKI{
		CACsr:                   ae.options.CASigningRequest,
		CAConfigFile:            ae.options.CAConfigFile,
		CASigningProfile:        ae.options.CASigningProfile,
		GeneratedCertsDirectory: filepath.Join(ae.options.GeneratedAssetsDirectory, "keys"),
		Log: ae.out,
	}

	// Generate or read cluster Certificate Authority
	util.PrintHeader(ae.out, "Configuring Certificates")
	var ca *tls.CA
	var err error
	if !ae.options.SkipCAGeneration {
		util.PrettyPrintOk(ae.out, "Generating cluster Certificate Authority")
		ca, err = pki.GenerateClusterCA(p)
		if err != nil {
			return fmt.Errorf("error generating CA for the cluster: %v", err)
		}
	} else {
		util.PrettyPrint(ae.out, "Skipping Certificate Authority generation\n")
		ca, err = pki.ReadClusterCA(p)
		if err != nil {
			return fmt.Errorf("error reading cluster CA: %v", err)
		}
	}

	// Generate node and user certificates
	err = pki.GenerateClusterCerts(p, ca, []string{"admin"})
	if err != nil {
		return fmt.Errorf("error generating certificates for the cluster: %v", err)
	}
	util.PrettyPrintOkf(ae.out, "Generated cluster certificates at %q", pki.GeneratedCertsDirectory)

	// Build the ansible inventory
	inventory := buildInventoryFromPlan(p)

	dnsIP, err := getDNSServiceIP(p)
	if err != nil {
		return fmt.Errorf("error getting DNS service IP: %v", err)
	}

	// Need absolute path for ansible. Otherwise it looks in the wrong place.
	tlsDir, err := filepath.Abs(pki.GeneratedCertsDirectory)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path to %s: %v", pki.GeneratedCertsDirectory, err)
	}
	ev := ansible.ExtraVars{
		"kubernetes_cluster_name":   p.Cluster.Name,
		"kubernetes_admin_password": p.Cluster.AdminPassword,
		"tls_directory":             tlsDir,
		"calico_network_type":       p.Cluster.Networking.Type,
		"kubernetes_services_cidr":  p.Cluster.Networking.ServiceCIDRBlock,
		"kubernetes_pods_cidr":      p.Cluster.Networking.PodCIDRBlock,
		"kubernetes_dns_service_ip": dnsIP,
		"modify_hosts_file":         strconv.FormatBool(p.Cluster.HostsFileDNS),
	}

	if p.Cluster.LocalRepository != "" {
		ev["local_repoository_path"] = p.Cluster.LocalRepository
	}

	if ae.options.RestartServices {
		services := []string{"etcd", "apiserver", "controller", "scheduler", "proxy", "kubelet", "calico_node", "docker"}
		for _, s := range services {
			ev[fmt.Sprintf("force_%s_restart", s)] = strconv.FormatBool(true)
		}
	}

	// Start explainer for handling ansible's stdout stream
	var exp explain.StreamExplainer
	switch ae.outputFormat {
	case ansible.RawFormat:
		exp = &explain.RawExplainer{ae.out}
	case ansible.JSONLinesFormat:
		exp = &explain.AnsibleEventStreamExplainer{
			EventStream:    ansible.EventStream,
			Out:            ae.out,
			Verbose:        ae.options.Verbose,
			EventExplainer: &explain.DefaultEventExplainer{},
		}
	}
	go exp.Explain(ae.ansibleStdout)

	// Run the installation playbook
	err = ae.runner.RunPlaybook(inventory, "kubernetes.yaml", ev, ae.outputFormat, ae.options.Verbose)
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}
	return nil
}

func (ae *ansibleExecutor) RunPreflightCheck(p *Plan) error {
	// build inventory
	inventory := buildInventoryFromPlan(p)

	ev := ansible.ExtraVars{
		// TODO: attempt to clean up these paths somehow...
		"kismatic_preflight_checker":       filepath.Join("inspector", "linux", "amd64", "kismatic-inspector"),
		"kismatic_preflight_checker_local": filepath.Join("ansible", "playbooks", "inspector", runtime.GOOS, runtime.GOARCH, "kismatic-inspector"),
		"modify_hosts_file":                strconv.FormatBool(p.Cluster.HostsFileDNS),
	}

	// Set explainer for pre-flight checks
	var exp explain.StreamExplainer
	switch ae.outputFormat {
	case ansible.RawFormat:
		exp = &explain.RawExplainer{ae.out}
	case ansible.JSONLinesFormat:
		exp = &explain.AnsibleEventStreamExplainer{
			EventStream:    ansible.EventStream,
			Out:            ae.out,
			Verbose:        ae.options.Verbose,
			EventExplainer: &explain.PreflightEventExplainer{&explain.DefaultEventExplainer{}},
		}
	}
	go exp.Explain(ae.ansibleStdout)

	// run pre-flight playbook
	playbook := "preflight.yaml"
	err := ae.runner.RunPlaybook(inventory, playbook, ev, ae.outputFormat, ae.options.Verbose)
	if err != nil {
		return fmt.Errorf("error running pre-flight checks: %v", err)
	}
	return nil
}

func buildInventoryFromPlan(p *Plan) ansible.Inventory {
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

	return inventory
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
