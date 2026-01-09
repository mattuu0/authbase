export async function login(username: string, password: string): Promise<void> {
  const response = await fetch("/auth/admin/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
    credentials: "include",
  });
  
  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new Error(errorData.error || "Invalid credentials");
  }

  // Backend uses cookies, but we set a flag for frontend consistency if needed
  localStorage.setItem("auth_token", "session_active");
}

export async function logout(): Promise<void> {
  await fetch("/auth/admin/logout", { 
    method: "POST",
    credentials: "include",
  });
  // We don't really need to remove anything from localStorage if we don't use tokens there anymore,
  // but if the app relies on it for isAuthenticated() check, we should handle it.
  localStorage.removeItem("auth_token");
}

export async function isAuthenticated(): Promise<boolean> {
  try {
    const response = await fetch("/auth/admin/info", {
      credentials: "include",
    });
    if (!response.ok) {
      localStorage.removeItem("auth_token");
      return false;
    }
    return true;
  } catch {
    localStorage.removeItem("auth_token");
    return false;
  }
}

export async function getCurrentUser() {
  const response = await fetch("/auth/admin/info", {
    credentials: "include",
  });
  if (!response.ok) return null;
  return await response.json();
}