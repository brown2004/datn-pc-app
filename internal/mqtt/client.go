package mqtt

// TODO: implement mqtt client
import (
	"encoding/json"
	"fmt"
	"os"
	"pc-app/internal/domain"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTClient struct {
	client   mqtt.Client
	deviceID string
	enabled  bool // ????
}

func NewClientFromEnv() (*MQTTClient, error) {

	//doc env
	broker := os.Getenv("PCAPP_MQTT_BROKER")
	deviceID := os.Getenv("PCAPP_DEVICE_ID")

	if broker == "" || deviceID == "" {
		return &MQTTClient{enabled: false}, nil
	}

	clientID := os.Getenv("PCAPP_MQTT_CLIENT_ID")
	if clientID == "" {
		clientID = "pc-app-" + deviceID
	}

	// options cho ket noi mqtt
	opts := mqtt.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectRetry(true)
	opts.SetConnectRetryInterval(3 * time.Second)

	if username := os.Getenv("PCAPP_MQTT_USERNAME"); username != "" {
		opts.SetUsername(username)
		opts.SetPassword(os.Getenv("PCAPP_MQTT_PASSWORD"))
	}

	willPayload, err := buildUnexpectedDisconnectPayload(deviceID)
	if err != nil {
		return nil, fmt.Errorf("build mqtt will payload: %w", err)
	}

	// thiet lap LWT tai topic cap nhat status, qos = 1
	opts.SetWill(StatusTopic(deviceID), string(willPayload), 1, true)

	// khoi tao doi tuong mqtt client
	client := mqtt.NewClient(opts)

	token := client.Connect()
	if !token.WaitTimeout(10 * time.Second) {
		return nil, fmt.Errorf("mqtt connect timeout")
	}
	if err := token.Error(); err != nil {
		return nil, fmt.Errorf("mqtt connect failed: %w", err)
	}

	mqttClient := &MQTTClient{
		client:   client,
		deviceID: deviceID,
		enabled:  true,
	}

	// cap nhat trang thai active den mqtt broker
	if err := mqttClient.PublishStatus(StatusActive, ReasonConnected); err != nil {
		client.Disconnect(250)
		return nil, err
	}

	return mqttClient, nil
}

func (cli *MQTTClient) Enabled() bool {
	return cli != nil && cli.enabled
}

func (cli *MQTTClient) PublishStatus(status, reason string) error {
	if !cli.Enabled() {
		return nil
	}
	payload, err := buildStatusPayload(cli.deviceID, status, reason)
	if err != nil {
		return fmt.Errorf("mashal status payload: %w", err)
	}

	return cli.publish(StatusTopic(cli.deviceID), payload)

}

func (cli *MQTTClient) PublishAlert(event domain.DeviceEvent) error {
	if !cli.Enabled() {
		return nil
	}

	payload, err := json.Marshal(AlertPayload{
		DeviceID:  cli.deviceID,
		EventType: event.EventType,
		Message:   fmt.Sprintf("detected event %s", event.EventType),
		Timestamp: time.Unix(event.Timestamp, 0),
	})
	if err != nil {
		return fmt.Errorf("marshal alert payload: %w", err)
	}

	return cli.publish(AlertTopic(cli.deviceID), payload)
}

func (cli *MQTTClient) Close() error {
	if !cli.Enabled() {
		return nil
	}
	if err := cli.PublishStatus(StatusInactive, ReasonGracefulShutdown); err != nil {
		return nil
	}
	cli.client.Disconnect(250)
	return nil

}

func (cli *MQTTClient) publish(topic string, payload []byte) error {
	token := cli.client.Publish(topic, 1, false, payload)
	if !token.WaitTimeout(5 * time.Second) {
		return fmt.Errorf("mqtt publish timeout for topic %s", topic)
	}
	if err := token.Error(); err != nil {
		return fmt.Errorf("mqtt publish failed for topic %s: %w", topic, err)
	}
	return nil
}
