// Auth service mock
export async function login(email: string, password: string): Promise<void> {
  await new Promise((resolve) => setTimeout(resolve, 800));
  if (email === "admin@example.com" && password === "password") {
    localStorage.setItem("auth_token", "mock_token");
    return;
  }
  throw new Error("Invalid credentials");
}

export async function logout(): Promise<void> {
  localStorage.removeItem("auth_token");
}

export function isAuthenticated(): boolean {
  return !!localStorage.getItem("auth_token");
}
