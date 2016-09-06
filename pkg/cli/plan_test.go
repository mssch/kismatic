package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/install"
)

func TestPlanCmdPlanNotFound(t *testing.T) {
	tests := []struct {
		in             io.Reader
		shouldError    bool
		expectedEtcd   int
		expectedMaster int
		expectedWorker int
	}{
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

		err := doPlan(test.in, out, fp, &install.CliOpts{})

		if err == nil && test.shouldError {
			t.Error("expected an error, but did not get one")
		}
	}
}
