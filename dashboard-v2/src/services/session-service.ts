import type { Session } from "../lib/types";

export async function getSessions(): Promise<Session[]> {
  const response = await fetch("/api/session", {
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to fetch sessions");
  return await response.json();
}

export async function deleteSession(id: string): Promise<void> {
  const response = await fetch("/api/session", {
    method: "DELETE",
    headers: { 
      "Content-Type": "application/json",
      "sessionid": id
    },
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to delete session");
}