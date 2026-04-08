package agent

// TODO: implement agent run
import (
	"pc-app/internal/alarm"
	"pc-app/internal/device"
	"pc-app/internal/domain"
	"pc-app/internal/host"
	"time"
)

func Run() {
	eventChan := make(chan domain.DeviceEvent)

	// khoi tao doi tuong monitor va alarm

	usbTarget := device.USBTarget{
		VendorID:  "413c", // Thay doi theo USB can theo doi
		ProductID: "301a", // Thay doi theo USB can theo doi
	}
	monitor := device.NewMonitor(usbTarget, time.Millisecond*500, eventChan)
	if err := monitor.Start(); err != nil {
		panic("Failed to start device monitor: " + err.Error())
	}

	hostService := host.New()
	alarmService := alarm.NewAlarmService(hostService)

	monitor.Start()

	for event := range eventChan {
		alarmService.Handle(event)

	}
}
