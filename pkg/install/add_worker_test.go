package install

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strconv"
	"testing"

	"github.com/apprenda/kismatic/pkg/ansible"
	"github.com/apprenda/kismatic/pkg/install/explain"
	"github.com/apprenda/kismatic/pkg/tls"
)

func mustGetTempDir(t *testing.T) string {
	dir, err := ioutil.TempDir("", "add-worker-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	return dir
}

func TestAddWorkerCertMissingCAMissing(t *testing.T) {
	e := ansibleExecutor{
		options:             ExecutorOptions{RestartServices: true, RunsDirectory: mustGetTempDir(t)},
		stdout:              ioutil.Discard,
		consoleOutputFormat: ansible.RawFormat,
		pki:                 &fakePKI{},
		certsDir:            mustGetTempDir(t),
	}
	originalPlan := &Plan{
		Worker: NodeGroup{
			Nodes: []Node{},
		},
	}
	newWorker := Node{}
	newPlan, err := e.AddWorker(originalPlan, newWorker)
	if newPlan != nil {
		t.Errorf("add worker returned an updated plan")
	}
	if err != errMissingClusterCA {
		t.Errorf("AddWorker did not return the expected error. Instead returned: %v", err)
	}
}

// Verify that cert gets generated
func TestAddWorkerCertMissingCAExists(t *testing.T) {
	pki := &fakePKI{
		caExists: true,
	}
	e := ansibleExecutor{
		options:             ExecutorOptions{RestartServices: true, RunsDirectory: mustGetTempDir(t)},
		stdout:              ioutil.Discard,
		consoleOutputFormat: ansible.RawFormat,
		pki:                 pki,
		runnerExplainerFactory: fakeRunnerExplainer(nil),
		certsDir:               mustGetTempDir(t),
	}
	originalPlan := &Plan{
		Worker: NodeGroup{
			Nodes: []Node{},
		},
		Cluster: Cluster{
			Networking: NetworkConfig{
				ServiceCIDRBlock: "10.0.0.0/16",
			},
		},
	}
	newWorker := Node{}
	_, err := e.AddWorker(originalPlan, newWorker)
	if err != nil {
		t.Errorf("unexpected error while adding worker: %v", err)
	}
	if pki.generateCACalled {
		t.Errorf("CA was generated, even though it already existed")
	}
	if !pki.generateNodeCertCalled {
		t.Error("node certificate was not generated")
	}
}

func TestAddWorkerPlanIsUpdated(t *testing.T) {
	e := ansibleExecutor{
		options:             ExecutorOptions{RestartServices: true, RunsDirectory: mustGetTempDir(t)},
		stdout:              ioutil.Discard,
		consoleOutputFormat: ansible.RawFormat,
		pki: &fakePKI{
			caExists: true,
		},
		runnerExplainerFactory: fakeRunnerExplainer(nil),
		certsDir:               mustGetTempDir(t),
	}
	originalPlan := &Plan{
		Worker: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "existingWorker",
				},
			},
		},
		Cluster: Cluster{
			Networking: NetworkConfig{
				ServiceCIDRBlock: "10.0.0.0/16",
			},
		},
	}
	newWorker := Node{
		Host: "test",
	}
	updatedPlan, err := e.AddWorker(originalPlan, newWorker)
	if err != nil {
		t.Errorf("unexpected error while adding worker: %v", err)
	}
	if updatedPlan.Worker.ExpectedCount != 2 {
		t.Errorf("expected count was not incremented")
	}
	found := false
	for _, w := range updatedPlan.Worker.Nodes {
		if w.Host == newWorker.Host {
			found = true
		}
	}
	if !found {
		t.Errorf("the updated plan does not include the new worker")
	}
}

func TestAddWorkerPlanNotUpdatedAfterFailure(t *testing.T) {
	e := ansibleExecutor{
		options:             ExecutorOptions{RestartServices: true, RunsDirectory: mustGetTempDir(t)},
		stdout:              ioutil.Discard,
		consoleOutputFormat: ansible.RawFormat,
		pki: &fakePKI{
			caExists: true,
		},
		runnerExplainerFactory: fakeRunnerExplainer(errors.New("exec error")),
		certsDir:               mustGetTempDir(t),
	}
	originalPlan := &Plan{
		Worker: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "existingWorker",
				},
			},
		},
		Cluster: Cluster{
			Networking: NetworkConfig{
				ServiceCIDRBlock: "10.0.0.0/16",
			},
		},
	}
	newWorker := Node{
		Host: "test",
	}
	updatedPlan, err := e.AddWorker(originalPlan, newWorker)
	if err == nil {
		t.Errorf("expected an error, but didn't get one")
	}
	if updatedPlan != nil {
		t.Error("plan was updated, even though adding worker failed")
	}
}

