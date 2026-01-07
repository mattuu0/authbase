import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { ChevronLeft, Save, Loader2, User as UserIcon, Mail, Tag, ShieldCheck } from "lucide-react";
import type { User } from "../lib/types";
import { getUserById, updateUser } from "../services/user-service";

export default function UserEditPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    labels: ""
  });

  useEffect(() => {
    if (id) {
      fetchUser(id);
    }
  }, [id]);

  const fetchUser = async (userId: string) => {
    try {
      const data = await getUserById(userId);
      if (data) {
        setUser(data);
        setFormData({
          name: data.name,
          email: data.email,
          labels: data.labels.join(", ")
        });
      }
    } catch (error) {
      console.error("Failed to fetch user:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!id) return;
    
    setSaving(true);
    try {
      await updateUser(id, {
        name: formData.name,
        email: formData.email,
        labels: formData.labels.split(",").map(s => s.trim()).filter(Boolean)
      });
      navigate("/dashboard/users");
    } catch (error) {
      console.error("Failed to update user:", error);
      alert("更新に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="flex h-64 items-center justify-center">
        <Loader2 className="h-8 w-8 animate-spin text-blue-600" />
      </div>
    );
  }

  if (!user) {
    return (
      <div className="text-center py-12">
        <h2 className="text-2xl font-bold text-gray-900">ユーザーが見つかりませんでした</h2>
        <Link to="/dashboard/users" className="mt-4 text-blue-600 hover:underline">
          一覧に戻る
        </Link>
      </div>
    );
  }

  return (
    <div className="max-w-3xl mx-auto space-y-6">
      <div className="flex items-center gap-4">
        <Link
          to="/dashboard/users"
          className="p-2 rounded-full hover:bg-gray-100 transition-colors text-gray-500"
        >
          <ChevronLeft className="h-6 w-6" />
        </Link>
        <div>
          <h2 className="text-3xl font-bold tracking-tight text-gray-900">ユーザーを編集</h2>
          <p className="text-gray-500 text-sm">ID: {user.id}</p>
        </div>
      </div>

      <div className="bg-white rounded-2xl border border-gray-200 shadow-sm overflow-hidden">
        <div className="p-8">
          <form onSubmit={handleSubmit} className="space-y-6">
            <div className="grid gap-6 sm:grid-cols-2">
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

              <div className="sm:col-span-2 space-y-2">
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

              <div className="sm:col-span-2 space-y-2">
                <label className="flex items-center gap-2 text-xs font-semibold uppercase text-gray-500">
                  <ShieldCheck className="h-3.5 w-3.5" />
                  プロバイダ情報
                </label>
                <div className="rounded-lg bg-gray-50 p-4 border border-gray-100">
                  <div className="flex justify-between items-center text-sm">
                    <span className="text-gray-500">プロバイダ</span>
                    <span className="font-medium text-gray-900 capitalize">{user.provider}</span>
                  </div>
                  <div className="flex justify-between items-center text-sm mt-2">
                    <span className="text-gray-500">プロバイダID</span>
                    <span className="font-mono text-xs text-gray-900">{user.providerId}</span>
                  </div>
                </div>
              </div>
            </div>

            <div className="flex justify-end gap-3 pt-6 border-t">
              <Link
                to="/dashboard/users"
                className="px-6 py-2.5 rounded-xl border border-gray-200 font-semibold text-gray-600 hover:bg-gray-50 transition-colors"
              >
                キャンセル
              </Link>
              <button
                type="submit"
                disabled={saving}
                className="flex items-center gap-2 px-6 py-2.5 rounded-xl bg-blue-600 font-semibold text-white shadow-lg shadow-blue-200 hover:bg-blue-500 active:scale-[0.98] transition-all disabled:opacity-50"
              >
                {saving ? (
                  <Loader2 className="h-5 w-5 animate-spin" />
                ) : (
                  <>
                    <Save className="h-5 w-5" />
                    変更を保存
                  </>
                )}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}
