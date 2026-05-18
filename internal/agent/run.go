package agent

// TODO: implement agent run
import (
	"fmt"
	"pc-app/internal/alarm"
	"pc-app/internal/device"
	"pc-app/internal/domain"
	"pc-app/internal/host"
	"pc-app/internal/mqtt"
)

func Run() {
	eventChan := make(chan domain.DeviceEvent)

	expectedDeviceID := "DATN-F103-1A5A5CCA0000000000048CD2"

	serialMonitor := device.NewSerialMonitor(expectedDeviceID, eventChan)
	serialMonitor.Start()

	// usbTarget := device.USBTarget{
	// 	VendorID:  "093a",
	// 	ProductID: "2510",
	// }

	// monitor := device.NewMonitor(usbTarget, time.Millisecond*500, eventChan)
	// if err := monitor.Start(); err != nil {
	// 	panic("Failed to start device monitor: " + err.Error())
	// }

	//host: xu ly can thiep sau cac phan ung tuy thuoc vao he dieu hanh
	hostService := host.New()

	//alarm: logic xu ly su kien canh bao
	alarmService := alarm.NewAlarmService(hostService)

	//mqtt: khoi tao mqtt client va defer don dep
	mqttClient, err := mqtt.NewClientFromEnv()
	if err != nil {
		fmt.Printf("[MQTT] init failed: %v\n", err)
	}

	if mqttClient != nil && mqttClient.Enabled() {
		defer func() {
			if err := mqttClient.Close(); err != nil {
				fmt.Printf("[MQTT] close failed: %v\n", err)
			}
		}()
	} else {
		fmt.Println("[MQTT] disabled, set PCAPP_MQTT_BROKER and PCAPP_DEVICE_ID to enable publishing")
	}

	// consumer eventChan de xu ly canh bao
	for event := range eventChan {
		alarmService.Handle(event)

		if mqttClient != nil {
			if err := mqttClient.PublishAlert(event); err != nil {
				fmt.Printf("[MQTT] publish alert failed: %v\n", err)
			}
		}

	}
}
