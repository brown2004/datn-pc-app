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
		VendorID:  "24AE",
		ProductID: "2013",// Thay doi theo USB can theo doi
	}
	monitor := device.NewMonitor(usbTarget, time.Millisecond*500, eventChan)
	if err := monitor.Start(); err != nil {
		panic("Failed to start device monitor: " + err.Error())
	}

	hostService := host.New()
	alarmService := alarm.NewAlarmService(hostService)

	

	for event := range eventChan {
		alarmService.Handle(event)

	}
}
