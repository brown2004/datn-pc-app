package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"pc-app/internal/backend"
	"pc-app/internal/store"
	"pc-app/pkg/ipc"
)

const (
	defaultBackendURL = "http://localhost:8081"
	requestTimeout    = 8 * time.Second
)

// App struct
type App struct {
	ctx context.Context
	mu  sync.RWMutex

	client             *http.Client
	backendClient      *backend.Client
	backendURL         string
	agentClient        *ipc.Client
	localStore         *store.Store
	protectionEnabled  bool
	serverStatus       string
	pairingStatus      string
	deviceCode         string
	pairingSessionID   string
	pairingExpiresAt   time.Time
	pendingCredentials *agentCredentials
	lastError          string
}

// NewApp creates a new App application struct
func NewApp() *App {
	backendURL := strings.TrimRight(envOrDefault(defaultBackendURL, "PCAPP_BACKEND_URL", "DATN_BACKEND_URL"), "/")
	app := &App{
		client:        backend.NewHTTPClient(requestTimeout),
		backendURL:    backendURL,
		serverStatus:  "unknown",
		pairingStatus: "not started",
		backendClient: backend.NewClient(backendURL, requestTimeout),
	}
	app.initLocalStore()
	app.initAgentClient()
	return app
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initAgentClient()
}

type Status struct {
	AgentStatus       string `json:"agent_status"`
	DeviceStatus      string `json:"device_status"`
	ProtectionStatus  string `json:"protection_status"`
	BackendURL        string `json:"backend_url"`
	ServerStatus      string `json:"server_status"`
	ProtectionEnabled bool   `json:"protection_enabled"`
	PairingStatus     string `json:"pairing_status"`
	DeviceCode        string `json:"device_code"`
	PairingExpiresAt  string `json:"pairing_expires_at"`
	MQTTStatus        string `json:"mqtt_status"`
	LastError         string `json:"last_error"`
	LastEvent         string `json:"last_event"`
	LastEventAt       string `json:"last_event_at"`
}

type agentCredentials struct {
	PCAgentID   string `json:"pc_agent_id"`
	AgentSecret string `json:"agent_secret"`
}

type startPairingResponse struct {
	PairingSessionID string `json:"pairing_session_id"`
	DeviceCode       string `json:"device_code"`
	ExpiresIn        int64  `json:"expires_in"`
}

type pairingStatusResponse struct {
	Status           string `json:"status"`
	PCAgentID        string `json:"pc_agent_id"`
	AgentSecret      string `json:"agent_secret"`
	CredentialIssued bool   `json:"credential_issued"`
}

type backendError struct {
	StatusCode int
	Code       string `json:"error"`
	Message    string `json:"message"`
}

func (e *backendError) Error() string {
	message := e.Message
	if message == "" {
		message = e.Code
	}
	if message == "" {
		return fmt.Sprintf("backend returned %d", e.StatusCode)
	}

	return fmt.Sprintf("backend %d: %s", e.StatusCode, message)
}

func (a *App) GetStatus() Status {
	a.refreshServerStatus()
	a.refreshBackendState()
	if err := a.syncPendingAgentCredentials(); err != nil {
		a.setLastError("Da lien ket backend, nhung chua dong bo PC Agent local: " + err.Error())
	}
	a.refreshBackendCredentialStatus()
	return a.mergeAgentStatus(a.snapshot())
}

func (a *App) SetProtectionMode(enabled bool) Status {
	agentClient, err := a.ensureAgentClient()
	if err != nil {
		a.setLastError("Khong ket noi duoc PC Agent local: " + err.Error())
		return a.snapshot()
	}

	ctx, cancel := a.requestContext()
	defer cancel()

	if _, err := agentClient.SetProtection(ctx, enabled); err != nil {
		a.setLastError("Khong cap nhat duoc protection qua PC Agent local: " + err.Error())
		return a.snapshot()
	}

	return a.GetStatus()
}

