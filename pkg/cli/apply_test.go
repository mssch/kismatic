package cli

import (
	"bytes"
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/install"
)

func TestApplyCmdInvalidPlanFound(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		exists: true,
		plan:   &install.Plan{},
	}
	fe := &fakeExecutor{}

	applyCmd := &applyCmd{
		out:      out,
		planner:  fp,
		executor: fe,
	}

	err := applyCmd.run(false, "table")

	// expect an error here... we don't care about testing validation
	if err == nil {
		t.Error("expected error due to invalid plan, but did not get one")
	}

	if !fp.readCalled {
		t.Error("the plan was not read")
	}

	if fe.installCalled {
		t.Error("install was called with an invalid plan")
	}
}

// TODO: put plan validation behind interface to enable these tests
// func TestApplyCmdSkipCAGeneration(t *testing.T) {
// 	out := &bytes.Buffer{}
// 	fp := &fakePlanner{
// 		exists: true,
// 		plan:   &install.Plan{},
// 	}
// 	fe := &fakeExecutor{}
// 	fpki := &fakePKI{}

// 	applyCmd := &applyCmd{
// 		out:      out,
// 		planner:  fp,
// 		executor: fe,
// 		pki:      fpki,
// 		skipCAGeneration: true
// 	}

// 	applyCmd.run()
// 	if fpki.generateCACalled {
// 		t.Errorf("generated CA when skip CA generation was set to true")
// 	}

// 	if !fpki.readClusterCACalled {
// 		t.Errorf("did not read CA cert when skip CA generation was set to true")
// 	}
// }
