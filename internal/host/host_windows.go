package host

import (
	"fmt"
	"os/exec"
	"syscall"
)

type WindowsService struct{}

func newHostService() HostService {
	return &WindowsService{}
}

func (w *WindowsService) LockScreen() error {
	user32 := syscall.NewLazyDLL("user32.dll")
	proc := user32.NewProc("LockWorkStation")

	r1, _, err := proc.Call()
	if r1 == 0 {
		return fmt.Errorf("LockWorkStation failed: %v", err)
	}
	return nil
}

func (w *WindowsService) PlayAlarm() error {
	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-Command",
		`$p = New-Object System.Media.SoundPlayer (Resolve-Path "assets\siren_sound.wav"); $p.PlaySync()`,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("play alarm failed: %w | output: %s", err, string(out))
	}
	return nil
}

func (w *WindowsService) SetMaxVolume() error {
	script := `
Add-Type -TypeDefinition @"
using System;
using System.Runtime.InteropServices;

[Guid("5CDF2C82-841E-4546-9722-0CF74078229A"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IAudioEndpointVolume {
    int RegisterControlChangeNotify(IntPtr pNotify);
    int UnregisterControlChangeNotify(IntPtr pNotify);
    int GetChannelCount(out uint pnChannelCount);
    int SetMasterVolumeLevel(float fLevelDB, Guid pguidEventContext);
    int SetMasterVolumeLevelScalar(float fLevel, Guid pguidEventContext);
    int GetMasterVolumeLevel(out float pfLevelDB);
    int GetMasterVolumeLevelScalar(out float pfLevel);
    int SetChannelVolumeLevel(uint nChannel, float fLevelDB, Guid pguidEventContext);
    int SetChannelVolumeLevelScalar(uint nChannel, float fLevel, Guid pguidEventContext);
    int GetChannelVolumeLevel(uint nChannel, out float pfLevelDB);
    int GetChannelVolumeLevelScalar(uint nChannel, out float pfLevel);
    int SetMute([MarshalAs(UnmanagedType.Bool)] bool bMute, Guid pguidEventContext);
    int GetMute(out bool pbMute);
    int GetVolumeStepInfo(out uint pnStep, out uint pnStepCount);
    int VolumeStepUp(Guid pguidEventContext);
    int VolumeStepDown(Guid pguidEventContext);
    int QueryHardwareSupport(out uint pdwHardwareSupportMask);
    int GetVolumeRange(out float pflVolumeMindB, out float pflVolumeMaxdB, out float pflVolumeIncrementdB);
}

[Guid("A95664D2-9614-4F35-A746-DE8DB63617E6"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IMMDeviceEnumerator {
    int NotImpl1();
    int GetDefaultAudioEndpoint(int dataFlow, int role, out IMMDevice ppDevice);
}

[Guid("D666063F-1587-4E43-81F1-B948E807363F"), InterfaceType(ComInterfaceType.InterfaceIsIUnknown)]
interface IMMDevice {
    int Activate(ref Guid iid, int dwClsCtx, IntPtr pActivationParams, out IAudioEndpointVolume ppInterface);
}

[ComImport, Guid("BCDE0395-E52F-467C-8E3D-C4579291692E")]
class MMDeviceEnumerator { }

public class VolumeControl {
    public static void SetMax() {
        var enumerator = (IMMDeviceEnumerator)(new MMDeviceEnumerator());
        IMMDevice device;
        Marshal.ThrowExceptionForHR(enumerator.GetDefaultAudioEndpoint(0, 1, out device));
        Guid iid = typeof(IAudioEndpointVolume).GUID;
        IAudioEndpointVolume volume;
        Marshal.ThrowExceptionForHR(device.Activate(ref iid, 23, IntPtr.Zero, out volume));
        Marshal.ThrowExceptionForHR(volume.SetMute(false, Guid.Empty));
        Marshal.ThrowExceptionForHR(volume.SetMasterVolumeLevelScalar(1.0f, Guid.Empty));
    }
}
"@;
[VolumeControl]::SetMax()
`

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-Command",
		script,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("set max volume failed: %w | output: %s", err, string(out))
	}
	return nil
}
