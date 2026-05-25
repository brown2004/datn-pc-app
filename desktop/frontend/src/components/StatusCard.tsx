type StatusTone = "ok" | "warn" | "danger" | "idle";

type StatusCardProps = {
  label: string;
  value: string;
  tone?: StatusTone;
};

export function StatusCard({ label, value, tone = "idle" }: StatusCardProps) {
  const toneClass = {
    ok: "border-l-[#15803d]",
    warn: "border-l-[#b45309]",
    danger: "border-l-[#b42318]",
    idle: "border-l-[#64748b]",
  }[tone];

  return (
    <div className={`flex min-h-0 min-w-0 flex-col justify-center rounded-md border border-l-[3px] border-[#d8dee8] bg-[#fbfbfc] px-4 py-2.5 ${toneClass}`}>
      <span className="mb-2 block text-xs leading-none text-[#5f6b7a]">{label}</span>
      <strong className="block truncate text-[17px] font-bold leading-tight text-[#111827]">{value}</strong>
    </div>
  );
}