func TestAddWorkerRestartServicesEnabled(t *testing.T) {
	fakeRunner := fakeRunner{}
	e := ansibleExecutor{
		certsDir:            mustGetTempDir(t),
		options:             ExecutorOptions{RestartServices: true, RunsDirectory: mustGetTempDir(t)},
		stdout:              ioutil.Discard,
		consoleOutputFormat: ansible.RawFormat,
		pki: &fakePKI{
			caExists: true,
		},
		runnerExplainerFactory: func(explain.AnsibleEventExplainer, io.Writer) (ansible.Runner, *explain.AnsibleEventStreamExplainer, error) {
			return &fakeRunner, &explain.AnsibleEventStreamExplainer{}, nil
		},
	}
	originalPlan := &Plan{
		Worker: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "existingWorker",
				},
			},
		},
		Cluster: Cluster{
			Networking: NetworkConfig{
				ServiceCIDRBlock: "10.0.0.0/16",
			},
		},
	}
	newWorker := Node{
		Host: "test",
	}
	_, err := e.AddWorker(originalPlan, newWorker)
	if err != nil {
		t.Errorf("unexpected error")
	}
	expectedServicesToRestart := []string{"proxy", "kubelet", "calico_node", "docker"}
	for _, svc := range expectedServicesToRestart {
		found := false
		for k, v := range fakeRunner.incomingVars["kubernetes-worker.yaml"] {
			if k == fmt.Sprintf("force_%s_restart", svc) && v == strconv.FormatBool(true) {
				found = true
			}
		}
		if !found {
			t.Errorf("missing restart flag for service %s", svc)
		}
	}
}

func TestAddWorkerHostsFilesDNSEnabled(t *testing.T) {
	fakeRunner := fakeRunner{}
	e := ansibleExecutor{
		options:             ExecutorOptions{RunsDirectory: mustGetTempDir(t)},
		stdout:              ioutil.Discard,
		consoleOutputFormat: ansible.RawFormat,
		pki: &fakePKI{
			caExists: true,
		},
		runnerExplainerFactory: func(explain.AnsibleEventExplainer, io.Writer) (ansible.Runner, *explain.AnsibleEventStreamExplainer, error) {
			return &fakeRunner, &explain.AnsibleEventStreamExplainer{}, nil
		},
		certsDir: mustGetTempDir(t),
	}
	originalPlan := &Plan{
		Worker: NodeGroup{
			ExpectedCount: 1,
			Nodes: []Node{
				{
					Host: "existingWorker",
				},
			},
		},
		Cluster: Cluster{
			Networking: NetworkConfig{
				ServiceCIDRBlock: "10.0.0.0/16",
				UpdateHostsFiles: true,
			},
		},
	}
	newWorker := Node{
		Host: "test",
	}
	_, err := e.AddWorker(originalPlan, newWorker)
	if err != nil {
		t.Errorf("unexpected error")
	}
	expectedPlaybook := "_hosts.yaml"
	if fakeRunner.allNodesPlaybook != expectedPlaybook {
		t.Errorf("expected playbook %s, but ran %s", expectedPlaybook, fakeRunner.allNodesPlaybook)
	}
}

//// Fakes for testing
type fakePKI struct {
	caExists               bool
	nodeCertExists         bool
	err                    error
	generateCACalled       bool
	generateNodeCertCalled bool
}

func (f *fakePKI) CertificateAuthorityExists() (bool, error)     { return f.caExists, f.err }
func (f *fakePKI) NodeCertificateExists(node Node) (bool, error) { return f.nodeCertExists, f.err }
func (f *fakePKI) GenerateNodeCertificate(plan *Plan, node Node, ca *tls.CA) error {
	f.generateNodeCertCalled = true
	return f.err
}
func (f *fakePKI) GetClusterCA() (*tls.CA, error) { return nil, f.err }
func (f *fakePKI) GenerateClusterCA(p *Plan) (*tls.CA, error) {
	f.generateCACalled = true
	return nil, f.err
}
func (f *fakePKI) GenerateClusterCertificates(p *Plan, ca *tls.CA, users []string) error { return f.err }

type fakeRunner struct {
	eventChan        chan ansible.Event
	err              error
	incomingVars     map[string]ansible.ExtraVars
	allNodesPlaybook string
}

func (f *fakeRunner) StartPlaybook(playbookFile string, inventory ansible.Inventory, vars ansible.ExtraVars) (<-chan ansible.Event, error) {
	f.allNodesPlaybook = playbookFile
	return f.eventChan, f.err
}
func (f *fakeRunner) WaitPlaybook() error { return f.err }
func (f *fakeRunner) StartPlaybookOnNode(playbookFile string, inventory ansible.Inventory, vars ansible.ExtraVars, node string) (<-chan ansible.Event, error) {
	if f.incomingVars == nil {
		f.incomingVars = map[string]ansible.ExtraVars{}
	}
	f.incomingVars[playbookFile] = vars
	return f.eventChan, f.err
}

func fakeRunnerExplainer(execError error) func(explain.AnsibleEventExplainer, io.Writer) (ansible.Runner, *explain.AnsibleEventStreamExplainer, error) {
	return func(explain.AnsibleEventExplainer, io.Writer) (ansible.Runner, *explain.AnsibleEventStreamExplainer, error) {
		return &fakeRunner{err: execError}, &explain.AnsibleEventStreamExplainer{}, nil
	}
}
