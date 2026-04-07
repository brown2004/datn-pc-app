package host

import "fmt"
// TODO: implement windows-specific host service

type WindowsService struct {
}

func newHostService() HostService {
	return &WindowsService{}
}

// LockScreen implements [HostService].
func (w *WindowsService) LockScreen() error {
	fmt.Println("Locking screen (Windows)")
	return nil
}

// PlayAlarm implements [HostService].
func (w *WindowsService) PlayAlarm() error {
	fmt.Println("Playing alarm sound (Windows)")
	return nil		
}

// SetMaxVolume implements [HostService].
func (w *WindowsService) SetMaxVolume() error {
	fmt.Println("Setting max volume (Windows)")
	return nil
}
