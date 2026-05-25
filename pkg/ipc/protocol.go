package ipc

const (
	DefaultAddr    = "127.0.0.1:17831"
	DefaultBaseURL = "http://" + DefaultAddr
)

type Status struct {
	AgentStatus       string `json:"agent_status"`
	DeviceStatus      string `json:"device_status"`
	ProtectionStatus  string `json:"protection_status"`
	BackendURL        string `json:"backend_url"`
	ProtectionEnabled bool   `json:"protection_enabled"`
	MQTTStatus        string `json:"mqtt_status"`
	LastError         string `json:"last_error"`
	LastEvent         string `json:"last_event"`
	LastEventAt       string `json:"last_event_at"`
}

type ProtectionRequest struct {
	Enabled bool `json:"enabled"`
}

type CredentialsRequest struct {
	PCAgentID   string `json:"pc_agent_id"`
	AgentSecret string `json:"agent_secret"`
}

type Identity struct {
	PCAgentID string `json:"pc_agent_id"`
}

type CredentialsStatus struct {
	HasCredentials bool   `json:"has_credentials"`
	Linked         bool   `json:"linked"`
	Error          string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
