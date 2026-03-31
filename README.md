# PC App - Desktop Anti-Theft Application

PC App is a desktop anti-theft application for Windows and Linux that protects your computer from theft through local motion detection and alarm, with optional cloud notification.

## Architecture

The PC App consists of two separate binaries:

### **pc-agent** (Background Service)
- Monitors the USB-connected hardware device (includes accelerometer sensor)
- Detects abnormal motion/lift patterns (theft events)
- Detects USB device disconnection
- Triggers local alarms (sound, max volume, screen lock)
- Maintains MQTT connection for heartbeat and telemetry
- Sends alerts to backend API for user notifications
- Exposes gRPC IPC server for UI communication

### **pc-ui** (Desktop User Interface)
- Fyne-based GUI for status monitoring and settings
- Shows device connection status, last event, and system health
- Allows user to adjust sensitivity and volume settings
- Provides test alarm button
- Communicates with agent via gRPC IPC

## Project Structure

```
pc-app/
├── cmd/
│   ├── pc-agent/          # Agent binary entry point
│   │   └── main.go
│   └── pc-ui/             # UI binary entry point
│       └── main.go
├── internal/              # Private application code (not importable by other modules)
│   ├── agent/             # Agent app lifecycle and logic
│   │   ├── app.go
│   │   └── run.go
│   ├── ui/                # UI app and Fyne window
│   │   ├── app.go
│   │   └── window_main.go
│   ├── config/            # Configuration management
│   │   └── config.go
│   ├── domain/            # Core domain models
│   │   ├── event.go       # TheftEvent, AlertEvent types
│   │   ├── status.go      # SystemStatus type
│   │   ├── settings.go    # UserSettings type
│   │   └── sensor.go      # AccelerometerReading type
│   ├── logx/              # Logging utilities
│   │   └── logger.go
│   ├── host/              # OS-level abstractions (volume, lock, device ID)
│   │   ├── host.go        # Interface and factory
│   │   ├── host_windows.go
│   │   └── host_linux.go
│   ├── device/            # USB device communication
│   │   ├── monitor.go     # Event generation from device readings
│   │   ├── usb.go         # USB device interface
│   │   └── parser.go      # Sensor data parsing
│   ├── alarm/             # Local alarm service
│   │   └── service.go
│   ├── mqtt/              # MQTT client and payloads
│   │   ├── client.go
│   │   ├── topics.go
│   │   ├── payload.go
│   │   └── will.go
│   ├── backend/           # Backend API client
│   │   └── client.go
│   ├── ipc/               # Inter-process communication (gRPC)
│   │   ├── protocol.go    # Service definitions
│   │   ├── server.go      # Agent-side gRPC server
│   │   └── client.go      # UI-side gRPC client
│   └── store/             # Local data persistence
│       └── store.go
├── assets/                # Static resources
│   ├── alarm.wav          # Alarm sound
│   └── icon.png           # App icon
├── deploy/                # Deployment scripts and configs
│   ├── windows/
│   └── linux/
├── go.mod                 # Go module definition
├── go.sum
└── README.md              # This file
```

## Technology Stack

- **Language**: Go 1.25.0
- **UI Framework**: Fyne v2 (cross-platform GUI)
- **IPC Protocol**: gRPC (agent ↔ UI communication)
- **Messaging**: MQTT (telemetry, alerts, heartbeat)
- **USB**: github.com/karalabe/hid (hardware device access)
- **Configuration**: TOML (embedded defaults + file overrides)

## Building

### Prerequisites
- Go 1.25.0 or later
- Linux/Windows development environment
- Basic build tools (gcc, Make, etc.)

### Build pc-agent
```bash
go build -o pc-agent ./cmd/pc-agent
```

### Build pc-ui
```bash
go build -o pc-ui ./cmd/pc-ui
```

### Build both with a script
```bash
# TODO: Add build.sh script for convenience
```

## Running

### Start the agent (background service)
```bash
./pc-agent
```

### Start the UI (in another terminal)
```bash
./pc-ui
```

## Configuration

Configuration is loaded from:
1. **Embedded defaults** in `internal/config/config.go`
2. **File override** at `~/.config/pc-app/config.toml` (Linux) or `%APPDATA%\pc-app\config.toml` (Windows)

### Example config.toml
```toml
[device]
vendor_id = 0x1234
product_id = 0x5678

[mqtt]
broker = "mqtt.example.com"
port = 1883
username = "user"
password = "password"

[backend]
url = "https://api.example.com"
api_key = "your-api-key"

[thresholds]
motion_threshold = 5.0  # m/s²
check_interval_ms = 100
alarm_volume = 100

[logging]
level = "info"
```

## Features (Planned/TODO)

- [x] Project skeleton with minimal compilable code
- [ ] USB device communication
- [ ] Motion detection algorithm
- [ ] Local alarm system (sound, volume, screen lock)
- [ ] MQTT client and message handling
- [ ] Backend API integration
- [ ] gRPC IPC protocol
- [ ] Full Fyne UI with real-time status
- [ ] Configuration persistence
- [ ] Cross-platform testing (Windows, Linux)
- [ ] Packaging (installers, systemd service)

## Development Notes

### TODO Markers
Throughout the codebase, `// TODO:` comments indicate where business logic will be implemented. These are placeholder functions with stubs.

### Testing
Currently, there are no tests. Unit tests should be added as features are implemented.

### Logging
Uses custom logger in `internal/logx/` with levels: Debug, Info, Warn, Error.

## Contributing

- Keep code idiomatic and readable
- Add TODO comments for incomplete work
- Test on both Windows and Linux before committing
- Use proper Go naming conventions and package organization

## License

TODO: Add license information

## Support

For issues, questions, or contributions, please open an issue in the repository.
