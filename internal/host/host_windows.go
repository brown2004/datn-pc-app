package host

// TODO: implement windows-specific host service

type WindowsService struct {
}

func New() HostService {
	return &WindowsService{}
}

// LockScreen implements [HostService].
func (w *WindowsService) LockScreen() error {
	panic("unimplemented")
}

// PlayAlarm implements [HostService].
func (w *WindowsService) PlayAlarm() error {
	panic("unimplemented")
}

// SetMaxVolume implements [HostService].
func (w *WindowsService) SetMaxVolume() error {
	panic("unimplemented")
}
