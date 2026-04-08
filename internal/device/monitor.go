package device

import (
	"fmt"
	"pc-app/internal/domain"
	"time"
)

// TODO: implement device monitor
type Monitor struct {
	target     USBTarget
	interval   time.Duration
	out        chan<- domain.DeviceEvent //out la mot cai channel chua ten cua event
	wasPresent bool
}

func NewMonitor(target USBTarget, interval time.Duration, out chan<- domain.DeviceEvent) *Monitor {
	return &Monitor{
		target:   target,
		interval: interval,
		out:      out,
	}
}

func (m *Monitor) Start() error {
	present, err := IsUSBTargetPresent(m.target)
	if err != nil {
		return err
	}

	m.wasPresent = present

	go m.loop()

	return nil
}

func (m *Monitor) loop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for range ticker.C {
		present, err := IsUSBTargetPresent(m.target)
		if err != nil {
			fmt.Println("Error checking USB target:", err)
			continue
		}

		if m.wasPresent && !present {
			m.out <- domain.DeviceEvent{
				EventType: domain.EventUSBRemoved,
				Timestamp: time.Now().Unix(),
			}
		}
		if !m.wasPresent && present {
			fmt.Println("USB target connected")
		}

		m.wasPresent = present
	}
}
