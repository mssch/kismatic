package ansible

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	// RawFormat is the raw Ansible output formatting
	RawFormat = OutputFormat("raw")
	// JSONLinesFormat is a JSON Lines representation of Ansible events
	JSONLinesFormat = OutputFormat("json_lines")
)

// OutputFormat is used for controlling the STDOUT format of the Ansible runner
type OutputFormat string

// Runner for running Ansible playbooks
type Runner interface {
	RunPlaybook(inventory Inventory, playbookFile string, vars ExtraVars) error
}

type runner struct {
	// Out is the stdout writer for the Ansible process
	out io.Writer
	// ErrOut is the stderr writer for the Ansible process
	errOut io.Writer

	pythonPath   string
	ansibleDir   string
	outputFormat OutputFormat
	verbose      bool
}

// ExtraVars is a map of variables that are used when executing a playbook
type ExtraVars map[string]string

func (v ExtraVars) commandLineVars() (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", fmt.Errorf("error marshaling ansible vars")
	}
	return string(b), nil
}

// NewRunner returns a new runner for executing Ansible commands.
func NewRunner(out, errOut io.Writer, ansibleDir string, outputFormat OutputFormat, verbose bool) (Runner, error) {
	ppath, err := getPythonPath()
	if err != nil {
		return nil, err
	}

	return &runner{
		out:          out,
		errOut:       errOut,
		pythonPath:   ppath,
		ansibleDir:   ansibleDir,
		outputFormat: outputFormat,
		verbose:      verbose,
	}, nil
}

// RunPlaybook with the given inventory and extra vars
func (r *runner) RunPlaybook(inv Inventory, playbookFile string, vars ExtraVars) error {
	extraVars, err := vars.commandLineVars()
	if err != nil {
		return fmt.Errorf("error building extra vars: %v", err)
	}

	inventoryFile := filepath.Join(r.ansibleDir, "inventory.ini")
	if err := ioutil.WriteFile(inventoryFile, inv.toINI(), 0644); err != nil {
		return fmt.Errorf("error writing inventory file to %q: %v", inventoryFile, err)
	}

	playbook := filepath.Join(r.ansibleDir, "playbooks", playbookFile)
	cmd := exec.Command(filepath.Join(r.ansibleDir, "bin", "ansible-playbook"), "-i", inventoryFile, "-s", playbook, "--extra-vars", extraVars)
	cmd.Stdout = r.out
	cmd.Stderr = r.errOut
	os.Setenv("PYTHONPATH", r.pythonPath)
	os.Setenv("ANSIBLE_HOST_KEY_CHECKING", "False")
	os.Setenv("ANSIBLE_CALLBACK_PLUGINS", filepath.Join(r.ansibleDir, "playbooks", "callback"))
	if r.outputFormat != RawFormat {
		os.Setenv("ANSIBLE_STDOUT_CALLBACK", string(r.outputFormat))
	}

	if r.verbose {
		cmd.Args = append(cmd.Args, "-vvvv")
	}

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running playbook: %v", err)
	}

	return nil
}

func getPythonPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting working dir: %v", err)
	}
	lib := filepath.Join(wd, "ansible", "lib", "python2.7", "site-packages")
	lib64 := filepath.Join(wd, "ansible", "lib64", "python2.7", "site-packages")
	return fmt.Sprintf("%s:%s", lib, lib64), nil
}
