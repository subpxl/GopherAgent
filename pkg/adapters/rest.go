package adapters

import (
	"agent-in-go/pkg/session"
	"encoding/json"
	"fmt"
	"net/http"
)

type RESTAdapter struct {
	port  string
	store *session.SessionStore
}

func NewRESTAdapter(port string, store *session.SessionStore) *RESTAdapter {
	return &RESTAdapter{port: port, store: store}
}

func (a *RESTAdapter) Name() string { return "REST" }

func (a *RESTAdapter) Start() error {
	http.HandleFunc("/ask", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}
		var req struct {
			SessionID string `json:"session_id"`
			Question  string `json:"question"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Question == "" {
			http.Error(w, `need {"session_id":"...","question":"..."}`, http.StatusBadRequest)
			return
		}
		if req.SessionID == "" {
			req.SessionID = "rest-default"
		}
		answer := a.store.Ask(req.SessionID, req.Question)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(struct {
			Answer    string `json:"answer"`
			SessionID string `json:"session_id"`
		}{Answer: answer, SessionID: req.SessionID})
	})
	
	fmt.Printf("[REST] listening on :%s\n", a.port)
	return http.ListenAndServe(":"+a.port, nil)
}
