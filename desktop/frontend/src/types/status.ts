export type Status = {
  agent_status: string;
  device_status: string;
  protection_status: string;
  backend_url: string;
  server_status: string;
  protection_enabled: boolean;
  pairing_status: string;
  device_code: string;
  pairing_expires_at: string;
  mqtt_status: string;
  last_error: string;
};

export const emptyStatus: Status = {
  agent_status: "loading",
  device_status: "loading",
  protection_status: "loading",
  backend_url: "-",
  server_status: "unknown",
  protection_enabled: false,
  pairing_status: "not started",
  device_code: "",
  pairing_expires_at: "",
  mqtt_status: "not configured",
  last_error: "",
};
