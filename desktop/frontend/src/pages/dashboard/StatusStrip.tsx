import { formatStatusLabel, formatUsbStatus } from "./formatters";
import type { StatusTone } from "./types";
import type { Status } from "../../types/status";

type StatusStripProps = {
  status: Status;
  mqttStatus: string;
};

export function StatusStrip({ status, mqttStatus }: StatusStripProps) {
  const serverOnline = status.server_status === "online";
  const usbConnected = status.device_status === "connected";

  return (
    <section className="grid min-h-0 grid-cols-4 border-b border-[#1e2328] bg-[#eef2f4]">
      <StatusBlock
        title="USB chong trom"
        value={formatUsbStatus(status.device_status)}
        detail={status.device_status || "-"}
        tone={usbConnected ? "ok" : "danger"}
      />
      <StatusBlock
        title="Server"
        value={serverOnline ? "Online" : "Offline"}
        detail={status.backend_url || "-"}
        tone={serverOnline ? "ok" : "danger"}
      />
      <StatusBlock
        title="MQTT"
        value={formatStatusLabel(mqttStatus)}
        detail={mqttStatus === "not configured" ? "Chua cau hinh broker" : "MQTT broker"}
        tone={mqttStatus === "connected" || mqttStatus === "online" ? "ok" : mqttStatus === "not configured" ? "warn" : "danger"}
      />
      <StatusBlock
        title="Protection"
        value={status.protection_enabled ? "Active" : "Inactive"}
        detail={status.protection_status || "-"}
        tone={status.protection_enabled ? "ok" : "warn"}
        last
      />
    </section>
  );
}

function StatusBlock({
  title,
  value,
  detail,
  tone,
  last,
}: {
  title: string;
  value: string;
  detail: string;
  tone: StatusTone;
  last?: boolean;
}) {
  const toneClass = {
    ok: "text-[#087c68]",
    warn: "text-[#9a5b00]",
    danger: "text-[#a51d1d]",
  }[tone];

  return (
    <article className={`flex min-w-0 flex-col justify-center bg-[#eef2f4] px-8 text-center ${last ? "" : "border-r border-[#d6dee5]"}`}>
      <h3 className="m-0 truncate text-lg font-semibold text-[#14181c]">{title}</h3>
      <p className={`mt-3 truncate text-[24px] font-semibold tracking-tight ${toneClass}`}>{value}</p>
      <p className="mt-3 truncate text-xs text-[#69727a]">{detail}</p>
    </article>
  );
}
