package agent

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"pc-app/internal/backend"
	"pc-app/internal/config"
	"pc-app/internal/domain"
	"pc-app/internal/store"
	"pc-app/pkg/ipc"
)

type App struct {
	mu sync.RWMutex

	cfg     config.Config
	backend *backend.Client
	store   *store.Store

	credentials       *store.Credentials
	protectionEnabled bool
	lastError         string

	deviceConnected bool
	deviceID        string
	lastEvent       string
	lastEventAt     time.Time
	mqttStatus      string
}

func NewApp(cfg config.Config, backendClient *backend.Client, localStore *store.Store) (*App, error) {
	app := &App{
		cfg:        cfg,
		backend:    backendClient,
		store:      localStore,
		mqttStatus: "not configured",
	}

	if credentials, err := localStore.LoadCredentials(); err == nil {
		app.credentials = credentials
		app.logf("loaded local credentials pc_agent_id=%s", credentials.PCAgentID)
	} else if !errors.Is(err, os.ErrNotExist) {
		app.lastError = err.Error()
		app.logErrorf("%s", err.Error())
	} else {
		app.logf("no local credentials found")
	}

	return app, nil
}

func (a *App) Status(ctx context.Context) (ipc.Status, error) {
	return a.snapshot(), nil
}

func (a *App) Identity(ctx context.Context) (ipc.Identity, error) {
	pcAgentID, err := a.store.LoadPCAgentID()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return ipc.Identity{}, nil
		}
		return ipc.Identity{}, err
	}

	return ipc.Identity{PCAgentID: pcAgentID}, nil
}

func (a *App) CredentialsStatus(ctx context.Context) (ipc.CredentialsStatus, error) {
	credentials := a.currentCredentials()
	if credentials == nil {
		a.logf("credentials status requested: no local credentials")
		return ipc.CredentialsStatus{}, nil
	}

	a.logf("verifying credentials with backend pc_agent_id=%s", credentials.PCAgentID)
	if _, err := a.backend.Verify(ctx, credentials.PCAgentID, credentials.AgentSecret); err != nil {
		if isBackendStatus(err, http.StatusUnauthorized, http.StatusNotFound) {
			a.clearCredentials("Lien ket thiet bi khong con hop le. Vui long lien ket lai.")
			return ipc.CredentialsStatus{}, nil
		}

		a.setLastError(err.Error())
		return ipc.CredentialsStatus{
			HasCredentials: true,
			Linked:         false,
			Error:          err.Error(),
		}, nil
	}

	a.mu.Lock()
	a.lastError = ""
	a.mu.Unlock()
	a.logf("credentials verified with backend pc_agent_id=%s", credentials.PCAgentID)

	return ipc.CredentialsStatus{
		HasCredentials: true,
		Linked:         true,
	}, nil
}

func (a *App) SetProtection(ctx context.Context, enabled bool) (ipc.Status, error) {
	a.logf("protection update requested enabled=%t", enabled)

	a.mu.Lock()
	a.protectionEnabled = enabled
	a.lastError = ""
	a.mu.Unlock()

	credentials := a.currentCredentials()
	if credentials == nil {
		a.logf("protection updated locally only enabled=%t; missing credentials", enabled)
		return a.snapshot(), nil
	}

	a.logf("syncing protection to backend pc_agent_id=%s enabled=%t", credentials.PCAgentID, enabled)
	_, err := a.backend.SetProtection(ctx, credentials.PCAgentID, credentials.AgentSecret, enabled)
	if err != nil {
		if isBackendStatus(err, http.StatusUnauthorized, http.StatusNotFound) {
			a.clearCredentials("Lien ket thiet bi khong con hop le. Vui long lien ket lai.")
			return a.snapshot(), nil
		}
		a.setLastError(err.Error())
		return a.snapshot(), nil
	}
	a.logf("protection synced to backend pc_agent_id=%s enabled=%t", credentials.PCAgentID, enabled)

	return a.snapshot(), nil
}