func (a *App) StartPairing() Status {
	pcAgentID, err := a.getLocalPCAgentID()
	if err != nil {
		a.setLastError("Khong lay duoc ID laptop local: " + err.Error())
		return a.snapshot()
	}
	if strings.TrimSpace(pcAgentID) == "" {
		a.setLastError("Chua co ID laptop local")
		return a.snapshot()
	}

	hostname, err := os.Hostname()
	if err != nil || strings.TrimSpace(hostname) == "" {
		hostname = "DATN PC"
	}

	var response startPairingResponse
	err = a.postJSON("/api/pc-agents/pairing/start", map[string]string{
		"pc_agent_id": pcAgentID,
		"device_name": hostname,
		"os_type":     runtime.GOOS,
	}, &response)
	if err != nil {
		a.setLastError(err.Error())
		return a.snapshot()
	}

	a.mu.Lock()
	a.deviceCode = response.DeviceCode
	a.pairingSessionID = response.PairingSessionID
	a.pairingExpiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	a.pairingStatus = "waiting for user"
	a.lastError = ""
	a.mu.Unlock()

	return a.GetStatus()
}

func (a *App) refreshBackendState() {
	pairingSessionID, deviceCode := a.currentBackendInputs()
	if pairingSessionID != "" && deviceCode != "" {
		a.pollPairingStatus(pairingSessionID, deviceCode)
	}
}

func (a *App) pollPairingStatus(pairingSessionID, deviceCode string) {
	var response pairingStatusResponse
	query := url.Values{}
	query.Set("pairing_session_id", pairingSessionID)
	query.Set("device_code", deviceCode)
	path := "/api/pc-agents/pairing/status?" + query.Encode()
	if err := a.getJSON(path, &response); err != nil {
		a.setLastError(err.Error())
		return
	}

	if response.Status != "confirmed" || response.PCAgentID == "" || response.AgentSecret == "" {
		a.mu.Lock()
		a.pairingStatus = normalizePairingStatus(response.Status)
		a.lastError = ""
		a.mu.Unlock()
		return
	}

	credentials := &agentCredentials{
		PCAgentID:   response.PCAgentID,
		AgentSecret: response.AgentSecret,
	}
	if err := a.saveLocalCredentials(credentials); err != nil {
		a.setLastError("Da lien ket backend, nhung chua luu duoc credential local: " + err.Error())
		return
	}

	a.mu.Lock()
	a.pendingCredentials = credentials
	a.deviceCode = ""
	a.pairingSessionID = ""
	a.pairingExpiresAt = time.Time{}
	a.pairingStatus = "syncing local agent"
	a.lastError = ""
	a.mu.Unlock()

	if err := a.syncAgentCredentials(credentials); err != nil {
		a.setLastError("Da lien ket backend, nhung chua dong bo PC Agent local: " + err.Error())
		return
	}

	a.mu.Lock()
	a.pendingCredentials = nil
	a.pairingStatus = "linked"
	a.lastError = ""
	a.mu.Unlock()
}

func (a *App) snapshot() Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	protectionStatus := "disabled"
	if a.protectionEnabled {
		protectionStatus = "enabled"
	}

	agentStatus := "offline"
	deviceStatus := "unknown"
	pairingStatus := a.pairingStatus
	if pairingStatus == "" {
		pairingStatus = "not started"
	}
	expiresAt := ""

	if a.pendingCredentials != nil {
		pairingStatus = "syncing local agent"
	} else if a.deviceCode != "" {
		if time.Now().After(a.pairingExpiresAt) {
			pairingStatus = "expired"
		} else {
			pairingStatus = "waiting for user"
			expiresAt = a.pairingExpiresAt.Format(time.RFC3339)
		}
	}

	return Status{
		AgentStatus:       agentStatus,
		DeviceStatus:      deviceStatus,
		ProtectionStatus:  protectionStatus,
		BackendURL:        a.backendURL,
		ServerStatus:      a.serverStatus,
		ProtectionEnabled: a.protectionEnabled,
		PairingStatus:     pairingStatus,
		DeviceCode:        a.deviceCode,
		PairingExpiresAt:  expiresAt,
		MQTTStatus:        "unknown",
		LastError:         a.lastError,
	}
}

