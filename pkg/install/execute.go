package install

import (
	"fmt"
	"io"
	"io/ioutil"
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

// The PreFlightExecutor will run pre-flight checks against the
// environment defined in the plan file
type PreFlightExecutor interface {
	RunPreFlightCheck(*Plan) error
}

// The Executor will carry out the installation plan
type Executor interface {
	PreFlightExecutor
	Install(p *Plan) error
	RunSmokeTest(*Plan) error
}

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

// NewExecutor returns an executor for performing installations according to the installation plan.
func NewExecutor(stdout io.Writer, errOut io.Writer, options ExecutorOptions) (Executor, error) {
	// TODO: Is there a better way to handle this path to the ansible install dir?
	ansibleDir := "ansible"

	// Validate options
	if options.CASigningRequest == "" {
		return nil, fmt.Errorf("CASigningRequest option cannot be empty")
	}
	if options.CAConfigFile == "" {
		return nil, fmt.Errorf("CAConfigFile option cannot be empty")
	}
	if options.CASigningProfile == "" {
		return nil, fmt.Errorf("CASigningProfile option cannot be empty")
	}
	if options.GeneratedAssetsDirectory == "" {
		return nil, fmt.Errorf("GeneratedAssetsDirectory option cannot be empty")
	}
	if options.RunsDirectory == "" {
		options.RunsDirectory = "./runs"
	}

	// Setup the console output format
	var outFormat ansible.OutputFormat
	switch options.OutputFormat {
	case "raw":
		outFormat = ansible.RawFormat
	case "simple":
		outFormat = ansible.JSONLinesFormat
	default:
		return nil, fmt.Errorf("Output format %q is not supported", options.OutputFormat)
	}

	return &ansibleExecutor{
		options:             options,
		stdout:              stdout,
		consoleOutputFormat: outFormat,
		ansibleDir:          ansibleDir,
	}, nil
}

func NewPreFlightExecutor(stdout io.Writer, errOut io.Writer, options ExecutorOptions) (PreFlightExecutor, error) {
	ansibleDir := "ansible"
	if options.RunsDirectory == "" {
		options.RunsDirectory = "./runs"
	}
	// Setup the console output format
	var outFormat ansible.OutputFormat
	switch options.OutputFormat {
	case "raw":
		outFormat = ansible.RawFormat
	case "simple":
		outFormat = ansible.JSONLinesFormat
	default:
		return nil, fmt.Errorf("Output format %q is not supported", options.OutputFormat)
	}

	return &ansibleExecutor{
		options:             options,
		stdout:              stdout,
		consoleOutputFormat: outFormat,
		ansibleDir:          ansibleDir,
	}, nil
}

type ansibleExecutor struct {
	options             ExecutorOptions
	stdout              io.Writer
	consoleOutputFormat ansible.OutputFormat
	ansibleDir          string
}

// Install the cluster according to the installation plan
func (ae *ansibleExecutor) Install(p *Plan) error {
	runDirectory, err := ae.createRunDirectory("install")
	if err != nil {
		return fmt.Errorf("error creating working directory for installation: %v", err)
	}

	// Save the plan file that was used for this execution
	fp := FilePlanner{
		File: filepath.Join(runDirectory, "kismatic-cluster.yaml"),
	}
	if err = fp.Write(p); err != nil {
		return fmt.Errorf("error recording plan file to %s: %v", fp.File, err)
	}

	// Generate private keys and certificates for the cluster
	generatedCertsDir, err := ae.generateTLSAssets(p)
	if err != nil {
		return err
	}

	// Build the ansible inventory
	inventory := buildInventoryFromPlan(p)

	dnsIP, err := getDNSServiceIP(p)
	if err != nil {
		return fmt.Errorf("error getting DNS service IP: %v", err)
	}

	// Need absolute path for ansible. Otherwise ansible looks for it in the wrong place.
	tlsDir, err := filepath.Abs(generatedCertsDir)
	if err != nil {
		return fmt.Errorf("failed to determine absolute path to %s: %v", generatedCertsDir, err)
	}
	ev := ansible.ExtraVars{
		"kubernetes_cluster_name":    p.Cluster.Name,
		"kubernetes_admin_password":  p.Cluster.AdminPassword,
		"tls_directory":              tlsDir,
		"calico_network_type":        p.Cluster.Networking.Type,
		"kubernetes_services_cidr":   p.Cluster.Networking.ServiceCIDRBlock,
		"kubernetes_pods_cidr":       p.Cluster.Networking.PodCIDRBlock,
		"kubernetes_dns_service_ip":  dnsIP,
		"modify_hosts_file":          strconv.FormatBool(p.Cluster.Networking.UpdateHostsFiles),
		"enable_calico_policy":       strconv.FormatBool(p.Cluster.Networking.PolicyEnabled),
		"enable_docker_registry":     strconv.FormatBool(p.DockerRegistry.UseInternal),
		"allow_package_installation": strconv.FormatBool(p.Cluster.AllowPackageInstallation),
	}

	if ae.options.RestartServices {
		services := []string{"etcd", "apiserver", "controller_manager", "scheduler", "proxy", "kubelet", "calico_node", "docker"}
		for _, s := range services {
			ev[fmt.Sprintf("force_%s_restart", s)] = strconv.FormatBool(true)
		}
	}

	ansibleLogFilename := filepath.Join(runDirectory, "ansible.log")
	ansibleLogFile, err := os.Create(ansibleLogFilename)
	if err != nil {
		return fmt.Errorf("error creating ansible log file %q: %v", ansibleLogFilename, err)
	}

	// Run the installation playbook
	util.PrintHeader(ae.stdout, "Installing Cluster", '=')
	playbook := "kubernetes.yaml"
	eventExplainer := &explain.DefaultEventExplainer{}
	if err = ae.runPlaybookWithExplainer(playbook, eventExplainer, inventory, ev, ansibleLogFile); err != nil {
		return err
	}

	return nil
}

func (ae *ansibleExecutor) RunSmokeTest(p *Plan) error {
	runDirectory, err := ae.createRunDirectory("smoketest")
	if err != nil {
		return fmt.Errorf("error creating working directory for smoke test: %v", err)
	}

	ansibleLogFilename := filepath.Join(runDirectory, "ansible.log")
	ansibleLogFile, err := os.Create(ansibleLogFilename)
	if err != nil {
		return fmt.Errorf("error creating ansible log file %q: %v", ansibleLogFilename, err)
	}

	ev := ansible.ExtraVars{
		"kuberang_path":     filepath.Join("kuberang", "linux", "amd64", "kuberang"),
		"modify_hosts_file": strconv.FormatBool(p.Cluster.Networking.UpdateHostsFiles),
	}
	inventory := buildInventoryFromPlan(p)

	// run the preflight playbook with preflight explainer
	util.PrintHeader(ae.stdout, "Running Smoke Test", '=')
	playbook := "smoketest.yaml"
	explainer := &explain.PreflightEventExplainer{
		DefaultExplainer: &explain.DefaultEventExplainer{},
	}
	if err = ae.runPlaybookWithExplainer(playbook, explainer, inventory, ev, ansibleLogFile); err != nil {
		return fmt.Errorf("error running smoketest: %v", err)
	}
	return nil
}

// RunPreflightCheck against the nodes defined in the plan
func (ae *ansibleExecutor) RunPreFlightCheck(p *Plan) error {
	runDirectory, err := ae.createRunDirectory("preflight")
	if err != nil {
		return fmt.Errorf("error creating working directory for preflight: %v", err)
	}
	// Save the plan file that was used for this execution
	fp := FilePlanner{
		File: filepath.Join(runDirectory, "kismatic-cluster.yaml"),
	}
	if err = fp.Write(p); err != nil {
		return fmt.Errorf("error recording plan file to %s: %v", fp.File, err)
	}

	ansibleLogFilename := filepath.Join(runDirectory, "ansible.log")
	ansibleLogFile, err := os.Create(ansibleLogFilename)
	if err != nil {
		return fmt.Errorf("error creating ansible log file %q: %v", ansibleLogFilename, err)
	}

	// Build inventory and save it in runs directory
	inventory := buildInventoryFromPlan(p)

	pwd, _ := os.Getwd()
	ev := ansible.ExtraVars{
		// TODO: attempt to clean up these paths somehow...
		"kismatic_preflight_checker":       filepath.Join("inspector", "linux", "amd64", "kismatic-inspector"),
		"kismatic_preflight_checker_local": filepath.Join(pwd, "ansible", "playbooks", "inspector", runtime.GOOS, runtime.GOARCH, "kismatic-inspector"),
		"modify_hosts_file":                strconv.FormatBool(p.Cluster.Networking.UpdateHostsFiles),
		"allow_package_installation":       strconv.FormatBool(p.Cluster.AllowPackageInstallation),
	}

	// run the pre-flight playbook with pre-flight explainer
	playbook := "preflight.yaml"
	explainer := &explain.PreflightEventExplainer{
		DefaultExplainer: &explain.DefaultEventExplainer{},
	}
	if err = ae.runPlaybookWithExplainer(playbook, explainer, inventory, ev, ansibleLogFile); err != nil {
		return fmt.Errorf("error running preflight: %v", err)
	}
	return nil
}

func (ae *ansibleExecutor) createRunDirectory(runName string) (string, error) {
	start := time.Now()
	runDirectory := filepath.Join(ae.options.RunsDirectory, runName, start.Format("2006-01-02-15-04-05"))
	if err := os.MkdirAll(runDirectory, 0777); err != nil {
		return "", fmt.Errorf("error creating directory: %v", err)
	}
	return runDirectory, nil
}

func (ae *ansibleExecutor) generateTLSAssets(p *Plan) (certsDir string, err error) {
	keysDir := filepath.Join(ae.options.GeneratedAssetsDirectory, "keys")
	if err = os.MkdirAll(keysDir, 0777); err != nil {
		return "", fmt.Errorf("error creating directory %s for storing TLS assets: %v", keysDir, err)
	}
	pki := LocalPKI{
		CACsr:                   ae.options.CASigningRequest,
		CAConfigFile:            ae.options.CAConfigFile,
		CASigningProfile:        ae.options.CASigningProfile,
		GeneratedCertsDirectory: filepath.Join(ae.options.GeneratedAssetsDirectory, "keys"),
		Log: ae.stdout,
	}

	// Generate or read cluster Certificate Authority
	util.PrintHeader(ae.stdout, "Configuring Certificates", '=')
	var ca *tls.CA
	if !ae.options.SkipCAGeneration {
		util.PrettyPrintOk(ae.stdout, "Generating cluster Certificate Authority")
		ca, err = pki.GenerateClusterCA(p)
		if err != nil {
			return "", fmt.Errorf("error generating CA for the cluster: %v", err)
		}
	} else {
		util.PrettyPrintOk(ae.stdout, "Skipping Certificate Authority generation\n")
		ca, err = pki.ReadClusterCA(p)
		if err != nil {
			return "", fmt.Errorf("error reading cluster CA: %v", err)
		}
	}

	// Generate node and user certificates
	err = pki.GenerateClusterCerts(p, ca, []string{"admin"})
	if err != nil {
		return "", fmt.Errorf("error generating certificates for the cluster: %v", err)
	}
	util.PrettyPrintOk(ae.stdout, "Generated cluster certificates at %q", pki.GeneratedCertsDirectory)
	return keysDir, nil
}

func (ae *ansibleExecutor) runPlaybookWithExplainer(playbook string, explainer explain.AnsibleEventExplainer, inv ansible.Inventory, ev ansible.ExtraVars, ansibleLog io.Writer) error {
	// Setup sinks for explainer and ansible stdout
	var explainerOut, ansibleOut io.Writer
	switch ae.consoleOutputFormat {
	case ansible.JSONLinesFormat:
		explainerOut = ae.stdout
		ansibleOut = ansibleLog
	case ansible.RawFormat:
		explainerOut = ioutil.Discard
		ansibleOut = io.MultiWriter(ae.stdout, ansibleLog)
	}

	// Send stdout and stderr to ansibleOut
	runner, err := ansible.NewRunner(ansibleOut, ansibleOut, ae.ansibleDir)
	if err != nil {
		return fmt.Errorf("error creating ansible runner: %v", err)
	}

	// Start running ansible with the given playbook
	eventStream, err := runner.StartPlaybook(playbook, inv, ev)
	if err != nil {
		return fmt.Errorf("error running ansible playbook: %v", err)
	}

	// Setup explainer for handling event stream. Ansible will block
	// if the stream is not consumed
	streamExplainer := &explain.AnsibleEventStreamExplainer{
		Out:            explainerOut,
		Verbose:        ae.options.Verbose,
		EventExplainer: explainer,
	}
	go streamExplainer.Explain(eventStream)

	// Wait until ansible exits
	if err = runner.WaitPlaybook(); err != nil {
		return fmt.Errorf("error running playbook: %v", err)
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
