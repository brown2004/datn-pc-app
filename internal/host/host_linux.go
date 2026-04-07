package host

import (
	"fmt"
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
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("lock screen failed: %w | output: %s", err, string(out))
	}
	return nil
}

// PlayAlarm implements [HostService].
func (l *LinuxService) PlayAlarm() error {
	cmd := exec.Command("aplay", "assets/siren_sound.wav")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("play alarm failed: %w | output: %s", err, string(out))
	}
	return nil
}

// SetMaxVolume implements [HostService].
func (l *LinuxService) SetMaxVolume() error {
	cmd := exec.Command("amixer", "sset", "Master", "100%")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("set max volume failed: %w | output: %s", err, string(out))
	}
	return nil
}
