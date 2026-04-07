package host

import (
	"os/exec"
)

// TODO: implement linux-specific host service
type LinuxService struct {
}

func New() HostService {
	return &LinuxService{}
}

// LockScreen implements [HostService].
func (l *LinuxService) LockScreen() error {
	cmd := exec.Command("loginctl", "lock-session")
	return cmd.Run()
}

// PlayAlarm implements [HostService].
func (l *LinuxService) PlayAlarm() error {
	cmd := exec.Command("aplay", "assets/siren_sound.wav")
	return cmd.Run()
}

// SetMaxVolume implements [HostService].
func (l *LinuxService) SetMaxVolume() error {
	cmd := exec.Command("amixer", "sset", "Master", "100%")
	return cmd.Run()
}
