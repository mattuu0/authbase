import { type User, type CreateUserRequest } from "../lib/types";

export async function getUsers(): Promise<User[]> {
  const response = await fetch("/auth/api/user/all", {
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
  const response = await fetch("/auth/api/user", {
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
  const response = await fetch("/auth/api/user", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(data),
    credentials: "include"
  });
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Failed to create user");
  }
  return await response.json();
}

export async function deleteUser(userId: string): Promise<void> {
  const response = await fetch("/auth/api/user", {
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

  const response = await fetch("/auth/api/user/ban", {
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