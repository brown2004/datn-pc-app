import { AppShell } from "../../components/AppShell";
import type { Status } from "../../types/status";
import { PairingSection } from "./PairingSection";
import { StatusStrip } from "./StatusStrip";
import { TopSection } from "./TopSection";
import type { NoticeTone } from "./types";

type DashboardViewProps = {
  status: Status;
  noticeTone: NoticeTone;
  noticeMessage: string;
  mqttStatus: string;
  pairingTimeout: string;
  protectionBusy: boolean;
  pairingBusy: boolean;
  onProtectionChange: (enabled: boolean) => void;
  onCreatePairingCode: () => void;
};

export function DashboardView({
  status,
  noticeTone,
  noticeMessage,
  mqttStatus,
  pairingTimeout,
  protectionBusy,
  pairingBusy,
  onProtectionChange,
  onCreatePairingCode,
}: DashboardViewProps) {
  return (
    <AppShell>
      <div className="grid h-full min-h-0 grid-rows-[306px_186px_minmax(0,1fr)] overflow-hidden bg-[#e8ecef] text-[#14181c]">
        <TopSection
          protectionEnabled={status.protection_enabled}
          protectionBusy={protectionBusy}
          noticeTone={noticeTone}
          noticeMessage={noticeMessage}
          onProtectionChange={onProtectionChange}
        />
        <StatusStrip status={status} mqttStatus={mqttStatus} />
        <PairingSection
          pairingStatus={status.pairing_status}
          deviceCode={status.device_code}
          timeout={pairingTimeout}
          busy={pairingBusy}
          onCreateCode={onCreatePairingCode}
        />
      </div>
    </AppShell>
  );
}
