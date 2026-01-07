import { useState, useRef } from "react";
import { Tag, Plus, Pipette } from "lucide-react";
import { createLabel } from "../services/label-service";
import type { Label } from "../lib/types";
import { BaseModal } from "./ui/BaseModal";
import { Button } from "./ui/Button";
import { Input } from "./ui/Input";

interface LabelCreateModalProps {
  isOpen: boolean;
  onClose: () => void;
  onCreated: (newLabel: Label) => void;
  existingLabels: Label[];
}

const PRESET_COLORS = [
  "#3b82f6", // Blue
  "#ef4444", // Red
  "#22c55e", // Green
  "#f59e0b", // Amber
  "#a855f7", // Purple
  "#ec4899", // Pink
  "#06b6d4", // Cyan
];

export function LabelCreateModal({ isOpen, onClose, onCreated, existingLabels }: LabelCreateModalProps) {
  const [name, setName] = useState("");
  const [color, setColor] = useState(PRESET_COLORS[0]);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const colorInputRef = useRef<HTMLInputElement>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const trimmedName = name.trim();
    
    if (!trimmedName) {
      setError("ラベル名を入力してください。");
      return;
    }

    // 重複チェック
    if (existingLabels.some(l => l.name.toLowerCase() === trimmedName.toLowerCase())) {
      setError("このラベル名は既に登録されています。");
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const newLabel = await createLabel({ name: trimmedName, color });
      onCreated(newLabel);
      setName("");
      setColor(PRESET_COLORS[0]);
      onClose();
    } catch (err) {
      setError("ラベルの作成に失敗しました。");
    } finally {
      setIsSubmitting(false);
    }
  };

  const isCustomColor = !PRESET_COLORS.includes(color);

  return (
    <BaseModal
      isOpen={isOpen}
      onClose={onClose}
      title="ラベル作成"
      description="新しいラベルを作成してユーザーを分類します"
      footer={
        <div className="flex gap-3">
          <Button variant="outline" onClick={onClose} className="flex-1" type="button">
            キャンセル
          </Button>
          <Button 
            form="label-create-form"
            type="submit"
            isLoading={isSubmitting} 
            className="flex-[2]"
          >
            {!isSubmitting && <Plus className="h-4 w-4" />}
            ラベルを作成
          </Button>
        </div>
      }
    >
      <form id="label-create-form" onSubmit={handleSubmit} className="space-y-6">
        {error && <div className="text-sm text-red-600 bg-red-50 p-3 rounded-lg">{error}</div>}

        <Input
          label="ラベル名"
          icon={<Tag className="h-3 w-3" />}
          placeholder="例: プレミアム会員"
          value={name}
          onChange={(e) => {
            setName(e.target.value);
            if (error) setError(null);
          }}
          required
        />

        <div className="space-y-2">
          <label className="text-xs font-bold uppercase text-gray-500 px-1">
            カラー
          </label>
          <div className="grid grid-cols-4 gap-2">
            {PRESET_COLORS.map((presetColor) => (
              <button
                key={presetColor}
                type="button"
                className={`h-10 rounded-lg border-2 transition-all ${
                  color === presetColor ? "border-blue-600 ring-2 ring-blue-100" : "border-transparent"
                }`}
                style={{ backgroundColor: presetColor }}
                onClick={() => setColor(presetColor)}
              />
            ))}
            
            {/* Custom Color Picker */}
            <div className="relative">
              <button
                type="button"
                className={`h-10 w-full rounded-lg border-2 flex items-center justify-center transition-all ${
                  isCustomColor ? "border-blue-600 ring-2 ring-blue-100" : "border-gray-200 bg-white hover:bg-gray-50"
                }`}
                style={isCustomColor ? { backgroundColor: color } : {}}
                onClick={() => colorInputRef.current?.click()}
              >
                <Pipette className={`h-4 w-4 ${isCustomColor ? "text-white drop-shadow-sm" : "text-gray-400"}`} />
              </button>
              <input
                ref={colorInputRef}
                type="color"
                className="absolute inset-0 opacity-0 cursor-pointer pointer-events-none"
                value={color}
                onChange={(e) => setColor(e.target.value)}
              />
            </div>
          </div>

          <div className="mt-4 flex items-center gap-3 rounded-xl border border-gray-100 bg-gray-50 p-3">
             <div 
               className="h-10 w-10 rounded-lg flex items-center justify-center transition-colors"
               style={{ backgroundColor: `${color}20`, color: color }}
             >
               <Tag className="h-5 w-5" />
             </div>
             <div className="flex-1">
               <div className="text-sm font-medium text-gray-900">プレビュー</div>
               <div className="text-xs text-gray-500 truncate" style={{ color: color }}>
                 {name || "ラベル名"}
               </div>
             </div>
             <div className="text-[10px] font-mono font-bold text-gray-400 uppercase bg-white px-2 py-1 rounded border border-gray-100">
               {color}
             </div>
          </div>
        </div>
      </form>
    </BaseModal>
  );
}
