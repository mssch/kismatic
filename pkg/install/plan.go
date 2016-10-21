package install

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/apprenda/kismatic-platform/pkg/util"
	garbler "github.com/michaelbironneau/garbler/lib"

	yaml "gopkg.in/yaml.v2"
)

// PlanReadWriter is capable of reading/writing a Plan
type PlanReadWriter interface {
	Read() (*Plan, error)
	Write(*Plan) error
}

// Planner is used to plan the installation
type Planner interface {
	PlanReadWriter
	PlanExists() bool
}

// FilePlanner is a file-based installation planner
type FilePlanner struct {
	File string
}

// Read the plan from the file system
func (fp *FilePlanner) Read() (*Plan, error) {
	d, err := ioutil.ReadFile(fp.File)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	p := &Plan{}
	if err = yaml.Unmarshal(d, p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %v", err)
	}
	return p, nil
}

var yamlKeyRE = regexp.MustCompile(`[^a-zA-Z]*([a-z_\-A-Z]+)[ ]*:`)

// Write the plan to the file system
func (fp *FilePlanner) Write(p *Plan) error {
	oneTimeComments := commentMap
	bytez, marshalErr := yaml.Marshal(p)
	if marshalErr != nil {
		return fmt.Errorf("error marshalling plan to yaml", marshalErr)
	}

	f, err := os.Create(fp.File)
	if err != nil {
		return fmt.Errorf("error making plan file", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(bytes.NewReader(bytez))
	for scanner.Scan() {
		text := scanner.Text()
		matched := yamlKeyRE.FindStringSubmatch(text)

		if matched != nil && len(matched) > 1 {
			if thiscomment, ok := oneTimeComments[matched[1]]; ok {
				f.WriteString(fmt.Sprintf("%-40s # %s\n", text, thiscomment))
				delete(oneTimeComments, matched[1])
				continue
			}
		}
		f.WriteString(text + "\n")
	}

	return nil
}

// PlanExists return true if the plan exists on the file system
func (fp *FilePlanner) PlanExists() bool {
	_, err := os.Stat(fp.File)
	return !os.IsNotExist(err)
}

// WritePlanTemplate writes an installation plan with pre-filled defaults.
func WritePlanTemplate(p Plan, w PlanReadWriter) error {
	// Set sensible defaults
	p.Cluster.Name = "kubernetes"
	generatedAdminPass, _ := garbler.NewPassword(nil)
	// Cannot contain ":" in kubecofig
	generatedAdminPass = strings.Replace(generatedAdminPass, ":", "new", -1)
	p.Cluster.AdminPassword = generatedAdminPass

	// Set SSH defaults
	p.Cluster.SSH.User = "kismaticuser"
	p.Cluster.SSH.Key = "kismaticuser.key"
	p.Cluster.SSH.Port = 22

	// Set Networking defaults
	p.Cluster.Networking.Type = "overlay"
	p.Cluster.Networking.PodCIDRBlock = "172.16.0.0/16"
	p.Cluster.Networking.ServiceCIDRBlock = "172.17.0.0/16"
	p.Cluster.Networking.UpdateHostsFiles = false
	p.Cluster.Networking.PolicyEnabled = false

	// Set Certificate defaults
	p.Cluster.Certificates.Expiry = "17520h"
	p.DockerRegistry.SetupInternal = false

	// Set DockerRegistry defaults
	p.DockerRegistry.Port = 443

	// Generate entries for all node types
	n := Node{}
	for i := 0; i < p.Etcd.ExpectedCount; i++ {
		p.Etcd.Nodes = append(p.Etcd.Nodes, n)
	}

	for i := 0; i < p.Master.ExpectedCount; i++ {
		p.Master.Nodes = append(p.Master.Nodes, n)
	}

	for i := 0; i < p.Worker.ExpectedCount; i++ {
		p.Worker.Nodes = append(p.Worker.Nodes, n)
	}

	if err := w.Write(&p); err != nil {
		return fmt.Errorf("error writing installation plan template: %v", err)
	}
	return nil
}

func getKubernetesServiceIP(p *Plan) (string, error) {
	ip, err := util.GetIPFromCIDR(p.Cluster.Networking.ServiceCIDRBlock, 1)
	if err != nil {
		return "", fmt.Errorf("error getting kubernetes service IP: %v", err)
	}
	return ip.To4().String(), nil
}

func getDNSServiceIP(p *Plan) (string, error) {
	ip, err := util.GetIPFromCIDR(p.Cluster.Networking.ServiceCIDRBlock, 2)
	if err != nil {
		return "", fmt.Errorf("error getting DNS service IP: %v", err)
	}
	return ip.To4().String(), nil
}

var commentMap = map[string]string{
	"admin_password":           "This token is used for Kubernetes' administration.",
	"type":                     "Overlay or Routed. Routed pods can be addressed from outside the Kubernetes cluster; Overlay pods can only address each other.",
	"pod_cidr_block":           "Kubernetes will assign pods IPs in this range. Do not use a range that is already in use on your local network!",
	"service_cidr_block":       "Kubernetes will assign services IPs in this range. Do not use a range that is already in use by your local network or pod network!",
	"policy_enabled":           "When true, enables network policy enforcement on the Kubernetes Pod network. This is an advanced feature.",
	"update_hosts_files":       "When true, the installer will add entries for all nodes to other nodes' hosts files. Use when you don't have access to DNS.",
	"expiry":                   "Self-signed certificate expiration period in hours; default is 2 years.",
	"ssh_key":                  "Absolute path to the ssh public key we should use to manage nodes.",
	"etcd":                     "Here you will identify all of the nodes that should play the etcd role on your cluster.",
	"master":                   "Here you will identify all of the nodes that should play the master role.",
	"worker":                   "Here you will identify all of the nodes that will be workers.",
	"host":                     "The (short) hostname of a node, e.g. etcd01",
	"ip":                       "The ip address the installer should use to manage this node, e.g. 8.8.8.8.",
	"internalip":               "If the node has an IP for internal traffic, enter it here; otherwise leave blank.",
	"load_balanced_fqdn":       "If using a load balanced master fqdn enter it here; otherwise use the fqdn of any master node.",
	"load_balanced_short_name": "If using a load balanced master enter its short name here; otherwise use the short name of any master node.",
	"docker_registry":          "Here you will provide the details of your Docker registry or setup an internal one to run in the cluster. This is optional and the cluster will always have access to the Docker Hub.",
	"setup_internal":           "When true, a Docker Registry will be installed on top of your cluster and used to host Docker images needed for its installation.",
	"address":                  "IP or hostname for your Docker registry. An internal registry will NOT be setup when this field is provided. Must be accessible from all the nodes in the cluster.",
	"port":                     "Port for your Docker registry.",
	"CA":                       "Absolute path to the CA that was used when starting your Docker registry. The docker daemons on all nodes in the cluster will be configured with this CA.",
}
