import { Session } from "../lib/types";

const mockSessions: Session[] = [
  {
    id: "sess_1",
    userId: "usr_123456789",
    ipAddress: "192.168.1.1",
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124",
    createdAt: "2023-06-15T10:30:00Z",
    expiresAt: "2023-07-15T10:30:00Z",
    isActive: true,
  },
  {
    id: "sess_2",
    userId: "usr_987654321",
    ipAddress: "192.168.1.3",
    userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/91.0.4472.114",
    createdAt: "2023-06-15T09:00:00Z",
    expiresAt: "2023-07-15T09:00:00Z",
    isActive: true,
  },
];

export async function getSessions(): Promise<Session[]> {
  await new Promise((resolve) => setTimeout(resolve, 300));
  return [...mockSessions];
}

export async function deleteSession(id: string): Promise<void> {
  const index = mockSessions.findIndex((s) => s.id === id);
  if (index !== -1) {
    mockSessions.splice(index, 1);
  }
}
