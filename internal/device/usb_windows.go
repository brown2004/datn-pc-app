package device

import (
	"fmt"
	"os/exec"
	"strings"
)

func isUSBTargetPresent(target USBTarget) (bool, error) {
	pattern := fmt.Sprintf("USB\\VID_%s&PID_%s", strings.ToUpper(target.VendorID), strings.ToUpper(target.ProductID))

	script := fmt.Sprintf(`
$pattern = "%s"
$devices = Get-CimInstance Win32_PnPEntity | Where-Object { $_.PNPDeviceID -like "$pattern*" }
if ($devices) { Write-Output "present" } else { Write-Output "absent" }
`, pattern)

	cmd := exec.Command(
		"powershell",
		"-NoProfile",
		"-ExecutionPolicy", "Bypass",
		"-Command",
		script,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return false, fmt.Errorf("windows usb detect failed: %w | output: %s", err, string(out))
	}

	raw := strings.TrimSpace(strings.ToLower(string(out)))
	fmt.Println("[DEVICE][WINDOWS] detect output:", raw)

	return raw == "present", nil
}