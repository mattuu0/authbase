import { useState, useRef } from "react";
import { 
  X, 
  User as UserIcon, 
  Mail, 
  Lock, 
  Camera, 
  Plus, 
  Loader2,
  AlertCircle
} from "lucide-react";
import type { CreateUserRequest, User } from "../lib/types";
import { createUser } from "../services/user-service";

interface UserCreateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated: (newUser: User) => void;
}

export function UserCreateModal({ isOpen, onClose, onCreated }: UserCreateModalProps) {
  const [formData, setFormData] = useState({
    name: "",
    email: "",
    password: "",
    avatar: ""
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewImage, setPreviewImage] = useState<string | null>(null);

  if (!isOpen) return null;

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (event) => {
        if (event.target?.result) {
          const base64 = event.target.result as string;
          setPreviewImage(base64);
          setFormData({ ...formData, avatar: base64 });
        }
      };
      reader.readAsDataURL(file);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSubmitting(true);
    setError(null);

    try {
      const payload: CreateUserRequest = {
        name: formData.name,
        email: formData.email,
        password: formData.password,
        avatar: formData.avatar,
        provider: "basic",
        providerId: formData.email,
        labels: ["一般ユーザー"]
      };

      const newUser = await createUser(payload);
      onCreated(newUser);
      // Reset and close
      setFormData({ name: "", email: "", password: "", avatar: "" });
      setPreviewImage(null);
      onClose();
    } catch (err) {
      setError("ユーザーの作成に失敗しました。入力内容を確認してください。");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-gray-900/60 backdrop-blur-sm" onClick={onClose} />
      
      <div className="relative w-full max-w-md overflow-hidden rounded-2xl bg-white shadow-2xl transition-all">
        {/* Header */}
        <div className="flex items-center justify-between border-b px-6 py-4 bg-gray-50/50">
          <div>
            <h3 className="text-lg font-bold text-gray-900">新規ユーザー追加</h3>
            <p className="text-xs text-gray-500">新しい管理ユーザーを作成します</p>
          </div>
          <button onClick={onClose} className="rounded-full p-1 text-gray-400 hover:bg-gray-100 transition-colors">
            <X className="h-5 w-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6 space-y-5">
          {error && (
            <div className="flex items-center gap-2 rounded-lg bg-red-50 p-3 text-sm text-red-700">
              <AlertCircle className="h-4 w-4" />
              <p>{error}</p>
            </div>
          )}

          {/* Avatar Upload */}
          <div className="flex flex-col items-center gap-3">
            <div className="relative group">
              <div className="h-20 w-20 rounded-full border-2 border-dashed border-gray-200 bg-gray-50 flex items-center justify-center overflow-hidden transition-all group-hover:border-blue-400">
                {previewImage ? (
                  <img src={previewImage} alt="Preview" className="h-full w-full object-cover" />
                ) : (
                  <UserIcon className="h-8 w-8 text-gray-300" />
                )}
              </div>
              <button 
                type="button"
                onClick={() => fileInputRef.current?.click()}
                className="absolute bottom-0 right-0 rounded-full bg-blue-600 p-1.5 text-white shadow-lg transition-transform hover:scale-110 active:scale-95"
              >
                <Camera className="h-3.5 w-3.5" />
              </button>
              <input 
                type="file" 
                ref={fileInputRef} 
                className="hidden" 
                accept="image/*" 
                onChange={handleFileChange} 
              />
            </div>
            <p className="text-[10px] font-bold uppercase tracking-wider text-gray-400">アイコン画像</p>
          </div>

          <div className="space-y-4">
            {/* Name Input */}
            <div className="space-y-1.5">
              <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-2">
                <UserIcon className="h-3 w-3" />
                氏名
              </label>
              <input
                type="text"
                required
                className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-blue-500 focus:bg-white focus:ring-4 focus:ring-blue-50/50 outline-none transition-all"
                placeholder="山田 太郎"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
              />
            </div>

            {/* Email Input */}
            <div className="space-y-1.5">
              <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-2">
                <Mail className="h-3 w-3" />
                メールアドレス
              </label>
              <input
                type="email"
                required
                className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-blue-500 focus:bg-white focus:ring-4 focus:ring-blue-50/50 outline-none transition-all"
                placeholder="example@mail.com"
                value={formData.email}
                onChange={(e) => setFormData({ ...formData, email: e.target.value })}
              />
            </div>

            {/* Password Input */}
            <div className="space-y-1.5">
              <label className="text-xs font-bold uppercase text-gray-500 flex items-center gap-2">
                <Lock className="h-3 w-3" />
                パスワード
              </label>
              <input
                type="password"
                required
                minLength={8}
                className="w-full rounded-xl border border-gray-200 bg-gray-50 px-4 py-2.5 text-sm focus:border-blue-500 focus:bg-white focus:ring-4 focus:ring-blue-50/50 outline-none transition-all"
                placeholder="••••••••"
                value={formData.password}
                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
              />
            </div>
          </div>

          <div className="mt-8 flex gap-3 pt-2">
            <button
              type="button"
              onClick={onClose}
              className="flex-1 rounded-xl border border-gray-200 bg-white px-4 py-2.5 text-sm font-semibold text-gray-600 hover:bg-gray-50 transition-colors"
            >
              キャンセル
            </button>
            <button
              type="submit"
              disabled={isSubmitting}
              className="flex-[2] flex items-center justify-center gap-2 rounded-xl bg-blue-600 px-4 py-2.5 text-sm font-bold text-white shadow-lg shadow-blue-100 hover:bg-blue-500 transition-all active:scale-[0.98] disabled:opacity-50"
            >
              {isSubmitting ? (
                <Loader2 className="h-4 w-4 animate-spin" />
              ) : (
                <>
                  <Plus className="h-4 w-4" />
                  ユーザーを作成
                </>
              )}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
