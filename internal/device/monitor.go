package device

import (
	"fmt"
	"pc-app/internal/domain"
	"time"
)

type Monitor struct {
	target     USBTarget
	interval   time.Duration
	out        chan<- domain.DeviceEvent
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
	present, err := isUSBTargetPresent(m.target)
	if err != nil {
		return err
	}

	fmt.Printf("[DEVICE] initial target state: vendor=%s product=%s present=%v\n",
		m.target.VendorID, m.target.ProductID, present)

	m.wasPresent = present
	go m.loop()
	return nil
}

func (m *Monitor) loop() {
	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	for range ticker.C {
		present, err := isUSBTargetPresent(m.target)
		if err != nil {
			fmt.Println("[DEVICE] error checking USB target:", err)
			continue
		}

		fmt.Printf("[DEVICE] poll target state: present=%v, wasPresent=%v\n", present, m.wasPresent)

		if m.wasPresent && !present {
			fmt.Println("[DEVICE] target USB removed -> emit event")
			m.out <- domain.DeviceEvent{
				EventType: domain.EventUSBRemoved,
				Timestamp: time.Now().Unix(),
			}
		}

		if !m.wasPresent && present {
			fmt.Println("[DEVICE] target USB connected again")
		}

		m.wasPresent = present
	}
}