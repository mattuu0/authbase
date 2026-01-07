// プロバイダ関連の操作を行うサービス
// 実際の実装ではバックエンドAPIとの通信を行う

export interface Provider {
  ProviderCode: string
  ClientID: string
  ClientSecret: string
  CallbackURL: string
  IsEnabled: number
}

// プロバイダ一覧を取得
export async function getProviders(): Promise<Provider[]> {
  // 実際の実装ではAPIからデータを取得
  return mockProviders
}

// プロバイダを更新
export async function updateProvider(provider: Provider): Promise<Provider> {
  // 実際の実装ではAPIを呼び出してプロバイダを更新
  console.log("Update provider:", provider)

  // モックデータを更新して返す
  const index = mockProviders.findIndex((p) => p.ProviderCode === provider.ProviderCode)
  if (index !== -1) {
    mockProviders[index] = { ...provider }
    return mockProviders[index]
  }

  throw new Error("Provider not found")
}

// プロバイダ設定を更新する関数を追加
export async function updateProviderSettings(
  providerCode: string,
  settings: {
    name: string
    clientId: string
    clientSecret: string
    redirectUri: string
    scopes: string
  },
): Promise<Provider> {
  console.log("Update provider settings:", providerCode, settings)

  // モックデータを更新して返す
  const index = mockProviders.findIndex((p) => p.ProviderCode === providerCode)
  if (index !== -1) {
    mockProviders[index] = {
      ...mockProviders[index],
      ClientID: settings.clientId,
      ClientSecret: settings.clientSecret,
      CallbackURL: settings.redirectUri,
    }
    return mockProviders[index]
  }

  throw new Error("Provider not found")
}

// プロバイダの有効/無効を切り替え
export async function toggleProvider(providerCode: string): Promise<Provider> {
  // 実際の実装ではAPIを呼び出してプロバイダの有効/無効を切り替え
  console.log("Toggle provider:", providerCode)

  // モックデータを更新して返す
  const index = mockProviders.findIndex((p) => p.ProviderCode === providerCode)
  if (index !== -1) {
    mockProviders[index].IsEnabled = mockProviders[index].IsEnabled === 1 ? 0 : 1
    return mockProviders[index]
  }

  throw new Error("Provider not found")
}

// Basic認証設定を取得
export async function getBasicSettings(): Promise<{
  enabled: boolean
  hashRounds: number
}> {
  // 実際の実装ではAPIからデータを取得
  return mockBasicSettings
}

// Basic認証設定を更新
export async function updateBasicSettings(settings: {
  enabled: boolean
  hashRounds: number
}): Promise<{
  enabled: boolean
  hashRounds: number
}> {
  // 実際の実装ではAPIを呼び出して設定を更新
  console.log("Update basic settings:", settings)

  // モックデータを更新して返す
  mockBasicSettings = { ...settings }
  return mockBasicSettings
}

// システム設定を取得
export async function getSystemSettings(): Promise<{
  secretKey: string
}> {
  // 実際の実装ではAPIからデータを取得
  return mockSystemSettings
}

// システム設定を更新
export async function updateSystemSettings(settings: {
  secretKey: string
}): Promise<{
  secretKey: string
}> {
  // 実際の実装ではAPIを呼び出して設定を更新
  console.log("Update system settings:", settings)

  // モックデータを更新して返す
  mockSystemSettings = { ...settings }
  return mockSystemSettings
}

// モックデータ
const mockProviders: Provider[] = [
  {
    ProviderCode: "google",
    ClientID: "123456789012-abcdefghijklmnopqrstuvwxyz123456.apps.googleusercontent.com",
    ClientSecret: "GOCSPX-abcdefghijklmnopqrstuvwxyz123456",
    CallbackURL: "https://example.com/api/auth/callback/google",
    IsEnabled: 1,
  },
  {
    ProviderCode: "discord",
    ClientID: "",
    ClientSecret: "",
    CallbackURL: "https://example.com/api/auth/callback/discord",
    IsEnabled: 0,
  },
  {
    ProviderCode: "github",
    ClientID: "abcdef1234567890abcd",
    ClientSecret: "abcdef1234567890abcdef1234567890abcdef12",
    CallbackURL: "https://example.com/api/auth/callback/github",
    IsEnabled: 1,
  },
  {
    ProviderCode: "microsoft",
    ClientID: "",
    ClientSecret: "",
    CallbackURL: "https://example.com/api/auth/callback/microsoft",
    IsEnabled: 0,
  },
]

// Basic認証設定のモックデータ
let mockBasicSettings = {
  enabled: true,
  hashRounds: 10,
}

// システム設定のモックデータ
let mockSystemSettings = {
  secretKey: "your-secret-key-for-jwt-and-encryption",
}
