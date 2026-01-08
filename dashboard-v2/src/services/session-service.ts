import type { Session } from "../lib/types";

const mockSessions: Session[] = [
  {
    id: "sess_1",
    userId: "usr_123456789",
    userName: "Admin User",
    userEmail: "admin@example.com",
    ipAddress: "192.168.1.1",
    userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/91.0.4472.124",
    createdAt: "2023-06-15T10:30:00Z",
    expiresAt: "2023-07-15T10:30:00Z",
    isActive: true,
  },
  {
    id: "sess_2",
    userId: "usr_987654321",
    userName: "Test User",
    userEmail: "test@example.com",
    ipAddress: "192.168.1.3",
    userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) Chrome/91.0.4472.114",
    createdAt: "2023-06-15T09:00:00Z",
    expiresAt: "2023-07-15T09:00:00Z",
    isActive: true,
  },
  {
    id: "sess_3",
    userId: "usr_111222333",
    userName: "Jane Smith",
    userEmail: "jane@example.com",
    ipAddress: "10.0.0.5",
    userAgent: "Mozilla/5.0 (iPhone; CPU iPhone OS 14_6 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
    createdAt: "2023-06-16T15:20:00Z",
    expiresAt: "2023-07-16T15:20:00Z",
    isActive: false,
  },
  {
    id: "sess_4",
    userId: "usr_123456789",
    userName: "Admin User",
    userEmail: "admin@example.com",
    ipAddress: "172.16.0.10",
    userAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
    createdAt: "2023-06-17T08:45:00Z",
    expiresAt: "2023-07-17T08:45:00Z",
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
