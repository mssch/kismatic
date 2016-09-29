package inspector

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// CheckResult contains the results of a check
type CheckResult struct {
	Name    string
	Success bool
	Error   string
}

// Server supports the execution of inspector rules from a remote node
type Server struct {
	// The Port the server will listen on
	Port int
	// NodeFacts are the facts that are passed to the rules engine
	NodeFacts   []string
	rulesEngine Engine
}

type serverError struct {
	Error string
}

// Start the server
func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/execute", func(w http.ResponseWriter, req *http.Request) {
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

		results, err := s.rulesEngine.ExecuteRules(rules, s.NodeFacts)
		if err != nil {
			err = json.NewEncoder(w).Encode(serverError{Error: err.Error()})
			if err != nil {
				log.Printf("error writing server response: %v\n", err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = json.NewEncoder(w).Encode(results)
		if err != nil {
			log.Printf("error writing server response: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("/close", func(w http.ResponseWriter, req *http.Request) {
		err := s.rulesEngine.CloseChecks()
		if err != nil {
			log.Printf("error closing checks: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	})

	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
}
