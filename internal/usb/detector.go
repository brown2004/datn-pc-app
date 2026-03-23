package usb

import (
	"fmt"
	"time"

	"github.com/karalabe/hid"

	// Cú pháp import này phụ thuộc vào tên module em gõ lúc đầu
	// Nếu file go.mod của em tên là "pc-app", thì import như sau:
	"pc-app/internal/osaction"
)

// Khai báo VID và PID của STM32
// Chú ý: Lúc test thử ở nhà, em đổi số này thành VID/PID của con chuột/bàn phím em đang xài nhé!
const (
	TargetVID = 0x0483 // Mặc định của STMicroelectronics
	TargetPID = 0x5710 // Thay đổi theo cấu hình CubeMX của em
)

// StartDetector chạy ngầm liên tục để quét và kết nối với thiết bị chống trộm
func StartDetector() {
	fmt.Println("[HID Module] Đang khởi động trình giám sát USB HID...")

	for {
		deviceInfo := findSecurityDevice()
		if deviceInfo == nil {
			fmt.Println("[CẢNH BÁO] Không tìm thấy thiết bị an toàn STM32! Chờ quét lại...")
			time.Sleep(2 * time.Second)
			continue
		}

		fmt.Printf("[HID Module] Đã phát hiện thiết bị (ID: %04X:%04X). Đang kết nối...\n", deviceInfo.VendorID, deviceInfo.ProductID)

		// Bắt đầu lắng nghe thiết bị
		err := listenToDevice(deviceInfo)

		// Bất cứ khi nào hàm listenToDevice văng lỗi (do kẻ trộm rút USB)
		if err != nil {
			fmt.Println("\n[ALARM] 🚨 THIẾT BỊ BỊ RÚT HOẶC MẤT KẾT NỐI ĐỘT NGỘT:", err)

			// HÀNH ĐỘNG 1: Khóa màn hình ngay lập tức!
			osaction.LockScreen()

			// TODO (HÀNH ĐỘNG 2): Bắn tín hiệu MQTT lên server báo mất kết nối an toàn!
		}

		time.Sleep(1 * time.Second) // Nghỉ 1 nhịp trước khi thử kết nối lại
	}
}

// findSecurityDevice lọc toàn bộ thiết bị USB HID đang cắm trên máy
func findSecurityDevice() *hid.DeviceInfo {
	devices := hid.Enumerate(TargetVID, TargetPID)
	if len(devices) == 0 {
		return nil
	}
	return &devices[0] // Lấy thiết bị đầu tiên khớp mã
}

// listenToDevice mở luồng đọc dữ liệu thời gian thực (Interrupt Transfer)
func listenToDevice(info *hid.DeviceInfo) error {
	device, err := info.Open()
	if err != nil {
		return fmt.Errorf("không thể mở thiết bị: %v", err)
	}
	defer device.Close()

	fmt.Println("[HID Module] Kết nối thành công! Đang chờ dữ liệu cảm biến...")

	// Kích thước buffer phải khớp với Report Size em cấu hình trên STM32CubeMX
	buf := make([]byte, 8)

	for {
		// Dòng này là "trái tim" của hệ thống.
		// Nó sẽ đứng im chờ dữ liệu từ STM32. Nếu rút cáp, nó văng lỗi lập tức!
		n, err := device.Read(buf)
		if err != nil {
			return err
		}

		if n > 0 {
			// Logic chống trộm: Giả sử STM32 gửi byte[0] = 0x01 khi phát hiện rung lắc
			if buf[0] == 0x01 {
				fmt.Printf("\n[ALARM] 🚨 PHÁT HIỆN RUNG LẮC (Report: %X)!!!\n", buf[:n])

				// HÀNH ĐỘNG 1: Khóa màn hình
				osaction.LockScreen()

				// TODO (HÀNH ĐỘNG 2): Bắn MQTT lên Server báo rung lắc

				// Chống spam: Đợi 3 giây rồi mới quét cảnh báo tiếp
				time.Sleep(3 * time.Second)
			}
		}
	}
}
