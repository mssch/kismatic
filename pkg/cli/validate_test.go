package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestValidateCmdSucess(t *testing.T) {
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
			in:             strings.NewReader("8\n\n\n"),
			expectedEtcd:   8,
			expectedMaster: 2,
			expectedWorker: 3,
		},
	}
	for _, test := range tests {
		out := &bytes.Buffer{}
		fp := &fakePlan{
			exists: true,
		}

		err := doPlan(test.in, out, fp, &installOpts{})

		if err != nil && !test.shouldError {
			t.Errorf("unexpected error running command: %v", err)
		}

		if !test.shouldError {
			// Verify defaults are set in the plan
			p := fp.plan
			if p.Etcd.ExpectedCount != test.expectedEtcd {
				t.Errorf("expected %d etcd nodes, got %d", test.expectedEtcd, p.Etcd.ExpectedCount)
			}
			if p.Master.ExpectedCount != test.expectedMaster {
				t.Errorf("expected %d master nodes, got %d", test.expectedMaster, p.Master.ExpectedCount)
			}
			if p.Worker.ExpectedCount != test.expectedWorker {
				t.Errorf("expected %d worker nodes, got %d", test.expectedWorker, p.Worker.ExpectedCount)
			}
		}
	}
}