func (a *App) refreshServerStatus() {
	ctx, cancel := a.requestContext()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.backendURL+"/health", nil)
	if err != nil {
		a.setServerStatus("offline")
		return
	}

	resp, err := a.client.Do(req)
	if err != nil {
		a.setServerStatus("offline")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		a.setServerStatus("online")
		return
	}

	a.setServerStatus("offline")
}

func (a *App) mergeAgentStatus(status Status) Status {
	agentClient, err := a.ensureAgentClient()
	if err != nil {
		return status
	}

	ctx, cancel := a.requestContext()
	defer cancel()

	agentStatus, err := agentClient.GetStatus(ctx)
	if err != nil {
		return status
	}

	status.AgentStatus = agentStatus.AgentStatus
	status.DeviceStatus = agentStatus.DeviceStatus
	status.ProtectionStatus = agentStatus.ProtectionStatus
	status.ProtectionEnabled = agentStatus.ProtectionEnabled
	status.MQTTStatus = agentStatus.MQTTStatus
	status.LastEvent = agentStatus.LastEvent
	status.LastEventAt = agentStatus.LastEventAt
	if status.LastError == "" {
		status.LastError = agentStatus.LastError
	}

	return status
}

func (a *App) refreshBackendCredentialStatus() {
	if a.hasActivePairingFlow() || a.serverStatusSnapshot() != "online" {
		return
	}

	localStore, err := a.ensureLocalStore()
	if err != nil {
		a.setLastError("Khong doc duoc credential local: " + err.Error())
		return
	}

	credentials, err := localStore.LoadCredentials()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			a.markUnlinkedIfIdle()
			return
		}
		a.setLastError("Khong doc duoc credential local: " + err.Error())
		return
	}

	ctx, cancel := a.requestContext()
	defer cancel()

	if _, err := a.backendClient.Verify(ctx, credentials.PCAgentID, credentials.AgentSecret); err != nil {
		if isBackendStatus(err, http.StatusUnauthorized, http.StatusNotFound) {
			_ = localStore.ClearCredentials()
			a.mu.Lock()
			a.pairingStatus = "not started"
			a.lastError = "Lien ket thiet bi khong con hop le. Vui long lien ket lai."
			a.mu.Unlock()
			return
		}
		a.setLastError(err.Error())
		return
	}

	a.mu.Lock()
	a.pairingStatus = "linked"
	a.lastError = ""
	a.mu.Unlock()
}

func (a *App) syncPendingAgentCredentials() error {
	a.mu.RLock()
	if a.pendingCredentials == nil {
		a.mu.RUnlock()
		return nil
	}
	credentials := *a.pendingCredentials
	a.mu.RUnlock()

	if err := a.syncAgentCredentials(&credentials); err != nil {
		return err
	}

	a.mu.Lock()
	if a.pendingCredentials != nil &&
		a.pendingCredentials.PCAgentID == credentials.PCAgentID &&
		a.pendingCredentials.AgentSecret == credentials.AgentSecret {
		a.pendingCredentials = nil
		a.pairingStatus = "linked"
	}
	a.lastError = ""
	a.mu.Unlock()

	return nil
}

func (a *App) syncAgentCredentials(credentials *agentCredentials) error {
	agentClient, err := a.ensureAgentClient()
	if err != nil {
		return err
	}

	ctx, cancel := a.requestContext()
	defer cancel()

	_, err = agentClient.SetCredentials(ctx, credentials.PCAgentID, credentials.AgentSecret)
	return err
}

func (a *App) saveLocalCredentials(credentials *agentCredentials) error {
	localStore, err := a.ensureLocalStore()
	if err != nil {
		return err
	}

	return localStore.SaveCredentials(store.Credentials{
		PCAgentID:   strings.TrimSpace(credentials.PCAgentID),
		AgentSecret: strings.TrimSpace(credentials.AgentSecret),
	})
}

