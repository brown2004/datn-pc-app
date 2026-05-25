import type { ReactNode } from "react";

type AppShellProps = {
  children: ReactNode;
};

export function AppShell({ children }: AppShellProps) {
  return <div className="h-full w-full overflow-hidden bg-[#f5fbfd] text-[#17324d]">{children}</div>;
}
