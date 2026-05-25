package device

import "testing"

func TestIsDATNDeviceID(t *testing.T) {
	tests := []struct {
		name     string
		deviceID string
		want     bool
	}{
		{
			name:     "firmware DATN id",
			deviceID: "DATN-F103-1A5A5CCA0000000000048CD2",
			want:     true,
		},
		{
			name:     "bare MCU UID is not enough",
			deviceID: "1A5A5CCA0000000000048CD2",
			want:     false,
		},
		{
			name:     "PC label is not a device firmware id",
			deviceID: "pc-vht-001",
			want:     false,
		},
		{
			name:     "empty device id",
			deviceID: "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDATNDeviceID(tt.deviceID); got != tt.want {
				t.Fatalf("isDATNDeviceID(%q) = %v, want %v", tt.deviceID, got, tt.want)
			}
		})
	}
}
