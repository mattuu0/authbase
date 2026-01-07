// セッション関連の操作を行うサービス
// 実際の実装ではバックエンドAPIとの通信を行う

export interface Session {
  id: string
  userId: string
  ipAddress: string
  userAgent: string
  createdAt: string
  expiresAt: string // 最終アクティブから有効期限に変更
  isActive: boolean
}

// セッション一覧を取得
export async function getSessions(): Promise<Session[]> {
  // 実際の実装ではAPIからデータを取得
  return mockSessions
}

// ユーザーIDでセッションを検索
export async function getSessionsByUserId(userId: string): Promise<Session[]> {
  // 実際の実装ではAPIからデータを取得
  return mockSessions.filter((session) => session.userId === userId)
}

// セッションを検索
export async function searchSessions(query: string): Promise<Session[]> {
  // 実際の実装ではAPIからデータを取得
  return mockSessions.filter(
    (session) =>
      session.id.toLowerCase().includes(query.toLowerCase()) ||
      session.userId.toLowerCase().includes(query.toLowerCase()) ||
      session.ipAddress.toLowerCase().includes(query.toLowerCase()) ||
      session.userAgent.toLowerCase().includes(query.toLowerCase()),
  )
}

// セッションを削除
export async function deleteSession(sessionId: string): Promise<void> {
  // 実際の実装ではAPIを呼び出してセッションを削除
  console.log("Delete session:", sessionId)

  // モックデータから削除
  const index = mockSessions.findIndex((s) => s.id === sessionId)
  if (index !== -1) {
    mockSessions.splice(index, 1)
    return
  }

  throw new Error("Session not found")
}

// ユーザーのすべてのセッションを削除
export async function deleteAllSessionsByUserId(userId: string): Promise<void> {
  // 実際の実装ではAPIを呼び出してユーザーのすべてのセッションを削除
  console.log("Delete all sessions for user:", userId)

  // モックデータから削除
  const newSessions = mockSessions.filter((s) => s.userId !== userId)
  mockSessions.length = 0
  mockSessions.push(...newSessions)
}

// モックデータ
const mockSessions: Session[] = [
  {
    id: "sess_123456789",
    userId: "usr_123456789",
    ipAddress: "192.168.1.1",
    userAgent:
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    createdAt: "2023-06-15T10:30:00Z",
    expiresAt: "2023-07-15T10:30:00Z",
    isActive: true,
  },
  {
    id: "sess_987654321",
    userId: "usr_123456789",
    ipAddress: "192.168.1.2",
    userAgent:
      "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
    createdAt: "2023-06-14T08:15:00Z",
    expiresAt: "2023-07-14T08:15:00Z",
    isActive: false,
  },
  {
    id: "sess_456789123",
    userId: "usr_987654321",
    ipAddress: "192.168.1.3",
    userAgent:
      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36",
    createdAt: "2023-06-15T09:00:00Z",
    expiresAt: "2023-07-15T09:00:00Z",
    isActive: true,
  },
  {
    id: "sess_789123456",
    userId: "usr_456789123",
    ipAddress: "192.168.1.4",
    userAgent:
      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
    createdAt: "2023-06-14T14:30:00Z",
    expiresAt: "2023-07-14T14:30:00Z",
    isActive: false,
  },
  {
    id: "sess_321654987",
    userId: "usr_789123456",
    ipAddress: "192.168.1.5",
    userAgent:
      "Mozilla/5.0 (iPad; CPU OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/604.1",
    createdAt: "2023-06-15T11:00:00Z",
    expiresAt: "2023-07-15T11:00:00Z",
    isActive: true,
  },
  {
    id: "sess_654987321",
    userId: "usr_321654987",
    ipAddress: "192.168.1.6",
    userAgent:
      "Mozilla/5.0 (Linux; Android 11; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.120 Mobile Safari/537.36",
    createdAt: "2023-06-14T16:00:00Z",
    expiresAt: "2023-07-14T16:00:00Z",
    isActive: false,
  },
  {
    id: "sess_987321654",
    userId: "usr_123456789",
    ipAddress: "192.168.1.7",
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:89.0) Gecko/20100101 Firefox/89.0",
    createdAt: "2023-06-13T13:30:00Z",
    expiresAt: "2023-07-13T13:30:00Z",
    isActive: false,
  },
  {
    id: "sess_654321987",
    userId: "usr_987654321",
    ipAddress: "192.168.1.8",
    userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:89.0) Gecko/20100101 Firefox/89.0",
    createdAt: "2023-06-13T10:00:00Z",
    expiresAt: "2023-07-13T10:00:00Z",
    isActive: false,
  },
]