func (a *App) getLocalPCAgentID() (string, error) {
	localStore, err := a.ensureLocalStore()
	if err != nil {
		return "", err
	}

	return localStore.LoadOrCreatePCAgentID()
}

func (a *App) initLocalStore() {
	localStore, err := store.New()
	if err != nil {
		a.lastError = "Khong khoi tao duoc local store: " + err.Error()
		return
	}

	a.localStore = localStore
}

func (a *App) ensureLocalStore() (*store.Store, error) {
	if a.localStore != nil {
		return a.localStore, nil
	}

	localStore, err := store.New()
	if err != nil {
		return nil, err
	}
	a.localStore = localStore
	return localStore, nil
}

func (a *App) initAgentClient() {
	a.agentClient = ipc.NewDefaultClient(5 * time.Second)
}

func (a *App) ensureAgentClient() (*ipc.Client, error) {
	if a.agentClient != nil {
		return a.agentClient, nil
	}

	client := ipc.NewDefaultClient(5 * time.Second)
	a.agentClient = client
	return client, nil
}

func (a *App) currentBackendInputs() (string, string) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.pairingSessionID, a.deviceCode
}

func (a *App) setLastError(message string) {
	a.mu.Lock()
	a.lastError = message
	a.mu.Unlock()
}

func (a *App) setServerStatus(status string) {
	a.mu.Lock()
	a.serverStatus = status
	a.mu.Unlock()
}

func (a *App) serverStatusSnapshot() string {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.serverStatus
}

func (a *App) hasActivePairingFlow() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.deviceCode != "" ||
		a.pendingCredentials != nil ||
		a.pairingStatus == "waiting for user" ||
		a.pairingStatus == "syncing local agent"
}

func (a *App) markUnlinkedIfIdle() {
	a.mu.Lock()
	defer a.mu.Unlock()

	if a.pairingStatus == "linked" {
		a.pairingStatus = "not started"
	}
}

func (a *App) getJSON(path string, out any) error {
	ctx, cancel := a.requestContext()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.backendURL+path, nil)
	if err != nil {
		return err
	}
	return a.doJSON(req, out)
}

func (a *App) postJSON(path string, body any, out any) error {
	return a.sendJSON(http.MethodPost, path, body, out)
}

func (a *App) sendJSON(method string, path string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}

	ctx, cancel := a.requestContext()
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, a.backendURL+path, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return a.doJSON(req, out)
}

func (a *App) doJSON(req *http.Request, out any) error {
	resp, err := a.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		backendErr := &backendError{StatusCode: resp.StatusCode}
		_ = json.NewDecoder(resp.Body).Decode(backendErr)
		_, _ = io.Copy(io.Discard, resp.Body)
		return backendErr
	}

	if out == nil {
		return nil
	}
	err = json.NewDecoder(resp.Body).Decode(out)
	_, _ = io.Copy(io.Discard, resp.Body)
	return err
}

func (a *App) requestContext() (context.Context, context.CancelFunc) {
	base := context.Background()
	if a.ctx != nil {
		base = a.ctx
	}
	return context.WithTimeout(base, requestTimeout)
}

func envOrDefault(fallback string, keys ...string) string {
	for _, key := range keys {
		value := strings.TrimSpace(os.Getenv(key))
		if value != "" {
			return value
		}
	}

	return fallback
}

func normalizePairingStatus(status string) string {
	switch strings.TrimSpace(status) {
	case "pending", "":
		return "waiting for user"
	case "confirmed":
		return "syncing local agent"
	default:
		return status
	}
}

func isBackendStatus(err error, statusCodes ...int) bool {
	var apiErr *backend.APIError
	if !errors.As(err, &apiErr) {
		return false
	}

	for _, statusCode := range statusCodes {
		if apiErr.StatusCode == statusCode {
			return true
		}
	}

	return false
}
