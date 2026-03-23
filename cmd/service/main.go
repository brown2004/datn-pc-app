package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	fmt.Println("=== BẮT ĐẦU CHẠY ANTI-THEFT CORE SERVICE ===")

	// Khởi tạo một Goroutine giả lập chạy ngầm giám sát
	go func() {
		for {
			fmt.Println("[Watchdog] Core Service vẫn đang hoạt động ổn định...")
			time.Sleep(5 * time.Second)
		}
	}()

	// Tạo channel hứng tín hiệu hệ điều hành
	quit := make(chan os.Signal, 1)

	// Lắng nghe tín hiệu tắt.
	// Dùng os.Interrupt (tương đương SIGINT - Ctrl+C)
	// Để bắt thêm SIGTERM chuẩn nền tảng, ta có thể bỏ trống để bắt mọi ngắt hợp lệ,
	// hoặc bắt đích danh os.Interrupt là đủ cho hầu hết trường hợp Graceful Shutdown cơ bản.
	signal.Notify(quit, os.Interrupt)

	// App sẽ "treo" ở dòng này, chờ đợi
	<-quit

	fmt.Println("\n[Hệ thống] Nhận được lệnh tắt Service. Đang dọn dẹp...")
	time.Sleep(1 * time.Second)
	fmt.Println("=== ĐÃ TẮT ANTI-THEFT CORE SERVICE ===")
}
