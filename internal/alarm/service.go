package alarm

import (
	"fmt"
	"pc-app/internal/domain"
	"pc-app/internal/host"
	"sync"
	"time"
)

// TODO: implement alarm service
type AlarmService struct {
	host            host.HostService
	isAlarming      bool
	mt              sync.Mutex
	activeEventType string
	activeSince     time.Time
	lastFinishedAt  time.Time
	cooldown        time.Duration
}

func NewAlarmService(host host.HostService) *AlarmService {
	return &AlarmService{
		host:     host,
		cooldown: time.Second * 3,
	}
}

func (s *AlarmService) Handle(event domain.DeviceEvent) {
	switch event.EventType {
	case domain.EventMotion, domain.EventUSBRemoved:
		fmt.Printf("Phat hien su kien: %s\n", event.EventType)

	default:
		fmt.Printf("Phat hien su kien khong xac dinh: %s", event.EventType)
		return
	}

	if !s.tryStartAlarm(event) {
		fmt.Printf("Alarm is already active for event type %s, ignoring new event\n", s.activeEventType)
		return
	}
	defer s.finishAlarm()

	s.triggerAlarm(event)
}

func (s *AlarmService) tryStartAlarm(event domain.DeviceEvent) bool {
	s.mt.Lock()
	defer s.mt.Unlock()

	if s.isAlarming {
		return false
	}

	if !s.lastFinishedAt.IsZero() {
		elapsed := time.Now().Sub(s.lastFinishedAt)
		if elapsed < s.cooldown {
			fmt.Printf("Alarm is in cooldown period (%.1f seconds remaining), ignoring event\n", (s.cooldown - elapsed).Seconds())
			return false
		}

	}

	s.isAlarming = true
	s.activeEventType = event.EventType
	s.activeSince = time.Now()
	return true

}

func (s *AlarmService) finishAlarm() {
	s.mt.Lock()
	s.isAlarming = false
	s.activeEventType = ""
	s.activeSince = time.Time{}
	s.lastFinishedAt = time.Now()
	s.mt.Unlock()
}

func (s *AlarmService) triggerAlarm(event domain.DeviceEvent) {
	fmt.Printf("Kich hoat alarm cho su kien %s\n", event.EventType)
	if err := s.host.LockScreen(); err != nil {
		fmt.Printf("Loi khi phat am thanh alarm: %s\n", err)
	}
	if err := s.host.SetMaxVolume(); err != nil {
		fmt.Printf("Loi khi dat muc am luong toi da: %s\n", err)
	}
	if err := s.host.PlayAlarm(); err != nil {
		fmt.Printf("Loi khi kich hoat alarm: %s\n", err)
	}

}
