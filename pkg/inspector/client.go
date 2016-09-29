package inspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
)

// The Client executes rules against a remote inspector
type Client struct {
	// TargetNode is the ip:port of the inspector running on the remote node
	TargetNode string
}

// ExecuteRules against the target inspector server
func (c Client) ExecuteRules(rules []rule.Rule) ([]rule.RuleResult, error) {
	d, err := json.Marshal(rules)
	if err != nil {
		return nil, fmt.Errorf("error marshaling check request: %v", err)
	}
	resp, err := http.Post(fmt.Sprintf("http://%s%s", c.TargetNode, executeEndpoint), "application/json", bytes.NewReader(d))
	if err != nil {
		return nil, fmt.Errorf("error posting request to server: %v", err)
	}
	defer resp.Body.Close()
	// verify response status code
	if resp.StatusCode == http.StatusInternalServerError {
		errMsg := &serverError{}
		if err = json.NewDecoder(resp.Body).Decode(errMsg); err != nil {
			return nil, fmt.Errorf("failed to decode server response: %v. Server sent %q status", err, resp.Status)
		}
		return nil, fmt.Errorf("server sent %q status: error from server: %s", http.StatusInternalServerError, errMsg.Error)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server responded with non-successful status: %q", resp.Status)
	}

	// we got an OK - handle the response
	results := []rule.RuleResult{}
	err = json.NewDecoder(resp.Body).Decode(&results)
	if err != nil {
		return nil, fmt.Errorf("error decoding server response: %v", err)
	}

	// TODO: Run remote rules here...
	// node, _, err := net.SplitHostPort(c.TargetNode)
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting host from %q: %v", c.TargetNode, err)
	// }
	// // Run TCP checks if any
	// for _, p := range m.OpenTCPPorts {
	// 	// Build check
	// 	tcpCheck := TCPPortClientCheck{p, node}
	// 	err := tcpCheck.Check()
	// 	// Build result
	// 	r := CheckResult{
	// 		Name:    tcpCheck.Name(),
	// 		Success: err == nil,
	// 	}
	// 	if err != nil {
	// 		r.Error = fmt.Sprintf("%v", tcpCheck.Check())
	// 	}
	// 	results = append(results, r)
	// }

	// TODO: add retry logic here?
	resp, err = http.Get(fmt.Sprintf("http://%s%s", c.TargetNode, closeEndpoint))
	if err != nil {
		// TODO: Handle this error? Or just log it
	}

	return results, nil
}
