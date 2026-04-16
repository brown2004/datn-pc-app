package mqtt

// TODO: implement mqtt will message
import (
	"encoding/json"
	"time"
)

func buildStatusPayload(deviceId, status, reason string) ([]byte, error) {
	payload := StatusPayload{
		DeviceID:  deviceId,
		Status:    status,
		Reason:    reason,
		Timestamp: time.Now(),
	}
	return json.Marshal(payload)
} // nen cho vao file payload.go

func buildUnexpectedDisconnectPayload(deviceId string) ([]byte, error) {
	return buildStatusPayload(deviceId, StatusInactive, ReasonUnexpectedDisconnect)
}
