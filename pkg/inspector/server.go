package inspector

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// CheckResult contains the results of a check
type CheckResult struct {
	Name    string
	Success bool
	Error   string
}

// Server stands up an HTTP server that handles pre-flight checking.
type Server struct {
	ListenPort  int
	NodeLabels  []string
	rulesEngine Engine
}

// Start the CheckServer
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/run-checks", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		rules := []Rule{}
		err := json.NewDecoder(req.Body).Decode(rules)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		results := s.RunChecks(rules)
		json.NewEncoder(w).Encode(results)
	})

	mux.HandleFunc("/close-checks", func(w http.ResponseWriter, req *http.Request) {
		// TODO: Handle errors when closing
		s.rulesEngine.CloseChecks()
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.ListenPort), mux)
}

// RunChecks according to the check request and return the collection of results.
func (s *Server) RunChecks(rules []Rule) []RuleResult {
	res := s.rulesEngine.ExecuteRules(rules, s.NodeLabels)
	return res
}
