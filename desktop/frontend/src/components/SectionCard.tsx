import type { ReactNode } from "react";

type SectionCardProps = {
  title: string;
  description?: string;
  action?: ReactNode;
  children: ReactNode;
};

export function SectionCard({ title, description, action, children }: SectionCardProps) {
  return (
    <section className="grid min-h-0 grid-rows-[auto_minmax(0,1fr)] overflow-hidden rounded-lg border border-[#9adfed] bg-white shadow-[0_3px_10px_rgba(0,175,203,0.1)]">
      <header className="flex min-h-[42px] items-center justify-between gap-3 border-b border-[#9adfed] bg-[#eaf8fc] px-3.5 py-2.5">
        <div>
          <h2 className="m-0 text-sm font-bold text-[#0096b8]">{title}</h2>
          {description && <p className="mt-1 text-xs text-[#718397]">{description}</p>}
        </div>
        {action && <div className="flex items-center">{action}</div>}
      </header>
      <div className="min-h-0 p-3">{children}</div>
    </section>
  );
}
