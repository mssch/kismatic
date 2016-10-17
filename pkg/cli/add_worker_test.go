package cli

import (
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/install"
)

func TestAddWorkerToPlan(t *testing.T) {
	plan := install.Plan{}
	node := install.Node{
		Host:       "newWorker",
		IP:         "someIP",
		InternalIP: "someOtherIP",
	}
	newPlan := addWorkerToPlan(plan, node)
	if newPlan.Worker.ExpectedCount != 1 {
		t.Errorf("expected worker count was %d, wanted 1", newPlan.Worker.ExpectedCount)
	}
	if plan.Worker.ExpectedCount != 0 {
		t.Error("original plan was modified")
	}
}
