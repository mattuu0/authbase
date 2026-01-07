import { useState, useRef } from "react";
import { User as UserIcon, Mail, Lock, Camera, Plus, Loader2 } from "lucide-react";
import type { CreateUserRequest, User } from "../lib/types";
import { createUser } from "../services/user-service";
import { BaseModal } from "./ui/BaseModal";
import { Button } from "./ui/Button";
import { Input } from "./ui/Input";

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
      const newUser = await createUser({
        ...formData,
        provider: "basic",
        providerId: formData.email,
        labels: ["一般ユーザー"]
      });
      onCreated(newUser);
      setFormData({ name: "", email: "", password: "", avatar: "" });
      setPreviewImage(null);
      onClose();
    } catch (err) {
      setError("ユーザーの作成に失敗しました。");
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <BaseModal
      isOpen={isOpen}
      onClose={onClose}
      title="新規ユーザー追加"
      description="新しい管理ユーザーを作成します"
      footer={
        <div className="flex gap-3">
          <Button variant="outline" onClick={onClose} className="flex-1">
            キャンセル
          </Button>
          <Button 
            onClick={handleSubmit} 
            isLoading={isSubmitting} 
            className="flex-[2]"
          >
            {!isSubmitting && <Plus className="h-4 w-4" />}
            ユーザーを作成
          </Button>
        </div>
      }
    >
      <div className="space-y-5">
        {error && <div className="text-sm text-red-600 bg-red-50 p-3 rounded-lg">{error}</div>}
        
        <div className="flex flex-col items-center gap-3">
          <div className="relative group cursor-pointer" onClick={() => fileInputRef.current?.click()}>
            <div className="h-20 w-20 rounded-full border-2 border-dashed border-gray-200 bg-gray-50 flex items-center justify-center overflow-hidden transition-all group-hover:border-blue-400">
              {previewImage ? (
                <img src={previewImage} alt="Preview" className="h-full w-full object-cover" />
              ) : (
                <UserIcon className="h-8 w-8 text-gray-300" />
              )}
            </div>
            <div className="absolute bottom-0 right-0 rounded-full bg-blue-600 p-1.5 text-white shadow-lg">
              <Camera className="h-3.5 w-3.5" />
            </div>
            <input type="file" ref={fileInputRef} className="hidden" accept="image/*" onChange={handleFileChange} />
          </div>
          <p className="text-[10px] font-bold uppercase tracking-wider text-gray-400">アイコン画像</p>
        </div>

        <Input
          label="氏名"
          icon={<UserIcon className="h-3 w-3" />}
          placeholder="山田 太郎"
          value={formData.name}
          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
          required
        />
        
        <Input
          label="メールアドレス"
          icon={<Mail className="h-3 w-3" />}
          type="email"
          placeholder="example@mail.com"
          value={formData.email}
          onChange={(e) => setFormData({ ...formData, email: e.target.value })}
          required
        />

        <Input
          label="パスワード"
          icon={<Lock className="h-3 w-3" />}
          type="password"
          placeholder="••••••••"
          value={formData.password}
          onChange={(e) => setFormData({ ...formData, password: e.target.value })}
          required
          minLength={8}
        />
      </div>
    </BaseModal>
  );
}