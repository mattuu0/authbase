import { useState, useRef, useEffect } from "react";
import { User as UserIcon, Mail, Lock, Camera, Plus, Loader2, Tag, X } from "lucide-react";
import type { CreateUserRequest, User, Label as LabelType } from "../lib/types";
import { createUser } from "../services/user-service";
import { getLabels } from "../services/label-service";
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
  const [selectedLabels, setSelectedLabels] = useState<string[]>([]);
  const [availableLabels, setAvailableLabels] = useState<LabelType[]>([]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [previewImage, setPreviewImage] = useState<string | null>(null);

  useEffect(() => {
    const fetchLabels = async () => {
      try {
        const labels = await getLabels();
        setAvailableLabels(labels);
      } catch (error) {
        console.error("Failed to fetch labels:", error);
      }
    };

    if (isOpen) {
      fetchLabels();
    }
  }, [isOpen]);

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

  const handleAddLabel = (labelName: string) => {
    if (!selectedLabels.includes(labelName)) {
      setSelectedLabels([...selectedLabels, labelName]);
    }
  };

  const handleRemoveLabel = (labelName: string) => {
    setSelectedLabels(selectedLabels.filter(l => l !== labelName));
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    // 基本的なバリデーション
    if (!formData.name.trim()) {
      setError("氏名を入力してください。");
      return;
    }
    if (!formData.email.trim()) {
      setError("メールアドレスを入力してください。");
      return;
    }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
      setError("有効なメールアドレスを入力してください。");
      return;
    }
    if (formData.password.length < 8) {
      setError("パスワードは8文字以上で入力してください。");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const newUser = await createUser({
        ...formData,
        provider: "basic",
        providerId: formData.email,
        labels: selectedLabels.length > 0 ? selectedLabels : ["一般ユーザー"]
      });
      onCreated(newUser);
      setFormData({ name: "", email: "", password: "", avatar: "" });
      setSelectedLabels([]);
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
      title="ユーザー追加"
      description="新しいユーザーを作成します"
      footer={
        <div className="flex gap-3">
          <Button variant="outline" onClick={onClose} className="flex-1" type="button">
            キャンセル
          </Button>
          <Button 
            form="user-create-form"
            type="submit"
            isLoading={isSubmitting} 
            className="flex-[2]"
          >
            {!isSubmitting && <Plus className="h-4 w-4" />}
            ユーザーを作成
          </Button>
        </div>
      }
    >
      <form id="user-create-form" onSubmit={handleSubmit} className="space-y-5">
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

        <div className="space-y-1.5">
          <label className="text-[11px] font-bold uppercase tracking-wider text-gray-500 flex items-center gap-1.5 px-1">
            <Tag className="h-3.5 w-3.5" />
            ラベル
          </label>
          <div className="flex flex-wrap gap-1.5 min-h-[32px] mb-2">
            {selectedLabels.map((label) => (
              <span
                key={label}
                className="inline-flex items-center gap-1 rounded-full bg-blue-50 px-2.5 py-0.5 text-xs font-medium text-blue-700 border border-blue-100"
              >
                {label}
                <button type="button" onClick={() => handleRemoveLabel(label)} className="hover:text-blue-900">
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
      </form>
    </BaseModal>
  );
}