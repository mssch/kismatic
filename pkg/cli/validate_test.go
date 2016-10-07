package cli

import (
	"bytes"
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/install"
)

func TestValidateCmdPlanNotFound(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		exists: false,
	}
	opts := &validateOpts{
		planFile:     "planFile",
		verbose:      false,
		outputFormat: "table",
	}
	err := doValidate(out, fp, opts)
	if err == nil {
		t.Errorf("validate did not return an error when the plan does not exist")
	}

	if fp.readCalled {
		t.Errorf("attempted to read a non-existent plan file")
	}
}

func TestValidateCmdPlanInvalid(t *testing.T) {
	out := &bytes.Buffer{}
	fp := &fakePlanner{
		exists: true,
		plan:   &install.Plan{},
	}
	opts := &validateOpts{
		planFile:     "planFile",
		verbose:      false,
		outputFormat: "table",
	}
	err := doValidate(out, fp, opts)
	if err == nil {
		t.Errorf("did not return an error with an invalid plan")
	}

	if !fp.readCalled {
		t.Errorf("did not read the plan file")
	}
}
