package install

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type NetworkConfig struct {
	Type             string
	PodCIDRBlock     string `yaml:"pod_cidr_block"`
	ServiceCIDRBlock string `yaml:"service_cidr_block"`
}

type CertsConfig struct {
	Expiry          string
	LocationCity    string `yaml:"location_city"`
	LocationState   string `yaml:"location_state"`
	LocationCountry string `yaml:"location_country"`
}

type SSHConfig struct {
	User string
	Key  string `yaml:"ssh_key"`
	Port int    `yaml:"ssh_port"`
}

type Cluster struct {
	Name         string
	Networking   NetworkConfig
	Certificates CertsConfig
	SSH          SSHConfig
}

type Node struct {
	Host   string
	IP     string
	Labels []string
}

type NodeGroup struct {
	ExpectedCount int `yaml:"expected_count"`
	Nodes         []Node
}

type MasterNodeGroup struct {
	NodeGroup
	LoadBalancedFQDN      string `yaml:"load_balanced_fqdn"`
	LoadBalancedShortName string `yaml:"load_balanced_short_name"`
}

type Plan struct {
	Cluster Cluster
	Etcd    NodeGroup
	Master  MasterNodeGroup
	Worker  NodeGroup
}

type PlanReaderWriter interface {
	Read() (*Plan, error)
	Write(*Plan) error
	Exists() bool
}

type PlanFile struct {
	File string
}

func (pf *PlanFile) Read() (*Plan, error) {
	d, err := ioutil.ReadFile(pf.File)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %v", err)
	}

	p := &Plan{}
	if err = yaml.Unmarshal(d, p); err != nil {
		return nil, fmt.Errorf("failed to unmarshal plan: %v", err)
	}
	return p, nil
}

func (pf *PlanFile) Write(p *Plan) error {
	d, err := yaml.Marshal(&p)
	if err != nil {
		return fmt.Errorf("error marshalling install plan template: %v", err)
	}

	err = ioutil.WriteFile(pf.File, d, 0644)
	if err != nil {
		return fmt.Errorf("error writing install plan template: %v", err)
	}
	return nil
}

func (pf *PlanFile) Exists() bool {
	_, err := os.Stat(pf.File)
	return !os.IsNotExist(err)
}

func WritePlanTemplate(p Plan, w PlanReaderWriter) error {
	// Set sensible defaults
	p.Cluster.Name = "kubernetes"

	// Set SSH defaults
	p.Cluster.SSH.User = "kismaticuser"
	p.Cluster.SSH.Key = "kismaticuser.key"
	p.Cluster.SSH.Port = 22

	// Set Networking defaults
	p.Cluster.Networking.Type = "bridged"
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

func ValidatePlan(p *Plan) error {
	if p.Cluster.SSH.Key == "" {
		return fmt.Errorf("SshKeyPath is empty. The path to the SSH key is required.")
	}

	return nil
}
