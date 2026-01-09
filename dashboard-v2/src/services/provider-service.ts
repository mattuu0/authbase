import type { Provider } from "../lib/types";

export async function getProviders(): Promise<Provider[]> {
  const response = await fetch("/api/providers/oauth", {
    credentials: "include"
  });
  if (!response.ok) throw new Error("Failed to fetch providers");
  return await response.json();
}

export async function toggleProvider(code: string): Promise<Provider> {
  const providers = await getProviders();
  const provider = providers.find((p) => p.ProviderCode === code);
  if (!provider) throw new Error("Provider not found");
  
  const updatedProvider = { ...provider, IsEnabled: provider.IsEnabled === 1 ? 0 : 1 };
  
  const response = await fetch("/api/providers/oauth", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify([updatedProvider]),
    credentials: "include"
  });
  
  if (!response.ok) throw new Error("Failed to update provider");
  return updatedProvider;
}

export async function updateProvider(code: string, updates: Partial<Provider>): Promise<Provider> {
  const providers = await getProviders();
  const index = providers.findIndex((p) => p.ProviderCode === code);
  if (index === -1) throw new Error("Provider not found");
  
  const updatedProvider = { ...providers[index], ...updates };
  
  const response = await fetch("/api/providers/oauth", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify([updatedProvider]),
    credentials: "include"
  });
  
  if (!response.ok) throw new Error("Failed to update provider");
  return updatedProvider;
}