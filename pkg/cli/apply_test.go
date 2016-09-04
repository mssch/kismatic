package cli

import (
	"bytes"
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/install"
)

func TestApplyCmdInvalidPlanFound(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlan{
		exists: true,
		plan:   &install.Plan{},
	}
	fe := &fakeExecutor{}
	fpki := &fakePKI{}
	err := doApply(out, fp, fe, fpki, &installOpts{})

	// expect an error here... we don't care about testing validation
	if err == nil {
		t.Error("expected error due to invalid plan, but did not get one")
	}

	if !fp.readCalled {
		t.Error("the plan was not read")
	}

	if fpki.called {
		t.Error("cert generation was called with an invalid plan")
	}

	if fe.installCalled {
		t.Error("install was called with an invalid plan")
	}
}
