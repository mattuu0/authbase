import { type User } from "../lib/types";

export interface CreateUserRequest {
  id?: string;
  name: string;
  email: string;
  password: string;
  provider: string;
  providerId: string;
  avatar: string;
  labels: string[];
}

const mockUsers: User[] = [
  {
    id: "usr_123456789",
    name: "山田太郎",
    email: "yamada@example.com",
    provider: "google",
    providerId: "109876543210",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Yamada",
    labels: ["管理者"],
    createdAt: "2023-01-15",
    banned: false,
  },
  {
    id: "usr_987654321",
    name: "佐藤花子",
    email: "sato@example.com",
    provider: "github",
    providerId: "sato123",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Sato",
    labels: ["一般ユーザー"],
    createdAt: "2023-02-20",
    banned: false,
  },
  {
    id: "usr_456789123",
    name: "鈴木一郎",
    email: "suzuki@example.com",
    provider: "basic",
    providerId: "suzuki_ichiro",
    avatar: "https://api.dicebear.com/7.x/avataaars/svg?seed=Suzuki",
    labels: ["プレミアム"],
    createdAt: "2023-03-10",
    banned: true,
  },
];

export async function getUsers(): Promise<User[]> {
  // Simulate API delay
  await new Promise((resolve) => setTimeout(resolve, 500));
  return [...mockUsers];
}

export async function searchUsers(query: string): Promise<User[]> {
  const lowerQuery = query.toLowerCase();
  return mockUsers.filter(
    (user) =>
      user.name.toLowerCase().includes(lowerQuery) ||
      user.email.toLowerCase().includes(lowerQuery) ||
      user.id.toLowerCase().includes(lowerQuery)
  );
}

export async function deleteUser(userId: string): Promise<void> {
  const index = mockUsers.findIndex((u) => u.id === userId);
  if (index !== -1) {
    mockUsers.splice(index, 1);
  }
}

export async function toggleUserBan(userId: string): Promise<User> {
  const user = mockUsers.find((u) => u.id === userId);
  if (!user) throw new Error("User not found");
  user.banned = !user.banned;
  return { ...user };
}
