package mqtt

import "time"

const (
	StatusActive   = "active"
	StatusInactive = "inactive"

	ReasonConnected            = "connected"
	ReasonGracefulShutdown     = "graceful_shutdown"
	ReasonUnexpectedDisconnect = "unexpected_disconnect"
)

type StatusPayload struct {
	DeviceID  string    `json:"device_id"`
	Status    string    `json:"status"`
	Reason    string    `json:"reason"`
	Timestamp time.Time `json:"timestamp"`
}

type AlertPayload struct {
	DeviceID  string    `json:"device_id"`
	EventType string    `json:"event_type"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}
