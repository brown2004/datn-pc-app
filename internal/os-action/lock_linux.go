//go:build linux

package osaction

import (
	"fmt"
	"os/exec"
)

// LockScreen thực thi lệnh khóa màn hình qua bash shell
func LockScreen() {
	fmt.Println("[OS Action] ⚠️ Kích hoạt khóa màn hình LINUX...")

	// loginctl lock-session là chuẩn của systemd, khóa mọi Desktop Environment
	cmd := exec.Command("loginctl", "lock-session")
	err := cmd.Run()

	// Phương án dự phòng (Fallback) nếu Linux công ty em dùng bản cũ
	if err != nil {
		fmt.Println("[OS Action] Fallback sang xdg-screensaver...")
		exec.Command("xdg-screensaver", "lock").Run()
	}
}
