package install

import (
	"fmt"
	"io/ioutil"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type Plan struct {
	ClusterName string `yaml:"clusterName"`

	SshUser    string `yaml:"sshUser"`
	SshPort    int    `yaml:"sshPort"`
	SshKeyPath string `yaml:"sshKeyPath"`

	EtcdNodeCount int      `yaml:"etcdNodeCount"`
	EtcdNodes     []string `yaml:"etcdNodes"`

	MasterNodeCount int      `yaml:"masterNodeCount"`
	MasterNodes     []string `yaml:"masterNodes"`

	WorkerNodeCount int      `yaml:"workerNodeCount"`
	WorkerNodes     []string `yaml:"workerNodes"`
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
	p.SshPort = 22

	// Generate entries for all node types
	for i := 0; i < p.EtcdNodeCount; i++ {
		p.EtcdNodes = append(p.EtcdNodes, fmt.Sprintf("etcd%d", i))
	}

	for i := 0; i < p.MasterNodeCount; i++ {
		p.MasterNodes = append(p.MasterNodes, fmt.Sprintf("master%d", i))
	}

	for i := 0; i < p.WorkerNodeCount; i++ {
		p.WorkerNodes = append(p.WorkerNodes, fmt.Sprintf("worker%d", i))
	}

	if err := w.Write(&p); err != nil {
		return fmt.Errorf("error writing installation plan template: %v", err)
	}
	return nil
}

func ValidatePlan(p *Plan) error {
	if p.SshKeyPath == "" {
		return fmt.Errorf("SshKeyPath is empty. The path to the SSH key is required.")
	}

	return nil
}
