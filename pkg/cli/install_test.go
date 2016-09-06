package cli

import (
	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/tls"
)

type fakePlan struct {
	exists     bool
	plan       *install.Plan
	err        error
	readCalled bool
}

func (fp *fakePlan) PlanExists() bool { return fp.exists }
func (fp *fakePlan) Read() (*install.Plan, error) {
	fp.readCalled = true
	return fp.plan, fp.err
}
func (fp *fakePlan) Write(p *install.Plan) error {
	fp.plan = p
	return fp.err
}

type fakeExecutor struct {
	installCalled bool
	err           error
}

func (fe *fakeExecutor) GetVars(p *install.Plan, options *install.CliOpts) (*install.AnsibleVars, error) {
	return &install.AnsibleVars{}, fe.err
}

func (fe *fakeExecutor) Install(p *install.Plan, av *install.AnsibleVars) error {
	fe.installCalled = true
	return fe.err
}

type fakePKI struct {
	called bool
	err    error
}

func (fp *fakePKI) ReadClusterCA(p *install.Plan) (*tls.CA, error) {
	return nil, fp.err
}
func (fp *fakePKI) GenerateClusterCA(p *install.Plan) (*tls.CA, error) {
	return nil, fp.err
}

func (fp *fakePKI) GenerateClusterCerts(p *install.Plan, ca *tls.CA, users []string) error {
	fp.called = true
	return fp.err
}
