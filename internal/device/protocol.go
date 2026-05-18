package device

import (
	"strings"
)

type DeviceMessage struct {
	Raw      string
	Kind     string
	Fields   map[string]string
	DeviceID string
}

func ParseDeviceLine(line string) (DeviceMessage, bool) {
	line = strings.TrimSpace(line)
	if line == "" {
		return DeviceMessage{}, false
	}

	parts := strings.Split(line, ";")
	if len(parts) == 0 {
		return DeviceMessage{}, false
	}

	kind := parts[0]

	if kind != "DATN_HELLO" && kind != "DATN_DATA" && kind != "DATN_EVENT" {
		return DeviceMessage{}, false
	}

	fields := make(map[string]string)

	for _, part := range parts[1:] {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		fields[key] = value
	}

	return DeviceMessage{
		Raw:      line,
		Kind:     kind,
		Fields:   fields,
		DeviceID: fields["device_id"],
	}, true
}
