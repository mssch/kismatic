package cli

import (
	"github.com/apprenda/kismatic-platform/pkg/install"
	"github.com/apprenda/kismatic-platform/pkg/tls"
)

type fakePlanner struct {
	exists     bool
	plan       *install.Plan
	err        error
	readCalled bool
}

func (fp *fakePlanner) PlanExists() bool { return fp.exists }
func (fp *fakePlanner) Read() (*install.Plan, error) {
	fp.readCalled = true
	return fp.plan, fp.err
}
func (fp *fakePlanner) Write(p *install.Plan) error {
	fp.plan = p
	return fp.err
}

type fakeExecutor struct {
	installCalled bool
	err           error
}

func (fe *fakeExecutor) Install(p *install.Plan) error {
	fe.installCalled = true
	return fe.err
}

func (fe *fakeExecutor) RunPreflightCheck(p *install.Plan) error {
	return nil
}

type fakePKI struct {
	called              bool
	generateCACalled    bool
	readClusterCACalled bool
	err                 error
}

func (fp *fakePKI) ReadClusterCA(p *install.Plan) (*tls.CA, error) {
	fp.readClusterCACalled = true
	return nil, fp.err
}
func (fp *fakePKI) GenerateClusterCA(p *install.Plan) (*tls.CA, error) {
	fp.generateCACalled = true
	return nil, fp.err
}

func (fp *fakePKI) GenerateClusterCerts(p *install.Plan, ca *tls.CA, users []string) error {
	fp.called = true
	return fp.err
}