func (a *App) SetCredentials(ctx context.Context, pcAgentID string, agentSecret string) (ipc.Status, error) {
	a.logf("credentials sync received pc_agent_id=%s", strings.TrimSpace(pcAgentID))

	credentials := store.Credentials{
		PCAgentID:   strings.TrimSpace(pcAgentID),
		AgentSecret: strings.TrimSpace(agentSecret),
	}
	if credentials.PCAgentID == "" || credentials.AgentSecret == "" {
		a.setLastError("credentials khong hop le")
		return a.snapshot(), nil
	}
	if err := a.store.SaveCredentials(credentials); err != nil {
		a.setLastError(err.Error())
		return a.snapshot(), nil
	}

	a.mu.Lock()
	a.credentials = &credentials
	a.lastError = ""
	a.mu.Unlock()
	a.logf("credentials saved pc_agent_id=%s", credentials.PCAgentID)

	return a.snapshot(), nil
}

func (a *App) HandleDeviceEvent(event domain.DeviceEvent) {
	a.mu.Lock()
	a.deviceConnected = event.EventType != domain.EventUSBRemoved
	a.deviceID = event.DeviceID
	a.lastEvent = event.EventType
	a.lastEventAt = time.Unix(event.Timestamp, 0).UTC()
	a.mu.Unlock()
	a.logf("device event type=%s device_id=%s score_mg=%d", event.EventType, event.DeviceID, event.ScoreMG)
}

func (a *App) MarkMQTTStatus(status string) {
	a.mu.Lock()
	previous := a.mqttStatus
	a.mqttStatus = status
	a.mu.Unlock()
	if status != previous {
		a.logf("mqtt status changed %q -> %q", previous, status)
	}
}

func (a *App) ProtectionEnabled() bool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.protectionEnabled
}

func (a *App) snapshot() ipc.Status {
	a.mu.RLock()
	defer a.mu.RUnlock()

	protectionStatus := "disabled"
	if a.protectionEnabled {
		protectionStatus = "enabled"
	}

	deviceStatus := "disconnected"
	if a.deviceConnected {
		deviceStatus = "connected"
	} else if a.deviceID != "" {
		deviceStatus = "removed"
	}

	lastEventAt := ""
	if !a.lastEventAt.IsZero() {
		lastEventAt = a.lastEventAt.Format(time.RFC3339)
	}

	return ipc.Status{
		AgentStatus:       "online",
		DeviceStatus:      deviceStatus,
		ProtectionStatus:  protectionStatus,
		BackendURL:        a.cfg.BackendURL,
		ProtectionEnabled: a.protectionEnabled,
		MQTTStatus:        a.mqttStatus,
		LastError:         a.lastError,
		LastEvent:         a.lastEvent,
		LastEventAt:       lastEventAt,
	}
}

func (a *App) currentCredentials() *store.Credentials {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.credentials == nil {
		return nil
	}

	copy := *a.credentials
	return &copy
}

func (a *App) setLastError(message string) {
	a.mu.Lock()
	previous := a.lastError
	a.lastError = message
	a.mu.Unlock()
	if strings.TrimSpace(message) != "" && message != previous {
		a.logErrorf("%s", message)
	}
}

func (a *App) clearCredentials(message string) {
	if err := a.store.ClearCredentials(); err != nil {
		message = fmt.Sprintf("%s Khong the xoa credentials local: %s", message, err)
	}

	a.mu.Lock()
	a.credentials = nil
	a.lastError = message
	a.mu.Unlock()
	a.logErrorf("%s", message)
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

func (a *App) logf(format string, args ...any) {
	fmt.Printf("[AGENT] "+format+"\n", args...)
}

func (a *App) logErrorf(format string, args ...any) {
	fmt.Printf("[AGENT][ERROR] "+format+"\n", args...)
}
