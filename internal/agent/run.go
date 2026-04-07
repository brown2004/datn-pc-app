package agent

// TODO: implement agent run
import (
	"pc-app/internal/alarm"
	"pc-app/internal/device"
	"pc-app/internal/domain"
	"pc-app/internal/host"
)

func Run() {
	eventChan := make(chan domain.DeviceEvent)

	// khoi tao doi tuong monitor va alarm
	monitor := device.NewMonitor(eventChan)
	hostService := host.New()
	alarmService := alarm.NewAlarmService(hostService)

	monitor.Start()

	for event := range eventChan {
		alarmService.Handle(event)

	}
}
