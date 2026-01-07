import { Provider } from "../lib/types";

const mockProviders: Provider[] = [
  {
    ProviderCode: "google",
    ClientID: "123456789012-apps.googleusercontent.com",
    ClientSecret: "GOCSPX-secret",
    CallbackURL: "http://localhost:3000/auth/callback/google",
    IsEnabled: 1,
  },
  {
    ProviderCode: "github",
    ClientID: "gh-client-id",
    ClientSecret: "gh-client-secret",
    CallbackURL: "http://localhost:3000/auth/callback/github",
    IsEnabled: 1,
  },
  {
    ProviderCode: "discord",
    ClientID: "",
    ClientSecret: "",
    CallbackURL: "http://localhost:3000/auth/callback/discord",
    IsEnabled: 0,
  },
];

export async function getProviders(): Promise<Provider[]> {
  await new Promise((resolve) => setTimeout(resolve, 300));
  return [...mockProviders];
}

export async function toggleProvider(code: string): Promise<Provider> {
  const provider = mockProviders.find((p) => p.ProviderCode === code);
  if (!provider) throw new Error("Provider not found");
  provider.IsEnabled = provider.IsEnabled === 1 ? 0 : 1;
  return { ...provider };
}
