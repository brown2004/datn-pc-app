import { GetStatus, SetProtectionMode, StartPairing } from "../../wailsjs/go/main/App";
import type { Status } from "../types/status";

export async function getStatus(): Promise<Status> {
  return GetStatus();
}

export async function setProtectionMode(enabled: boolean): Promise<Status> {
  return SetProtectionMode(enabled);
}

export async function startPairing(): Promise<Status> {
  return StartPairing();
}
