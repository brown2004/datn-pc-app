type PairingSectionProps = {
  pairingStatus: string;
  deviceCode: string;
  timeout: string;
  busy: boolean;
  onCreateCode: () => void;
};

export function PairingSection(props: PairingSectionProps) {
  const linked = props.pairingStatus === "linked";

  return (
    <section className="grid min-h-0 grid-cols-[minmax(0,1fr)_420px] overflow-hidden bg-[#eef2f4]">
      <PairingCopy linked={linked} />
      <PairingPanel {...props} />
    </section>
  );
}

function PairingCopy({ linked }: { linked: boolean }) {
  return (
    <div className="flex min-w-0 flex-col justify-center px-12">
      <p className="m-0 text-[11px] font-semibold uppercase tracking-[0.2em] text-[#6c747c]">Device pairing</p>
      <h2 className="mt-3 truncate text-[28px] font-semibold tracking-tight text-[#14181c]">
        {linked ? "Da lien ket thiet bi" : "Lien ket thiet bi"}
      </h2>
      <p className="mt-3 max-w-[720px] text-sm leading-6 text-[#535b63]">
        {linked
          ? "Backend da xac nhan ma lien ket va credential da duoc chuyen cho PC Agent local."
          : "Tao ma lien ket, sau do nhap ma tren mobile app. Khi backend xac nhan, UI se chuyen credential cho PC Agent local."}
      </p>
    </div>
  );
}

function PairingPanel({
  pairingStatus,
  deviceCode,
  timeout,
  busy,
  onCreateCode,
}: PairingSectionProps) {
  const linked = pairingStatus === "linked";
  const syncing = pairingStatus === "syncing local agent";

  return (
    <div className="grid min-h-0 content-center gap-3 border-l border-[#d6dee5] bg-[#f8fafb] p-6">
      {linked ? (
        <div className="grid h-20 place-items-center rounded-md border border-[#9bbcac] bg-white px-4 text-center text-[24px] font-semibold text-[#087c68]">
          Da lien ket
        </div>
      ) : (
        <div className="grid h-20 place-items-center rounded-md border border-[#bac5ce] bg-white font-mono text-[31px] font-semibold tracking-[6px] text-[#14181c]">
          {deviceCode || "------"}
        </div>
      )}
      <InfoLine label="Status" value={pairingStatus || "-"} />
      <InfoLine label="Timeout" value={linked || syncing ? "-" : timeout} />
      <button
        className="h-11 rounded-md border border-[#111820] bg-[#111820] px-4 text-sm font-semibold text-white transition hover:bg-[#2a3037] disabled:cursor-not-allowed disabled:opacity-50"
        type="button"
        disabled={busy || syncing || linked}
        onClick={onCreateCode}
      >
        {busy ? "Generating..." : linked ? "Linked" : "Generate Code"}
      </button>
    </div>
  );
}

function InfoLine({ label, value }: { label: string; value: string }) {
  return (
    <div className="grid min-w-0 grid-cols-[74px_minmax(0,1fr)] items-center gap-3 text-[13px] text-[#69727a]">
      <span>{label}</span>
      <strong className="min-w-0 truncate font-mono text-xs font-semibold text-[#14181c]">{value}</strong>
    </div>
  );
}
