package host

// TODO: implement host service
type HostService interface {
	SetMaxVolume() error
	PlayAlarm() error
	LockScreen() error
}
