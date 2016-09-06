package preflight

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/apprenda/kismatic-platform/pkg/preflight/check"
)

// CheckRequest contains the list of checks that should be run
type CheckRequest struct {
	BinaryDependencies  []string
	PackageDependencies []string
	TCPPorts            []int
}

// CheckResult contains the results of a check
type CheckResult struct {
	Name    string
	Success bool
	Error   string
}

// Server stands up an HTTP server that handles pre-flight checking.
type Server struct {
	ListenPort     int
	mu             sync.Mutex
	closableChecks []ClosableCheck
}

// Start the CheckServer
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/run-checks", func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		cr := &CheckRequest{}
		err := json.NewDecoder(req.Body).Decode(cr)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		results := s.RunChecks(cr)
		json.NewEncoder(w).Encode(results)
	})

	mux.HandleFunc("/close-checks", func(w http.ResponseWriter, req *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()

		for _, c := range s.closableChecks {
			// TODO: What to do with error here?
			c.Close()
		}

		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.ListenPort), mux)
}

// RunChecks according to the check request and return the collection of results.
func (s *Server) RunChecks(cr *CheckRequest) []CheckResult {
	checks := []Check{}
	for _, b := range cr.BinaryDependencies {
		checks = append(checks, &check.BinaryDependencyCheck{b})
	}

	for _, p := range cr.PackageDependencies {
		checks = append(checks, &check.PackageInstalledCheck{p})
	}

	closable := []ClosableCheck{}
	for _, p := range cr.TCPPorts {
		c := &check.TCPPortServerCheck{PortNumber: p}
		checks = append(checks, c)
		closable = append(closable, c)
	}

	results := []CheckResult{}
	for _, c := range checks {
		// Run the check
		err := c.Check()
		// Build result
		r := CheckResult{
			Name:    c.Name(),
			Success: err == nil,
		}
		if err != nil {
			r.Error = fmt.Sprintf("%v", c.Check())
		}
		results = append(results, r)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.closableChecks = closable

	return results
}
