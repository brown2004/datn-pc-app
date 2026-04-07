package device

import (
	"fmt"
	"pc-app/internal/domain"
	"time"
)

// TODO: implement device monitor
type Monitor struct {
	out chan<- domain.DeviceEvent //out la mot cai channel chua ten cua event

}

func NewMonitor(out chan<- domain.DeviceEvent) *Monitor {
	return &Monitor{
		out: out,
	}
}

func (m *Monitor) Start() {
	go func() {
		for {
			time.Sleep(3 * time.Second)

			event := domain.DeviceEvent{
				EventType: domain.EventMotion,
				Timestamp: time.Now().Unix(),
			}

			fmt.Println("Monitor: detected motion event %s", event.EventType)
			m.out <- event
			time.Sleep(3 * time.Second)

			event = domain.DeviceEvent{
				EventType: domain.EventUSBRemoved,
				Timestamp: time.Now().Unix(),
			}

			fmt.Println("Monitor: detected USB removed event %s", event.EventType)
			m.out <- event

		}
	}()
}
