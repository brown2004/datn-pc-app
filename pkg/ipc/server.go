package ipc

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"
)

type Handler interface {
	Status(ctx context.Context) (Status, error)
	Identity(ctx context.Context) (Identity, error)
	CredentialsStatus(ctx context.Context) (CredentialsStatus, error)
	SetProtection(ctx context.Context, enabled bool) (Status, error)
	SetCredentials(ctx context.Context, pcAgentID string, agentSecret string) (Status, error)
}

type Server struct {
	addr       string
	handler    Handler
	httpServer *http.Server
}

func NewServer(addr string, handler Handler) *Server {
	if addr == "" {
		addr = DefaultAddr
	}

	server := &Server{
		addr:    addr,
		handler: handler,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/status", server.handleStatus)
	mux.HandleFunc("/identity", server.handleIdentity)
	mux.HandleFunc("/protection", server.handleProtection)
	mux.HandleFunc("/credentials", server.handleCredentials)
	mux.HandleFunc("/credentials/status", server.handleCredentialsStatus)

	server.httpServer = &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return server
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = s.httpServer.Shutdown(shutdownCtx)
	}()

	err = s.httpServer.Serve(listener)
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return err
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	status, err := s.handler.Status(r.Context())
	writeStatusResponse(w, status, err)
}

func (s *Server) handleIdentity(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	identity, err := s.handler.Identity(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, identity)
}

func (s *Server) handleProtection(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch && r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	var req ProtectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}

	status, err := s.handler.SetProtection(r.Context(), req.Enabled)
	writeStatusResponse(w, status, err)
}

func (s *Server) handleCredentialsStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	status, err := s.handler.CredentialsStatus(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func (s *Server) handleCredentials(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed")
		return
	}

	var req CredentialsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json")
		return
	}
	if strings.TrimSpace(req.PCAgentID) == "" || strings.TrimSpace(req.AgentSecret) == "" {
		writeError(w, http.StatusBadRequest, "invalid_credentials")
		return
	}

	status, err := s.handler.SetCredentials(r.Context(), req.PCAgentID, req.AgentSecret)
	writeStatusResponse(w, status, err)
}

func writeStatusResponse(w http.ResponseWriter, status Status, err error) {
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, status)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, ErrorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
