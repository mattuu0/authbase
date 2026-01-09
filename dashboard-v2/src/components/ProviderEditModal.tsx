import { useState, useEffect } from "react";
import { Key, Shield, Globe, Save } from "lucide-react";
import type { Provider } from "../lib/types";
import { updateProvider } from "../services/provider-service";
import { BaseModal } from "./ui/BaseModal";
import { Button } from "./ui/Button";
import { Input } from "./ui/Input";

interface ProviderEditModalProps {
  provider: Provider | null;
  isOpen: boolean;
  onClose: () => void;
  onUpdated: (updatedProvider: Provider) => void;
}

export function ProviderEditModal({ provider, isOpen, onClose, onUpdated }: ProviderEditModalProps) {
  const [formData, setFormData] = useState({
    ClientID: "",
    ClientSecret: "",
    CallbackURL: ""
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (provider) {
      setFormData({
        ClientID: provider.ClientID || "",
        ClientSecret: provider.ClientSecret || "",
        CallbackURL: provider.CallbackURL || ""
      });
    }
  }, [provider, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!provider) return;

    setIsSubmitting(true);
    setError(null);

    try {
      const updated = await updateProvider(provider.ProviderCode, formData);
      onUpdated(updated);
      onClose();
    } catch (err) {
      setError("プロバイダ設定の更新に失敗しました。");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (!provider) return null;

  return (
    <BaseModal
      isOpen={isOpen}
      onClose={onClose}
      title={`${provider.ProviderCode.toUpperCase()} 設定編集`}
      description="認証プロバイダの資格情報とコールバックURLを設定します"
      footer={
        <div className="flex gap-3">
          <Button variant="outline" onClick={onClose} className="flex-1" type="button">
            キャンセル
          </Button>
          <Button 
            form="provider-edit-form"
            type="submit"
            isLoading={isSubmitting} 
            className="flex-[2]"
          >
            {!isSubmitting && <Save className="h-4 w-4" />}
            設定を保存
          </Button>
        </div>
      }
    >
      <form id="provider-edit-form" onSubmit={handleSubmit} className="space-y-5">
        {error && <div className="text-sm text-red-600 bg-red-50 p-3 rounded-lg">{error}</div>}

        <Input
          label="Client ID"
          icon={<Key className="h-3 w-3" />}
          placeholder="Client IDを入力"
          value={formData.ClientID}
          onChange={(e) => setFormData({ ...formData, ClientID: e.target.value })}
          required
        />
        
        <Input
          label="Client Secret"
          icon={<Shield className="h-3 w-3" />}
          type="password"
          placeholder="Client Secretを入力"
          value={formData.ClientSecret}
          onChange={(e) => setFormData({ ...formData, ClientSecret: e.target.value })}
          required
        />

        <Input
          label="Callback URL"
          icon={<Globe className="h-3 w-3" />}
          placeholder="https://example.com/auth/callback"
          value={formData.CallbackURL}
          onChange={(e) => setFormData({ ...formData, CallbackURL: e.target.value })}
          required
        />
      </form>
    </BaseModal>
  );
}
