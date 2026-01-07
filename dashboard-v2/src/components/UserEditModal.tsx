import { useState, useEffect } from "react";
import { X, Save, Loader2, User as UserIcon, Mail, Tag } from "lucide-react";
import type { User } from "../lib/types";
import { updateUser } from "../services/user-service";

interface UserEditModalProps {
  user: User | null;
  isOpen: boolean;
  onClose: () => void;
  onUpdate: (updatedUser: User) => void;
}

export function UserEditModal({ user, isOpen, onClose, onUpdate }: UserEditModalProps) {
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    labels: ""
  });

  useEffect(() => {
    if (user) {
      setFormData({
        name: user.name,
        email: user.email,
        labels: user.labels.join(", ")
      });
    }
  }, [user]);

  if (!isOpen || !user) return null;

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSaving(true);
    try {
      const updated = await updateUser(user.id, {
        name: formData.name,
        email: formData.email,
        labels: formData.labels.split(",").map(s => s.trim()).filter(Boolean)
      });
      onUpdate(updated);
      onClose();
    } catch (error) {
      console.error("Failed to update user:", error);
      alert("更新に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      {/* Backdrop */}
      <div 
        className="absolute inset-0 bg-gray-900/40 backdrop-blur-sm transition-opacity" 
        onClick={onClose}
      />
      
      {/* Modal Content */}
      <div className="relative w-full max-w-lg overflow-hidden rounded-2xl bg-white shadow-2xl transition-all">
        <div className="flex items-center justify-between border-b px-6 py-4">
          <h3 className="text-xl font-bold text-gray-900">ユーザーを編集</h3>
          <button 
            onClick={onClose}
            className="rounded-full p-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-colors"
          >
            <X className="h-6 w-6" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          <div className="space-y-2">
            <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
              <UserIcon className="h-3.5 w-3.5" />
              氏名
            </label>
            <input
              type="text"
              required
              className="block w-full rounded-lg border border-gray-200 bg-gray-50 px-4 py-2.5 text-gray-900 outline-none transition-all focus:border-blue-500 focus:bg-white focus:ring-4 focus:ring-blue-100"
              value={formData.name}
              onChange={(e) => setFormData({ ...formData, name: e.target.value })}
            />
          </div>

          <div className="space-y-2">
            <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
              <Mail className="h-3.5 w-3.5" />
              メールアドレス
            </label>
            <input
              type="email"
              required
              className="block w-full rounded-lg border border-gray-200 bg-gray-50 px-4 py-2.5 text-gray-900 outline-none transition-all focus:border-blue-500 focus:bg-white focus:ring-4 focus:ring-blue-100"
              value={formData.email}
              onChange={(e) => setFormData({ ...formData, email: e.target.value })}
            />
          </div>

          <div className="space-y-2">
            <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
              <Tag className="h-3.5 w-3.5" />
              ラベル（カンマ区切り）
            </label>
            <input
              type="text"
              className="block w-full rounded-lg border border-gray-200 bg-gray-50 px-4 py-2.5 text-gray-900 outline-none transition-all focus:border-blue-500 focus:bg-white focus:ring-4 focus:ring-blue-100"
              placeholder="管理者, 開発, プレミアム"
              value={formData.labels}
              onChange={(e) => setFormData({ ...formData, labels: e.target.value })}
            />
          </div>

          <div className="mt-8 flex justify-end gap-3 pt-4 border-t">
            <button
              type="button"
              onClick={onClose}
              className="px-6 py-2 rounded-xl border border-gray-200 font-semibold text-gray-600 hover:bg-gray-50 transition-colors"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={saving}
              className="flex items-center gap-2 px-6 py-2 rounded-xl bg-blue-600 font-semibold text-white shadow-lg shadow-blue-200 hover:bg-blue-500 active:scale-[0.98] transition-all disabled:opacity-50"
            >
              {saving ? (
                <Loader2 className="h-5 w-5 animate-spin" />
              ) : (
                <>
                  <Save className="h-5 w-5" />
                  保存する
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
