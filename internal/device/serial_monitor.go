package device

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	"pc-app/internal/domain"

	"go.bug.st/serial"
	"go.bug.st/serial/enumerator"
)

const defaultBaudRate = 115200

type SerialMonitor struct {
	expectedDeviceID string
	eventChan        chan<- domain.DeviceEvent
}

func NewSerialMonitor(expectedDeviceID string, eventChan chan<- domain.DeviceEvent) *SerialMonitor {
	return &SerialMonitor{
		expectedDeviceID: expectedDeviceID,
		eventChan:        eventChan,
	}
}

func (m *SerialMonitor) Start() {
	go m.loop()
}

func (m *SerialMonitor) loop() {
	for {
		portName, connectedDeviceID, err := m.findDevicePort()
		if err != nil {
			fmt.Printf("[SERIAL] device not found: %v\n", err)
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf("[SERIAL] connected to STM32 device on %s, device_id=%s\n", portName, connectedDeviceID)
		m.eventChan <- domain.DeviceEvent{
			EventType: domain.EventUSBConnected,
			Timestamp: time.Now().Unix(),
			DeviceID:  connectedDeviceID,
		}

		if err := m.readLoop(portName); err != nil {
			fmt.Printf("[SERIAL] disconnected or read failed: %v\n", err)

			m.eventChan <- domain.DeviceEvent{
				EventType: domain.EventUSBRemoved,
				Timestamp: time.Now().Unix(),
				DeviceID:  connectedDeviceID,
			}

			time.Sleep(2 * time.Second)
		}
	}
}

func getCandidatePorts() ([]string, error) {
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		return nil, err
	}

	stm32Ports := make([]string, 0)
	usbSerialPorts := make([]string, 0)

	for _, port := range ports {
		if !port.IsUSB {
			fmt.Printf("[SERIAL] skip non-usb port: %s\n", port.Name)
			continue
		}

		fmt.Printf(
			"[SERIAL] usb port: name=%s vid=%s pid=%s serial=%s product=%s\n",
			port.Name,
			port.VID,
			port.PID,
			port.SerialNumber,
			port.Product,
		)

		if strings.EqualFold(port.VID, "0483") && strings.EqualFold(port.PID, "5740") {
			stm32Ports = append(stm32Ports, port.Name)
			continue
		}

		usbSerialPorts = append(usbSerialPorts, port.Name)
	}

	if len(stm32Ports) > 0 {
		return stm32Ports, nil
	}

	return usbSerialPorts, nil
}

func (m *SerialMonitor) findDevicePort() (string, string, error) {
	ports, err := getCandidatePorts()
	if err != nil {
		return "", "", err
	}

	if len(ports) == 0 {
		return "", "", fmt.Errorf("no candidate USB serial ports found")
	}

	fmt.Printf("[SERIAL] candidate ports: %v\n", ports)

	for _, portName := range ports {
		fmt.Printf("[SERIAL] trying candidate port: %s\n", portName)

		port, err := openSerialPort(portName)
		if err != nil {
			fmt.Printf("[SERIAL] open %s failed: %v\n", portName, err)
			continue
		}

		msg, ok := waitHello(port, 5*time.Second)
		_ = port.Close()

		if !ok {
			fmt.Printf("[SERIAL] no DATN_HELLO from %s\n", portName)
			continue
		}

		if m.expectedDeviceID != "" && msg.DeviceID != m.expectedDeviceID {
			fmt.Printf("[SERIAL] ignore %s: device_id mismatch: expected=%s actual=%s\n", portName, m.expectedDeviceID, msg.DeviceID)
			continue
		}

		fmt.Printf("[SERIAL] DATN_HELLO OK: device_id=%s\n", msg.DeviceID)
		return portName, msg.DeviceID, nil
	}

	return "", "", fmt.Errorf("no DATN STM32 device found")
}

func (m *SerialMonitor) readLoop(portName string) error {
	port, err := openSerialPort(portName)
	if err != nil {
		return err
	}
	defer port.Close()

	reader := bufio.NewReader(port)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		msg, ok := ParseDeviceLine(line)
		if !ok {
			continue
		}

		if m.expectedDeviceID != "" && msg.DeviceID != m.expectedDeviceID {
			fmt.Printf("[SERIAL] ignore message from unknown device_id: expected=%s actual=%s\n", m.expectedDeviceID, msg.DeviceID)
			continue
		}

		m.handleMessage(msg)
	}
}

func (m *SerialMonitor) handleMessage(msg DeviceMessage) {
	switch msg.Kind {
	case "DATN_HELLO":
		fmt.Printf("[SERIAL] hello device_id=%s mpu=%s fw=%s\n",
			msg.DeviceID,
			msg.Fields["mpu"],
			msg.Fields["fw"],
		)

	case "DATN_DATA":
		fmt.Printf("[SERIAL] data device_id=%s score_mg=%s\n",
			msg.DeviceID,
			msg.Fields["score_mg"],
		)

	case "DATN_EVENT":
		eventType := msg.Fields["type"]
		if eventType != "MOTION" {
			return
		}

		scoreMG, err := strconv.Atoi(msg.Fields["score_mg"])
		if err != nil {
			scoreMG = 0
		}

		fmt.Printf("[SERIAL] motion event device_id=%s score_mg=%d\n", msg.DeviceID, scoreMG)

		m.eventChan <- domain.DeviceEvent{
			EventType: domain.EventMotion,
			Timestamp: time.Now().Unix(),
			DeviceID:  msg.DeviceID,
			ScoreMG:   scoreMG,
		}
	}
}

func openSerialPort(portName string) (serial.Port, error) {
	mode := &serial.Mode{
		BaudRate: defaultBaudRate,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	port, err := serial.Open(portName, mode)
	if err != nil {

		return nil, err
	}

	_ = port.SetReadTimeout(2 * time.Second)

	// Quan trọng: bật DTR để firmware STM32 set g_cdc_port_open = 1.
	_ = port.SetDTR(true)
	_ = port.SetRTS(true)

	time.Sleep(300 * time.Millisecond)
	fmt.Printf("Da mo cong %s\n", portName)
	return port, nil
}

func waitHello(port serial.Port, timeout time.Duration) (DeviceMessage, bool) {
	reader := bufio.NewReader(port)
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		line, err := reader.ReadString('\n')
		if err != nil {
			continue
		}

		fmt.Printf("[SERIAL] read line while waiting hello: %s", line)

		msg, ok := ParseDeviceLine(line)
		if !ok {
			continue
		}

		if msg.Kind == "DATN_HELLO" && msg.DeviceID != "" {
			return msg, true
		}
	}

	return DeviceMessage{}, false
}
