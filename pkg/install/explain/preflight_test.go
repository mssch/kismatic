package explain

import (
	"fmt"
	"testing"

	"github.com/apprenda/kismatic-platform/pkg/ansible"
)

func TestPreFlightExplainer(t *testing.T) {
	exp := PreflightEventExplainer{&DefaultEventExplainer{}}
	start := &ansible.PlayStartEvent{}
	start.Name = "Starting checks..."

	fail01 := &ansible.RunnerFailedEvent{}
	fail01.Host = "host01"

	fail02 := &ansible.RunnerFailedEvent{}
	fail02.Host = "host02"

	events := []ansible.Event{start, fail01, fail02}

	fmt.Println()
	fmt.Println()
	for _, e := range events {
		fmt.Print(exp.ExplainEvent(e, false))
	}
	fmt.Println()
	fmt.Println()

}
