import { useState, useRef, useEffect } from "react";
import { 
  X, 
  Camera, 
  Copy, 
  Check, 
  Plus, 
  Loader2, 
  Trash2, 
  AlertCircle,
  Mail,
  User as UserIcon,
  Shield,
  Calendar,
  Tag
} from "lucide-react";
import type { User, Label as LabelType } from "../lib/types";
import { updateUser, deleteUser } from "../services/user-service";
import { getLabels } from "../services/label-service";
import { cn } from "../lib/utils";

// プロバイダーの表示名マッピング
const providerNames: Record<string, string> = {
  google: "Google",
  github: "GitHub",
  microsoft: "Microsoft",
  discord: "Discord",
  basic: "Basic",
};

interface UserEditModalProps {
  user: User | null;
  isOpen: boolean;
  onClose: () => void;
  onUpdate: (updatedUser: User) => void;
  onDelete: (userId: string) => void;
}

export function UserEditModal({ user, isOpen, onClose, onUpdate, onDelete }: UserEditModalProps) {
  const [name, setName] = useState("");
  const [avatar, setAvatar] = useState("");
  const [selectedLabels, setSelectedLabels] = useState<string[]>([]);
  const [availableLabels, setAvailableLabels] = useState<LabelType[]>([]);
  const [loadingLabels, setLoadingLabels] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState<Record<string, boolean>>({});
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (user) {
      setName(user.name);
      setAvatar(user.avatar);
      setSelectedLabels(user.labels);
    }
    
    const fetchLabels = async () => {
      try {
        setLoadingLabels(true);
        const labels = await getLabels();
        setAvailableLabels(labels);
      } catch (error) {
        console.error("Failed to fetch labels:", error);
      } finally {
        setLoadingLabels(false);
      }
    };

    if (isOpen) {
      fetchLabels();
    }
  }, [user, isOpen]);

  if (!isOpen || !user) return null;

  const handleRemoveLabel = (label: string) => {
    setSelectedLabels(selectedLabels.filter((l) => l !== label));
  };

  const handleAddLabel = (labelName: string) => {
    if (!selectedLabels.includes(labelName)) {
      setSelectedLabels([...selectedLabels, labelName]);
    }
  };

  const copyToClipboard = (text: string, field: string) => {
    navigator.clipboard.writeText(text);
    setCopied({ ...copied, [field]: true });
    setTimeout(() => {
      setCopied({ ...copied, [field]: false });
    }, 2000);
  };

  const handleSave = async () => {
    setSaving(true);
    setError(null);
    try {
      const updated = await updateUser(user.id, {
        name,
        avatar,
        labels: selectedLabels,
      });
      onUpdate(updated);
      onClose();
    } catch (error) {
      setError("更新に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  const handleDelete = async () => {
    if (!confirm("このユーザーを完全に削除しますか？この操作は取り消せません。")) return;
    try {
      await deleteUser(user.id);
      onDelete(user.id);
      onClose();
    } catch (error) {
      setError("削除に失敗しました");
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-gray-900/60 backdrop-blur-sm" onClick={onClose} />
      
      <div className="relative w-full max-w-2xl overflow-hidden rounded-2xl bg-white shadow-2xl transition-all">
        {/* Header */}
        <div className="flex items-center justify-between border-b px-6 py-4">
          <h3 className="text-xl font-bold text-gray-900">ユーザー編集</h3>
          <button onClick={onClose} className="rounded-full p-1 text-gray-400 hover:bg-gray-100 transition-colors">
            <X className="h-6 w-6" />
          </button>
        </div>

        <div className="p-6">
          {error && (
            <div className="mb-6 flex items-center gap-2 rounded-lg bg-red-50 p-4 text-sm text-red-700">
              <AlertCircle className="h-4 w-4 shrink-0" />
              <p>{error}</p>
            </div>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
            {/* Left Column */}
            <div className="space-y-6">
              <div className="flex flex-col items-center gap-3">
                <div className="relative group">
                  <img
                    src={avatar || "https://api.dicebear.com/7.x/avataaars/svg?seed=placeholder"}
                    alt={name}
                    className="h-24 w-24 rounded-full border-2 border-gray-100 object-cover transition-opacity group-hover:opacity-75"
                  />
                  <button 
                    onClick={() => fileInputRef.current?.click()}
                    className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity"
                  >
                    <div className="bg-black/40 p-2 rounded-full text-white">
                      <Camera className="h-5 w-5" />
                    </div>
                  </button>
                  <input type="file" ref={fileInputRef} className="hidden" accept="image/*" />
                </div>
                <p className="text-xs text-gray-500">クリックして画像をアップロード</p>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  <UserIcon className="h-3 w-3" />
                  ユーザー名
                </label>
                <input
                  type="text"
                  className="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm focus:border-blue-500 focus:bg-white outline-none transition-all"
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  ID
                </label>
                <div className="flex group">
                  <input
                    disabled
                    className="flex-1 rounded-l-lg border border-gray-200 bg-gray-100 px-3 py-2 text-sm font-mono text-gray-500"
                    value={user.id}
                  />
                  <button 
                    onClick={() => copyToClipboard(user.id, "id")}
                    className="rounded-r-lg border border-l-0 border-gray-200 bg-gray-50 px-3 hover:bg-gray-100 transition-colors"
                  >
                    {copied["id"] ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4 text-gray-400" />}
                  </button>
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  <Mail className="h-3 w-3" />
                  メールアドレス
                </label>
                <div className="flex">
                  <input
                    disabled
                    className="flex-1 rounded-l-lg border border-gray-200 bg-gray-100 px-3 py-2 text-sm text-gray-500"
                    value={user.email}
                  />
                  <button 
                    onClick={() => copyToClipboard(user.email, "email")}
                    className="rounded-r-lg border border-l-0 border-gray-200 bg-gray-50 px-3 hover:bg-gray-100 transition-colors"
                  >
                    {copied["email"] ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4 text-gray-400" />}
                  </button>
                </div>
              </div>
            </div>

            {/* Right Column */}
            <div className="space-y-6">
              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  <Shield className="h-3 w-3" />
                  認証プロバイダ
                </label>
                <input
                  disabled
                  className="w-full rounded-lg border border-gray-200 bg-gray-100 px-3 py-2 text-sm text-gray-500 capitalize"
                  value={providerNames[user.provider] || user.provider}
                />
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  プロバイダID
                </label>
                <div className="flex">
                  <input
                    disabled
                    className="flex-1 rounded-l-lg border border-gray-200 bg-gray-100 px-3 py-2 text-sm font-mono text-gray-500"
                    value={user.providerId}
                  />
                  <button 
                    onClick={() => copyToClipboard(user.providerId, "providerId")}
                    className="rounded-r-lg border border-l-0 border-gray-200 bg-gray-50 px-3 hover:bg-gray-100 transition-colors"
                  >
                    {copied["providerId"] ? <Check className="h-4 w-4 text-green-600" /> : <Copy className="h-4 w-4 text-gray-400" />}
                  </button>
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  <Tag className="h-3.5 w-3.5" />
                  ラベル
                </label>
                <div className="flex flex-wrap gap-1.5 mb-2 min-h-[32px]">
                  {selectedLabels.map((label) => (
                    <span
                      key={label}
                      className="inline-flex items-center gap-1 rounded-full bg-blue-50 px-2.5 py-0.5 text-xs font-medium text-blue-700 border border-blue-100"
                    >
                      {label}
                      <button onClick={() => handleRemoveLabel(label)} className="hover:text-blue-900">
                        <X className="h-3 w-3" />
                      </button>
                    </span>
                  ))}
                </div>
                <div className="relative">
                  <select 
                    className="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm focus:border-blue-500 outline-none appearance-none"
                    onChange={(e) => handleAddLabel(e.target.value)}
                    value=""
                  >
                    <option value="" disabled>ラベルを選択して追加</option>
                    {availableLabels
                      .filter(l => !selectedLabels.includes(l.name))
                      .map(l => (
                        <option key={l.id} value={l.name}>{l.name}</option>
                      ))
                    }
                  </select>
                  <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-400">
                    <Plus className="h-4 w-4" />
                  </div>
                </div>
              </div>

              <div className="space-y-1.5">
                <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-1.5">
                  <Calendar className="h-3 w-3" />
                  作成日
                </label>
                <input
                  disabled
                  className="w-full rounded-lg border border-gray-200 bg-gray-100 px-3 py-2 text-sm text-gray-500"
                  value={user.createdAt}
                />
              </div>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="flex items-center justify-between border-t bg-gray-50 px-6 py-4">
          <button
            onClick={handleDelete}
            className="flex items-center gap-2 text-sm font-semibold text-red-600 hover:text-red-700 transition-colors"
          >
            <Trash2 className="h-4 w-4" />
            ユーザーを削除
          </button>
          <div className="flex gap-3">
            <button
              onClick={onClose}
              className="px-4 py-2 text-sm font-semibold text-gray-600 hover:text-gray-800 transition-colors"
            >
              キャンセル
            </button>
            <button
              onClick={handleSave}
              disabled={saving}
              className="flex items-center gap-2 rounded-lg bg-blue-600 px-6 py-2 text-sm font-bold text-white shadow-lg shadow-blue-200 hover:bg-blue-500 transition-all disabled:opacity-50"
            >
              {saving ? <Loader2 className="h-4 w-4 animate-spin" /> : "保存する"}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}