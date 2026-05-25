import { useEffect, useMemo, useState } from "react";
import { getStatus, setProtectionMode, startPairing } from "../services/desktopApi";
import { emptyStatus, type Status } from "../types/status";
import { DashboardView } from "./dashboard/DashboardView";
import { formatCountdown } from "./dashboard/formatters";
import { initialNotice, type Notice } from "./dashboard/types";

const LOCAL_STATUS_POLL_INTERVAL_MS = 2000;

export function DashboardPage() {
  const [status, setStatus] = useState<Status>(emptyStatus);
  const [error, setError] = useState("");
  const [busy, setBusy] = useState("");
  const [notice, setNotice] = useState<Notice>(initialNotice);
  const [now, setNow] = useState(() => Date.now());

  async function loadStatus(options: { silent?: boolean } = {}) {
    try {
      setError("");
      const next = await getStatus();
      setStatus(next);
      if (!options.silent) {
        setNotice({ tone: "success", message: "Cap nhat trang thai he thong thanh cong." });
      }
    } catch (err) {
      reportError(err, "Cannot load status");
    }
  }

  async function updateProtection(enabled: boolean) {
    try {
      setBusy("protection");
      setError("");
      const next = await setProtectionMode(enabled);
      setStatus(next);
      setNotice({ tone: "success", message: enabled ? "Da bat che do bao ve." : "Da tat che do bao ve." });
    } catch (err) {
      reportError(err, "Cannot update protection mode");
    } finally {
      setBusy("");
    }
  }

  async function createPairingCode() {
    try {
      setBusy("pairing");
      setError("");
      const next = await startPairing();
      setStatus(next);
      setNotice({ tone: "success", message: "Da tao ma lien ket. Nhap ma nay tren mobile app de xac nhan thiet bi." });
    } catch (err) {
      reportError(err, "Cannot start pairing");
    } finally {
      setBusy("");
    }
  }

  function reportError(err: unknown, fallback: string) {
    const message = err instanceof Error ? err.message : fallback;
    setError(message);
    setNotice({ tone: "error", message });
  }

  useEffect(() => {
    loadStatus();
  }, []);

  useEffect(() => {
    const timer = window.setInterval(() => {
      loadStatus({ silent: true });
    }, LOCAL_STATUS_POLL_INTERVAL_MS);
    return () => window.clearInterval(timer);
  }, []);

  useEffect(() => {
    if (status.pairing_status !== "waiting for user" || !status.pairing_expires_at) {
      return;
    }

    setNow(Date.now());
    const timer = window.setInterval(() => {
      setNow(Date.now());
    }, 1000);
    return () => window.clearInterval(timer);
  }, [status.pairing_expires_at, status.pairing_status]);

  const pairingTimeout = useMemo(
    () => formatCountdown(status.pairing_expires_at, now),
    [now, status.pairing_expires_at],
  );
  const mqttStatus = status.mqtt_status || "not configured";
  const noticeMessage = error || status.last_error || notice.message;

  return (
    <DashboardView
      status={status}
      noticeTone={error || status.last_error ? "error" : notice.tone}
      noticeMessage={noticeMessage}
      mqttStatus={mqttStatus}
      pairingTimeout={pairingTimeout}
      protectionBusy={busy === "protection"}
      pairingBusy={busy === "pairing"}
      onProtectionChange={updateProtection}
      onCreatePairingCode={createPairingCode}
    />
  );
}
