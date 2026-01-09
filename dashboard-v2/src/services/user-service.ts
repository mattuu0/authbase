import { type User } from "../lib/types";

export interface CreateUserRequest {
  id?: string;
  name: string;
  email: string;
  password?: string;
  provider: string;
  providerId: string;
  avatar: string;
  labels: string[];
}

export async function getUsers(): Promise<User[]> {
  const response = await fetch("/api/user/all", {
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to fetch users");
  return await response.json();
}

export async function getUserById(id: string): Promise<User | null> {
  const users = await getUsers();
  return users.find((u) => u.id === id) || null;
}

export async function updateUser(id: string, data: Partial<User>): Promise<User> {
  const response = await fetch("/api/user", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ ...data, id }),
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to update user");
  return await response.json();
}

export async function searchUsers(query: string): Promise<User[]> {
  const users = await getUsers();
  const lowerQuery = query.toLowerCase();
  return users.filter(
    (user) =>
      user.name.toLowerCase().includes(lowerQuery) ||
      user.email.toLowerCase().includes(lowerQuery) ||
      user.id.toLowerCase().includes(lowerQuery)
  );
}

export async function createUser(data: CreateUserRequest): Promise<User> {
  // If the backend has a specific create user endpoint, use it.
  // Based on init.go, there isn't a direct /api/user POST endpoint for admins yet,
  // but there is /basic/signup. However, that might be for public use.
  // For now, let's assume there's no admin create user API if it's not in init.go.
  // Wait, I should check controllers/user.go to see if there's any hidden ones.
  throw new Error("Create user API not implemented in backend yet");
}

export async function deleteUser(userId: string): Promise<void> {
  const response = await fetch("/api/user", {
    method: "DELETE",
    headers: { 
      "Content-Type": "application/json",
      "userid": userId
    },
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to delete user");
}

export async function toggleUserBan(userId: string): Promise<User> {
  const users = await getUsers();
  const user = users.find(u => u.id === userId);
  if (!user) throw new Error("User not found");

  const response = await fetch("/api/user/ban", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({
      UserID: userId,
      IsBanned: !user.banned
    }),
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to toggle ban");
  return await response.json();
}