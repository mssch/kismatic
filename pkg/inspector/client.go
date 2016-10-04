package inspector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/apprenda/kismatic-platform/pkg/inspector/rule"
)

// The Client executes rules against a remote inspector
type Client struct {
	// TargetNode is the ip:port of the inspector running on the remote node
	TargetNode string
	// TargetNodeRole is the role of the node we are inspecting
	TargetNodeRole string
	engine         *rule.Engine
}

// NewClient returns an inspector client for running checks against remote nodes.
func NewClient(targetNode string, nodeRole string) (*Client, error) {
	host, _, err := net.SplitHostPort(targetNode)
	if err != nil {
		return nil, err
	}
	engine := &rule.Engine{
		RuleCheckMapper: rule.DefaultCheckMapper{
			PackageManager: nil, // Use a no-op pkg manager here instead
			TargetNodeIP:   host,
		},
	}
	return &Client{
		TargetNode:     targetNode,
		TargetNodeRole: nodeRole,
		engine:         engine,
	}, nil
}

// ExecuteRules against the target inspector server
func (c Client) ExecuteRules(rules []rule.Rule) ([]rule.RuleResult, error) {
	serverSideRules := getServerSideRules(rules)
	d, err := json.Marshal(serverSideRules)
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

	// Execute the rules that should run from a remote node
	clientSideRules := getClientSideRules(rules)
	facts := []string{c.TargetNodeRole}
	remoteResults, err := c.engine.ExecuteRules(clientSideRules, facts)
	if err != nil {
		return nil, err
	}
	results = append(results, remoteResults...)

	endpoint := fmt.Sprintf("http://%s%s", c.TargetNode, closeEndpoint)
	resp, err = http.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("GET request to %q failed. You might have to restart the inspector server. Error was: %v", endpoint, err)
	}

	return results, nil
}

func getServerSideRules(rules []rule.Rule) []rule.Rule {
	localRules := []rule.Rule{}
	for _, r := range rules {
		if !r.IsRemoteRule() {
			localRules = append(localRules, r)
		}
	}
	return localRules
}

func getClientSideRules(rules []rule.Rule) []rule.Rule {
	remoteRules := []rule.Rule{}
	for _, r := range rules {
		if r.IsRemoteRule() {
			remoteRules = append(remoteRules, r)
		}
	}
	return remoteRules
}
