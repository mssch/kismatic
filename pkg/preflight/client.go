package preflight

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/apprenda/kismatic-platform/pkg/preflight/check"
)

// The Client executes preflight checks against a remote node. The client expects a Server listening
// on the remote node.
type Client struct {
	// TargetNode is the ip:port to the remote node
	TargetNode string
}

// RunChecks runs the checks against a remote node
func (c Client) RunChecks(cr *CheckRequest) ([]CheckResult, error) {
	d, err := json.Marshal(cr)
	if err != nil {
		return nil, fmt.Errorf("error marshaling check request: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/run-checks", c.TargetNode), "application/json", bytes.NewReader(d))
	if err != nil {
		return nil, fmt.Errorf("server responded with error: %v", err)
	}
	defer resp.Body.Close()

	results := []CheckResult{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, fmt.Errorf("error decoding server response: %v", err)
	}

	node, _, err := net.SplitHostPort(c.TargetNode)
	if err != nil {
		return nil, fmt.Errorf("error getting host from %q: %v", c.TargetNode, err)
	}

	// Run TCP checks if any
	for _, p := range cr.TCPPorts {
		// Build check
		tcpCheck := check.TCPPortClientCheck{p, node}
		err := tcpCheck.Check()
		// Build result
		r := CheckResult{
			Name:    tcpCheck.Name(),
			Success: err == nil,
		}
		if err != nil {
			r.Error = fmt.Sprintf("%v", tcpCheck.Check())
		}
		results = append(results, r)
	}

	// TODO: Run SSH check
	resp, err = http.Get(fmt.Sprintf("http://%s/close-checks", c.TargetNode))
	if err != nil {
		// TODO: Handle this error? Or just log it
	}

	return results, nil
}
