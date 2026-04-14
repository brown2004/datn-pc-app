package device

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func isUSBTargetPresent(target USBTarget) (bool, error) {
	paths, err := filepath.Glob("/sys/bus/usb/devices/*")
	if err != nil {
		return false, fmt.Errorf("failed to list USB devices: %w", err)
	}

	for _, path := range paths {
		vendorBytes, err1 := os.ReadFile(filepath.Join(path, "idVendor"))
		if err1 != nil {
			continue
		}
		productBytes, err2 := os.ReadFile(filepath.Join(path, "idProduct"))
		if err2 != nil {
			continue
		}

		vendor := strings.TrimSpace(string(vendorBytes))
		product := strings.TrimSpace(string(productBytes))

		if strings.EqualFold(vendor, target.VendorID) &&
			strings.EqualFold(product, target.ProductID) {
			return true, nil
		}
	}

	return false, nil
}