import { useState, useEffect } from "react";
import { Settings, Shield, Key, Globe } from "lucide-react";
import type { Provider } from "../lib/types";
import { getProviders, toggleProvider } from "../services/provider-service";
import { cn } from "../lib/utils";
import { ProviderEditModal } from "../components/ProviderEditModal";

export default function ProvidersPage() {
  const [providers, setProviders] = useState<Provider[]>([]);
  const [loading, setLoading] = useState(true);
  const [editingProvider, setEditingProvider] = useState<Provider | null>(null);

  useEffect(() => {
    fetchProviders();
  }, []);

  const fetchProviders = async () => {
    setLoading(true);
    try {
      const data = await getProviders();
      setProviders(data);
    } catch (error) {
      console.error("Failed to fetch providers:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggle = async (code: string) => {
    try {
      await toggleProvider(code);
      setProviders(providers.map(p => p.ProviderCode === code ? { ...p, IsEnabled: p.IsEnabled === 1 ? 0 : 1 } : p));
    } catch (error) {
      console.error("Failed to toggle provider:", error);
    }
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900">プロバイダ設定</h2>
      </div>

      <div className="grid gap-6 lg:grid-cols-2">
        {loading ? (
          <div className="col-span-full py-10 text-center text-gray-500">読み込み中...</div>
        ) : (
          providers.map((provider) => (
            <div
              key={provider.ProviderCode}
              className={cn(
                "rounded-xl border bg-white p-6 shadow-sm transition-all",
                provider.IsEnabled === 1 ? "border-blue-100" : "border-gray-200 grayscale opacity-75"
              )}
            >
              <div className="mb-6 flex items-center justify-between">
                <div className="flex items-center gap-3">
                  <div className="flex h-12 w-12 items-center justify-center rounded-xl bg-gray-50">
                    <Settings className="h-6 w-6 text-gray-600" />
                  </div>
                  <div>
                    <h3 className="text-lg font-bold capitalize text-gray-900">
                      {provider.ProviderCode === "basic" ? "Basic" : provider.ProviderCode}
                    </h3>
                    <p className="text-sm text-gray-500">
                      {provider.ProviderCode === "basic" ? "ID/パスワード認証" : "OAuth2認証プロバイダ"}
                    </p>
                  </div>
                </div>
                <button
                  onClick={() => handleToggle(provider.ProviderCode)}
                  className={cn(
                    "relative inline-flex h-6 w-11 shrink-0 cursor-pointer rounded-full border-2 border-transparent transition-colors duration-200 ease-in-out focus:outline-none",
                    provider.IsEnabled === 1 ? "bg-blue-600" : "bg-gray-200"
                  )}
                >
                  <span
                    className={cn(
                      "pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 transition duration-200 ease-in-out",
                      provider.IsEnabled === 1 ? "translate-x-5" : "translate-x-0"
                    )}
                  />
                </button>
              </div>

              {provider.ProviderCode !== "basic" && (
                <>
                  <div className="space-y-4">
                    <div className="space-y-1">
                      <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
                        <Key className="h-3 w-3" />
                        Client ID
                      </label>
                      <div className="rounded-md bg-gray-50 px-3 py-2 text-sm font-mono text-gray-700 break-all">
                        {provider.ClientID || "未設定"}
                      </div>
                    </div>
                    <div className="space-y-1">
                      <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
                        <Shield className="h-3 w-3" />
                        Client Secret
                      </label>
                      <div className="rounded-md bg-gray-50 px-3 py-2 text-sm font-mono text-gray-700">
                        {provider.ClientSecret ? "••••••••••••••••" : "未設定"}
                      </div>
                    </div>
                    <div className="space-y-1">
                      <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
                        <Globe className="h-3 w-3" />
                        Callback URL
                      </label>
                      <div className="rounded-md bg-gray-50 px-3 py-2 text-sm font-mono text-gray-700 break-all">
                        {provider.CallbackURL}
                      </div>
                    </div>
                  </div>

                  <div className="mt-6 flex justify-end">
                    <button 
                      onClick={() => setEditingProvider(provider)}
                      className="text-sm font-semibold text-blue-600 hover:text-blue-500 transition-colors"
                    >
                      設定を編集
                    </button>
                  </div>
                </>
              )}
            </div>
          ))
        )}
      </div>

      <ProviderEditModal
        provider={editingProvider}
        isOpen={!!editingProvider}
        onClose={() => setEditingProvider(null)}
        onUpdated={(updated) => {
          setProviders(providers.map(p => p.ProviderCode === updated.ProviderCode ? updated : p));
        }}
      />
    </div>
  );
}
