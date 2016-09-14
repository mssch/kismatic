package install

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apprenda/kismatic-platform/pkg/util"

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

// Write the plan to the file system
func (fp *FilePlanner) Write(p *Plan) error {
	d, err := yaml.Marshal(&p)
	if err != nil {
		return fmt.Errorf("error marshalling install plan template: %v", err)
	}

	err = ioutil.WriteFile(fp.File, d, 0644)
	if err != nil {
		return fmt.Errorf("error writing install plan template: %v", err)
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

	// Set SSH defaults
	p.Cluster.SSH.User = "kismaticuser"
	p.Cluster.SSH.Key = "kismaticuser.key"
	p.Cluster.SSH.Port = 22

	// Set Networking defaults
	p.Cluster.Networking.Type = "overlay"
	p.Cluster.Networking.PodCIDRBlock = "172.16.0.0/16"
	p.Cluster.Networking.ServiceCIDRBlock = "172.17.0.0/16"

	// Set Certificate defaults
	p.Cluster.Certificates.Expiry = "17520h"
	p.Cluster.Certificates.LocationCity = "Troy"
	p.Cluster.Certificates.LocationState = "New York"
	p.Cluster.Certificates.LocationCountry = "US"

	// Generate entries for all node types
	n := Node{Host: "shortname", IP: "127.0.0.1"}
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
