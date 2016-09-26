package inspector

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

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

		m := &Manifest{}
		err := json.NewDecoder(req.Body).Decode(m)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer req.Body.Close()

		results := s.RunChecks(m)
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
func (s *Server) RunChecks(m *Manifest) []CheckResult {
	checks := []Check{}
	for _, c := range m.AvailablePackageDependencies {
		checks = append(checks, c)
	}
	for _, c := range m.InstalledPackageDependencies {
		checks = append(checks, c)
	}
	for _, c := range m.BinaryDependencies {
		checks = append(checks, c)
	}

	closable := []ClosableCheck{}
	for _, p := range m.OpenTCPPorts {
		c := &TCPPortServerCheck{PortNumber: p}
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
