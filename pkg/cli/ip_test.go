package cli

import (
	"bytes"
	"testing"

	"github.com/apprenda/kismatic/pkg/install"
)

func TestIPCmdMissingPlan(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		exists: false,
	}
	opts := &ipOpts{
		planFilename: "planFile",
	}
	if err := doIP(out, fp, opts); err == nil {
		t.Errorf("ip did not return an error when the plan does not exist")
	}
}

func TestIPCmdEmptyAddress(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		plan:   &install.Plan{},
		exists: true,
	}
	opts := &ipOpts{
		planFilename: "planFile",
	}
	if err := doIP(out, fp, opts); err == nil {
		t.Errorf("ip did not return an error when LoadBalancer is empty")
	}
}

func TestIPCmdValidAddress(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		plan: &install.Plan{
			Master: install.MasterNodeGroup{
				LoadBalancer: "10.0.0.10:6443",
			},
		},
		exists: true,
	}
	opts := &ipOpts{
		planFilename: "planFile",
	}
	err := doIP(out, fp, opts)
	if err != nil {
		t.Errorf("ip returned an error %v", err)
	}
	if out.String() != fp.plan.Master.LoadBalancer+"\n" {
		t.Errorf("ip returned %s, but expectetd %s", out.String(), fp.plan.Master.LoadBalancer)
	}
}
