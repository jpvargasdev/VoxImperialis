package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// NotifyRequest is the payload received from MachinusCronus.
type NotifyRequest struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// NotifyServer runs an HTTP server that receives notifications
// from other services and forwards them as XMPP messages.
type NotifyServer struct {
	xmpp *XMPPClient
	port string
}

// NewNotifyServer creates a notify server.
func NewNotifyServer(xmpp *XMPPClient, port string) *NotifyServer {
	return &NotifyServer{xmpp: xmpp, port: port}
}

// Start begins listening for HTTP requests. Blocks until error.
func (s *NotifyServer) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/notify", s.handleNotify)
	mux.HandleFunc("/health", s.handleHealth)

	addr := ":" + s.port
	log.Printf("notify: HTTP server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *NotifyServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func (s *NotifyServer) handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var req NotifyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Status == "" {
		http.Error(w, "name and status are required", http.StatusBadRequest)
		return
	}

	log.Printf("notify: received %s for %q: %s", req.Type, req.Name, req.Status)

	// Format as XMPP message
	text := fmt.Sprintf("[watchface]\nname:    %s\nstatus:  %s\nmessage: %s", req.Name, req.Status, req.Message)

	if err := s.xmpp.SendToAllowed(text); err != nil {
		log.Printf("notify: failed to send XMPP message: %v", err)
		http.Error(w, "failed to send notification", http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"ok":true}`)
}
