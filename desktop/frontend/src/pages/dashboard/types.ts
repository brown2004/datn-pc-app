export type NoticeTone = "neutral" | "success" | "error";

export type Notice = {
  tone: NoticeTone;
  message: string;
};

export type StatusTone = "ok" | "warn" | "danger";

export const initialNotice: Notice = {
  tone: "neutral",
  message: "PC Agent da san sang. Bam Refresh de cap nhat trang thai moi nhat.",
};
