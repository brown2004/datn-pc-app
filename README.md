# PC App - Desktop Anti-Theft Application

PC App is a desktop anti-theft project for Windows and Linux. The current codebase keeps the background agent and removes the old mock UI. A proper desktop app should be rebuilt separately.

## Current Architecture

### pc-agent

- Monitors the connected hardware device.
- Detects motion and device disconnect events.
- Triggers local alarm behavior.
- Talks to the backend and MQTT services.
- Owns the security logic that a future desktop UI should control through IPC.

### Desktop UI

The previous browser/mock UI has been removed. The next UI should be built as a real desktop app with a clear boundary from the agent.

## Project Structure

```text
pc-app/
  cmd/
    pc-agent/              Agent entry point
  internal/
    agent/                 Agent lifecycle
    alarm/                 Local alarm service
    backend/               Backend API client
    config/                Configuration loading
    device/                Hardware/device monitoring
    domain/                Core domain models
    host/                  OS-level operations
    ipc/                   Future agent <-> UI communication
    logx/                  Logging helpers
    mqtt/                  MQTT client and payloads
    store/                 Local persistence
  assets/                  Static resources
  go.mod
  go.sum
```

## Build

```powershell
go build -o pc-agent.exe ./cmd/pc-agent
```

## Run

```powershell
.\pc-agent.exe
```

During development:

```powershell
go run ./cmd/pc-agent
```

## Desktop App Direction

Build the new UI as a separate desktop app and keep the agent as the background process. Recommended options:

- Wails: Go backend plus HTML/CSS/JS frontend, good for a Go desktop app.
- Tauri: Rust shell plus web frontend, smaller runtime but adds Rust tooling.
- Electron: easiest if the team is strong in frontend, heavier install size.
- Native C# WinUI/WPF: best Windows-native feel, but Windows-only.

For this project, Wails is the most practical fit because the backend is already Go.
