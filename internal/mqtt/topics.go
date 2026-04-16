package mqtt

const baseTopic = "pcapp"

func StatusTopic(deviceID string) string {
	return baseTopic + "/status/" + deviceID
}

func AlertTopic(deviceID string) string {
	return baseTopic + "/alert/" + deviceID
}
