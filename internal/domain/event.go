package domain

// TODO: implement event types
const (
	EventMotion     string = "motion_alert"
	EventUSBRemoved string = "usb_removed"
)

type DeviceEvent struct {
	EventType string
	Timestamp int64

	DeviceID string
	ScoreMG  int
}
