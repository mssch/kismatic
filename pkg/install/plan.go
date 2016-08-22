package install

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"time"

	yaml "gopkg.in/yaml.v2"
)

// NetworkConfig describes the cluster's networking configuration
type NetworkConfig struct {
	Type             string
	PodCIDRBlock     string `yaml:"pod_cidr_block"`
	ServiceCIDRBlock string `yaml:"service_cidr_block"`
}

// CertsConfig describes the cluster's trust and certificate configuration
type CertsConfig struct {
	Expiry          string
	LocationCity    string `yaml:"location_city"`
	LocationState   string `yaml:"location_state"`
	LocationCountry string `yaml:"location_country"`
}

// SSHConfig describes the cluster's SSH configuration for accessing nodes
type SSHConfig struct {
	User string
	Key  string `yaml:"ssh_key"`
	Port int    `yaml:"ssh_port"`
}

// Cluster describes a Kismatic cluster
type Cluster struct {
	Name         string
	Networking   NetworkConfig
	Certificates CertsConfig
	SSH          SSHConfig
}

// A Node is a compute unit, virtual or physical, that is part of the cluster
type Node struct {
	Host   string
	IP     string
	Labels []string
}

// A NodeGroup is a collection of nodes
type NodeGroup struct {
	ExpectedCount int `yaml:"expected_count"`
	Nodes         []Node
}

// MasterNodeGroup is the collection of master nodes
type MasterNodeGroup struct {
	NodeGroup
	LoadBalancedFQDN      string `yaml:"load_balanced_fqdn"`
	LoadBalancedShortName string `yaml:"load_balanced_short_name"`
}

// Plan is the installation plan that the user intends to execute
type Plan struct {
	Cluster Cluster
	Etcd    NodeGroup
	Master  MasterNodeGroup
	Worker  NodeGroup
}

// PlanReaderWriter is capable of reading/writing a Plan
type PlanReaderWriter interface {
	Read() (*Plan, error)
	Write(*Plan) error
	Exists() bool
}

// PlanFile is an installation plan backed by a file
type PlanFile struct {
	File string
}

// Read the plan from the file system
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

// Write the plan to the file system
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

// Exists return true if the plan exists on the file system
func (pf *PlanFile) Exists() bool {
	_, err := os.Stat(pf.File)
	return !os.IsNotExist(err)
}

// WritePlanTemplate writes an installation plan with pre-filled defaults.
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

// ValidatePlan runs validation against the installation plan to ensure
// that the plan contains valid user input. Returns true, nil if the validation
// is successful. Otherwise, returns false and a collection of validation errors.
func ValidatePlan(p *Plan) (bool, []error) {
	v := newValidator()
	v.validate(p)
	return v.valid()
}

type validatable interface {
	validate() (bool, []error)
}

type validator struct {
	errs []error
}

func newValidator() *validator {
	return &validator{
		errs: []error{},
	}
}

func (v *validator) addError(err ...error) {
	v.errs = append(v.errs, err...)
}

func (v *validator) validate(obj validatable) {
	if ok, err := obj.validate(); !ok {
		v.addError(err...)
	}
}

func (v *validator) validateWithErrPrefix(prefix string, obj validatable) {
	if ok, err := obj.validate(); !ok {
		newErrs := make([]error, len(err), len(err))
		for i, err := range err {
			newErrs[i] = fmt.Errorf("%s: %v", prefix, err)
		}
		v.addError(newErrs...)
	}
}

func (v *validator) valid() (bool, []error) {
	if len(v.errs) > 0 {
		return false, v.errs
	}
	return true, nil
}

func (p *Plan) validate() (bool, []error) {
	v := newValidator()

	v.validate(&p.Cluster)
	v.validateWithErrPrefix("Etcd nodes", &p.Etcd)
	v.validateWithErrPrefix("Master nodes", &p.Master)
	v.validateWithErrPrefix("Worker nodes", &p.Worker)

	return v.valid()
}

func (c *Cluster) validate() (bool, []error) {
	v := newValidator()
	if c.Name == "" {
		v.addError(errors.New("Cluster name cannot be empty"))
	}
	v.validate(&c.Networking)
	v.validate(&c.Certificates)
	v.validate(&c.SSH)

	return v.valid()
}

func (n *NetworkConfig) validate() (bool, []error) {
	v := newValidator()
	if n.Type == "" {
		v.addError(errors.New("Networking type cannot be empty"))
	}
	if n.Type != "bridged" && n.Type != "overlay" {
		v.addError(fmt.Errorf("Invalid networking type %q was provided", n.Type))
	}
	if n.PodCIDRBlock == "" {
		v.addError(errors.New("Pod CIDR block cannot be empty"))
	}
	if _, _, err := net.ParseCIDR(n.PodCIDRBlock); n.PodCIDRBlock != "" && err != nil {
		v.addError(fmt.Errorf("Invalid Pod CIDR block provided: %v", err))
	}

	if n.ServiceCIDRBlock == "" {
		v.addError(errors.New("Service CIDR block cannot be empty"))
	}
	if _, _, err := net.ParseCIDR(n.ServiceCIDRBlock); n.ServiceCIDRBlock != "" && err != nil {
		v.addError(fmt.Errorf("Invalid Service CIDR block provided: %v", err))
	}
	return v.valid()
}

func (c *CertsConfig) validate() (bool, []error) {
	v := newValidator()
	if _, err := time.ParseDuration(c.Expiry); err != nil {
		v.addError(fmt.Errorf("Invalid certificate expiry %q provided: %v", c.Expiry, err))
	}

	if c.LocationCity == "" {
		v.addError(errors.New("Certificate location_city field is required"))
	}
	if c.LocationState == "" {
		v.addError(errors.New("Certificate location_state field is required"))
	}
	if c.LocationCountry == "" {
		v.addError(errors.New("Certificate location_country field is required"))
	}
	return v.valid()
}

func (s *SSHConfig) validate() (bool, []error) {
	v := newValidator()
	if s.User == "" {
		v.addError(errors.New("SSH user field is required"))
	}
	if s.Key == "" {
		v.addError(errors.New("SSH key field is required"))
	}
	if s.Port < 1 || s.Port > 65535 {
		v.addError(fmt.Errorf("SSH port %d is invalid. Port must be in the range 1-65535", s.Port))
	}
	return v.valid()
}

func (ng *NodeGroup) validate() (bool, []error) {
	v := newValidator()
	if len(ng.Nodes) != ng.ExpectedCount {
		v.addError(fmt.Errorf("Expected node count (%d) does not match the number of nodes provided (%d)", ng.ExpectedCount, len(ng.Nodes)))
	}
	for i, n := range ng.Nodes {
		v.validateWithErrPrefix(fmt.Sprintf("Node #%d", i+1), &n)
	}
	return v.valid()
}

func (mng *MasterNodeGroup) validate() (bool, []error) {
	v := newValidator()
	v.validate(&mng.NodeGroup)

	if mng.LoadBalancedFQDN == "" {
		v.addError(fmt.Errorf("Load balanced FQDN is required"))
	}

	if mng.LoadBalancedShortName == "" {
		v.addError(fmt.Errorf("Load balanced shortname is required"))
	}

	return v.valid()
}

func (n *Node) validate() (bool, []error) {
	v := newValidator()
	if n.Host == "" {
		v.addError(fmt.Errorf("Node host field is required"))
	}
	if n.IP == "" {
		v.addError(fmt.Errorf("Node IP field is required"))
	}
	if ip := net.ParseIP(n.IP); ip == nil {
		v.addError(fmt.Errorf("Invalid IP provided"))
	}
	return v.valid()
}
