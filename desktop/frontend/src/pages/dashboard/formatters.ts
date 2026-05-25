export function formatStatusLabel(value: string) {
  switch (value) {
    case "connected":
    case "online":
      return "Online";
    case "disconnected":
    case "offline":
      return "Offline";
    case "not configured":
      return "Not set";
    default:
      return value || "-";
  }
}

export function formatUsbStatus(value: string) {
  switch (value) {
    case "connected":
      return "Da cam";
    case "removed":
    case "disconnected":
      return "Chua cam";
    case "loading":
      return "Loading";
    default:
      return "Unknown";
  }
}

export function formatCountdown(value: string, now: number) {
  if (!value) {
    return "-";
  }

  const expiresAt = new Date(value).getTime();
  if (Number.isNaN(expiresAt)) {
    return "-";
  }

  const remainingSeconds = Math.ceil((expiresAt - now) / 1000);
  if (remainingSeconds <= 0) {
    return "Expired";
  }

  const hours = Math.floor(remainingSeconds / 3600);
  const minutes = Math.floor((remainingSeconds % 3600) / 60);
  const seconds = remainingSeconds % 60;
  const paddedSeconds = seconds.toString().padStart(2, "0");

  if (hours > 0) {
    return `${hours}:${minutes.toString().padStart(2, "0")}:${paddedSeconds}`;
  }

  return `${minutes}:${paddedSeconds}`;
}
