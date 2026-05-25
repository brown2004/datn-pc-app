import * as Switch from "@radix-ui/react-switch";
import type { NoticeTone } from "./types";

type TopSectionProps = {
  protectionEnabled: boolean;
  protectionBusy: boolean;
  noticeTone: NoticeTone;
  noticeMessage: string;
  onProtectionChange: (enabled: boolean) => void;
};

export function TopSection({
  protectionEnabled,
  protectionBusy,
  noticeTone,
  noticeMessage,
  onProtectionChange,
}: TopSectionProps) {
  return (
    <section className="grid min-h-0 grid-cols-[340px_minmax(0,1fr)] border-b border-[#1e2328]">
      <div className="grid min-h-0 grid-rows-[72px_minmax(0,1fr)] border-r border-[#1e2328]">
        <BrandPanel />
        <ProtectionPanel enabled={protectionEnabled} busy={protectionBusy} onChange={onProtectionChange} />
      </div>

      <NoticePanel tone={noticeTone} message={noticeMessage} />
    </section>
  );
}

function BrandPanel() {
  return (
    <div className="flex items-center justify-center border-b border-[#1e2328] bg-[#f7f8f9] px-6">
      <div className="min-w-0 text-center">
        <p className="m-0 truncate text-[13px] font-semibold uppercase tracking-[0.22em] text-[#8f1d1d]">DATN</p>
        <h1 className="mt-1 truncate text-lg font-semibold tracking-tight text-[#14181c]">Laptop Guard</h1>
      </div>
    </div>
  );
}

function ProtectionPanel({
  enabled,
  busy,
  onChange,
}: {
  enabled: boolean;
  busy: boolean;
  onChange: (enabled: boolean) => void;
}) {
  return (
    <div className="flex min-h-0 flex-col items-center justify-center bg-[#eef2f4] px-8 text-center">
      <p className="m-0 text-[11px] font-semibold uppercase tracking-[0.2em] text-[#6c747c]">Protection</p>
      <h2 className="mt-3 text-[30px] font-semibold tracking-tight text-[#14181c]">{enabled ? "On" : "Off"}</h2>
      <Switch.Root
        className="relative mt-5 h-9 w-[68px] rounded-full bg-[#aeb7bf] outline-none transition data-[state=checked]:bg-[#0d8d72] disabled:cursor-not-allowed disabled:opacity-50"
        checked={enabled}
        disabled={busy}
        onCheckedChange={onChange}
      >
        <Switch.Thumb className="block h-7 w-7 translate-x-1 rounded-full bg-white shadow-[0_2px_8px_rgba(0,0,0,0.22)] transition-transform data-[state=checked]:translate-x-10" />
      </Switch.Root>
      <p className="mt-4 max-w-[220px] text-xs leading-5 text-[#606870]">Bat/tat che do bao ve tren agent local.</p>
    </div>
  );
}

function NoticePanel({ tone, message }: { tone: NoticeTone; message: string }) {
  const toneClass = {
    neutral: "text-[#8f1d1d]",
    success: "text-[#176b4d]",
    error: "text-[#a51d1d]",
  }[tone];

  return (
    <div className="min-h-0 overflow-auto bg-[#eef2f4] px-12 py-9">
      <div className="mx-auto flex min-h-full max-w-[820px] items-center justify-center">
        <p className={`m-0 whitespace-pre-wrap text-center text-[15px] font-medium leading-7 ${toneClass}`}>{message}</p>
      </div>
    </div>
  );
}
