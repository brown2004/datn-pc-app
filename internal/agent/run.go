package agent

import (
	"context"
	"fmt"

	"pc-app/internal/alarm"
	"pc-app/internal/backend"
	"pc-app/internal/config"
	"pc-app/internal/device"
	"pc-app/internal/domain"
	"pc-app/internal/host"
	"pc-app/internal/mqtt"
	"pc-app/internal/store"
	"pc-app/pkg/ipc"
)

func Run() {
	ctx := context.Background()
	cfg := config.Load()

	localStore, err := store.New()
	if err != nil {
		panic("failed to init local store: " + err.Error())
	}

	backendClient := backend.NewClient(cfg.BackendURL, cfg.RequestTimeout)
	app, err := NewApp(cfg, backendClient, localStore)
	if err != nil {
		panic("failed to init agent app: " + err.Error())
	}

	ipcServer := ipc.NewServer(ipc.DefaultAddr, app)
	go func() {
		if err := ipcServer.ListenAndServe(ctx); err != nil {
			fmt.Printf("[IPC] server stopped: %v\n", err)
			app.setLastError("IPC server stopped: " + err.Error())
		}
	}()
	fmt.Printf("[IPC] listening on %s\n", ipc.DefaultAddr)

	eventChan := make(chan domain.DeviceEvent)
	serialMonitor := device.NewSerialMonitor(eventChan)
	serialMonitor.Start()

	hostService := host.New()
	alarmService := alarm.NewAlarmService(hostService)

	mqttClient, err := mqtt.NewClientFromEnv()
	if err != nil {
		fmt.Printf("[MQTT] init failed: %v\n", err)
		app.MarkMQTTStatus("error")
	} else if mqttClient != nil && mqttClient.Enabled() {
		app.MarkMQTTStatus("connected")
		defer func() {
			app.MarkMQTTStatus("disconnected")
			if err := mqttClient.Close(); err != nil {
				fmt.Printf("[MQTT] close failed: %v\n", err)
			}
		}()
	} else {
		app.MarkMQTTStatus("not configured")
		fmt.Println("[MQTT] disabled, set PCAPP_MQTT_BROKER and PCAPP_DEVICE_ID to enable publishing")
	}

	for event := range eventChan {
		app.HandleDeviceEvent(event)
		if event.EventType == domain.EventUSBConnected {
			continue
		}
		if !app.ProtectionEnabled() {
			fmt.Printf("[PROTECTION] disabled, ignore event %s from device_id=%s\n", event.EventType, event.DeviceID)
			continue
		}

		alarmService.Handle(event)
		if mqttClient != nil {
			if err := mqttClient.PublishAlert(event); err != nil {
				fmt.Printf("[MQTT] publish alert failed: %v\n", err)
				app.MarkMQTTStatus("error")
			}
		}
	}
}
