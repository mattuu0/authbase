// ラベル関連の操作を行うサービス
// 実際の実装ではバックエンドAPIとの通信を行う

export interface Label {
  id: string
  name: string
  color: string
  createdAt: string
}

// ラベル一覧を取得
export async function getLabels(): Promise<Label[]> {
  // 実際の実装ではAPIからデータを取得
  return mockLabels
}

// ラベルを検索
export async function searchLabels(query: string): Promise<Label[]> {
  // 実際の実装ではAPIからデータを取得
  return mockLabels.filter((label) => label.name.toLowerCase().includes(query.toLowerCase()))
}

// ラベルを作成
export async function createLabel(label: Omit<Label, "id" | "createdAt">): Promise<Label> {
  // 実際の実装ではAPIを呼び出してラベルを作成
  console.log("Create label:", label)

  // 新しいラベルを作成して返す
  const newLabel: Label = {
    id: `lbl_${Math.random().toString(36).substring(2, 8)}`,
    name: label.name,
    color: label.color,
    createdAt: new Date().toISOString().split("T")[0],
  }

  mockLabels.push(newLabel)
  return newLabel
}

// ラベルを更新
export async function updateLabel(label: Partial<Label> & { id: string }): Promise<Label> {
  // 実際の実装ではAPIを呼び出してラベルを更新
  console.log("Update label:", label)

  // モックデータを更新して返す
  const index = mockLabels.findIndex((l) => l.id === label.id)
  if (index !== -1) {
    mockLabels[index] = { ...mockLabels[index], ...label }
    return mockLabels[index]
  }

  throw new Error("Label not found")
}

// ラベルを削除
export async function deleteLabel(labelId: string): Promise<void> {
  // 実際の実装ではAPIを呼び出してラベルを削除
  console.log("Delete label:", labelId)

  // モックデータから削除
  const index = mockLabels.findIndex((l) => l.id === labelId)
  if (index !== -1) {
    mockLabels.splice(index, 1)
    return
  }

  throw new Error("Label not found")
}

// モックデータ
const mockLabels: Label[] = [
  {
    id: "lbl_123456",
    name: "管理者",
    color: "#ef4444",
    createdAt: "2023-01-10",
  },
  {
    id: "lbl_234567",
    name: "一般ユーザー",
    color: "#3b82f6",
    createdAt: "2023-01-15",
  },
  {
    id: "lbl_345678",
    name: "プレミアム",
    color: "#a855f7",
    createdAt: "2023-02-20",
  },
  {
    id: "lbl_456789",
    name: "ベータテスター",
    color: "#22c55e",
    createdAt: "2023-03-05",
  },
]
