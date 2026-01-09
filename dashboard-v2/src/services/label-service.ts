import type { Label } from "../lib/types";

export async function getLabels(): Promise<Label[]> {
  const response = await fetch("/api/labels", {
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to fetch labels");
  return await response.json();
}

export async function deleteLabel(labelId: string): Promise<void> {
  const response = await fetch("/api/labels", {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id: labelId }),
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to delete label");
}

export async function createLabel(label: { name: string; color: string }): Promise<Label> {
  const response = await fetch("/api/labels", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(label),
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to create label");
  return await response.json();
}

export async function updateLabel(label: { id: string; name: string; color: string }): Promise<Label> {
  const response = await fetch("/api/labels", {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(label),
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to update label");
  return await response.json();
}