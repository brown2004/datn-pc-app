package alarm

import (
	"fmt"
	"pc-app/internal/domain"
	"pc-app/internal/host"
)

// TODO: implement alarm service
type AlarmService struct {
	host host.HostService
}

func NewAlarmService(host host.HostService) *AlarmService {
	return &AlarmService{
		host: host,
	}
}

func (s *AlarmService) Handle(event domain.DeviceEvent) {
	switch event.EventType {
	case domain.EventMotion:
		// Handle motion event
		s.handleMotionEvent(event)
	case domain.EventUSBRemoved:
		// Handle USB removed event
		s.handleUSBRemovedEvent(event)

	default:
		fmt.Println("Phat hien su kien khong xac dinh: %s", event.EventType)
	}
}

func (s *AlarmService) handleUSBRemovedEvent(event domain.DeviceEvent) {
	s.host.SetMaxVolume()
	s.host.PlayAlarm()
	s.host.LockScreen()
}

func (s *AlarmService) handleMotionEvent(event domain.DeviceEvent) {
	s.host.SetMaxVolume()
	s.host.PlayAlarm()
	s.host.LockScreen()
}
