package device

import "pc-app/internal/domain"

// TODO: implement device monitor
type Monitor struct {
	out chan domain.DeviceEvent //out la mot cai channel chua ten cua event
}
