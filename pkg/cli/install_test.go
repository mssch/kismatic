package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/install"
)

type fakePlan struct {
	exists     bool
	plan       *install.Plan
	err        error
	readCalled bool
}

func (fp *fakePlan) Exists() bool { return fp.exists }
func (fp *fakePlan) Read() (*install.Plan, error) {
	fp.readCalled = true
	return fp.plan, fp.err
}
func (fp *fakePlan) Write(p *install.Plan) error {
	fp.plan = p
	return fp.err
}

func TestInstallCmdPlanNotFound(t *testing.T) {
	tests := []struct {
		in             io.Reader
		shouldError    bool
		expectedEtcd   int
		expectedMaster int
		expectedWorker int
	}{
		{
			// User accepts default node counts
			in:             strings.NewReader("\n\n\n"),
			expectedEtcd:   3,
			expectedMaster: 2,
			expectedWorker: 3,
		},
		{
			// User enteres node countes
			in:             strings.NewReader("10\n10\n10\n"),
			expectedEtcd:   10,
			expectedMaster: 10,
			expectedWorker: 10,
		},
		{
			// User enters invalid numeric input
			in:          strings.NewReader("0\n1\n1\n"),
			shouldError: true,
		},
		{
			// User enters invalid input
			in:          strings.NewReader("badInput\n"),
			shouldError: true,
		},
		{
			// User enters invalid input
			in:          strings.NewReader("badInput\nother\nfail\n"),
			shouldError: true,
		},
	}

	for _, test := range tests {
		out := &bytes.Buffer{}
		fp := &fakePlan{}
		cmd := NewCmdInstall(test.in, out, fp)
		cmd.SilenceUsage = true
		cmd.SilenceErrors = true

		err := cmd.Execute()

		if err == nil && test.shouldError {
			t.Error("expected an error, but did not get one")
		}

		if err != nil && !test.shouldError {
			t.Errorf("unexpected error running command: %v", err)
		}

		if !test.shouldError {
			// Verify defaults are set in the plan
			p := fp.plan
			if p.Etcd.ExpectedCount != test.expectedEtcd {
				t.Errorf("expected %d etcd nodes, got %d", test.expectedEtcd, p.Etcd.ExpectedCount)
			}
			if p.Master.ExpectedCount != 2 {
				t.Errorf("expected %d master nodes, got %d", test.expectedMaster, p.Master.ExpectedCount)
			}
			if p.Worker.ExpectedCount != 3 {
				t.Errorf("expected %d worker nodes, got %d", test.expectedWorker, p.Worker.ExpectedCount)
			}
			// Verify we didn't attempt to read plan
			if fp.readCalled {
				t.Errorf("attempted to read plan, when the plan does not exist")
			}
		}
	}
}

func TestInstallCmdPlanFound(t *testing.T) {
	in := strings.NewReader("")
	out := &bytes.Buffer{}
	fp := &fakePlan{
		exists: true,
		plan:   &install.Plan{},
	}
	cmd := NewCmdInstall(in, out, fp)
	cmd.SilenceUsage = true
	cmd.SilenceErrors = true
	err := cmd.Execute()

	// expect an error here... we don't care about testing validation
	if err == nil {
		t.Error("expected error due to invalid plan, but did not get one")
	}

	if !fp.readCalled {
		t.Error("the plan was not read")
	}
}
