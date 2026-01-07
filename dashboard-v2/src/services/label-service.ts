import type { Label } from "../lib/types";

const mockLabels: Label[] = [
  { id: "lbl_1", name: "管理者", color: "#ef4444", createdAt: "2023-01-10" },
  { id: "lbl_2", name: "一般ユーザー", color: "#3b82f6", createdAt: "2023-01-15" },
  { id: "lbl_3", name: "プレミアム", color: "#a855f7", createdAt: "2023-02-20" },
  { id: "lbl_4", name: "ベータテスター", color: "#22c55e", createdAt: "2023-03-05" },
];

export async function getLabels(): Promise<Label[]> {
  await new Promise((resolve) => setTimeout(resolve, 300));
  return [...mockLabels];
}

export async function deleteLabel(labelId: string): Promise<void> {
  const index = mockLabels.findIndex((l) => l.id === labelId);
  if (index !== -1) {
    mockLabels.splice(index, 1);
  }
}
