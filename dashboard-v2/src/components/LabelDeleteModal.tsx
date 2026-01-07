import { useState } from "react";
import { AlertTriangle, Loader2, X, Tag } from "lucide-react";
import type { Label } from "../lib/types";

interface LabelDeleteModalProps {
  label: Label | null;
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => Promise<void>;
}

export function LabelDeleteModal({ label, isOpen, onClose, onConfirm }: LabelDeleteModalProps) {
  const [isDeleting, setIsDeleting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen || !label) return null;

  const handleConfirm = async () => {
    try {
      setIsDeleting(true);
      setError(null);
      await onConfirm();
      onClose();
    } catch (err) {
      console.error("Failed to delete label:", err);
      setError("削除に失敗しました。もう一度お試しください。");
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-[60] flex items-center justify-center p-4">
      <div className="absolute inset-0 bg-gray-900/60 backdrop-blur-sm" onClick={onClose} />
      
      <div className="relative w-full max-w-md overflow-hidden rounded-2xl bg-white shadow-2xl transition-all">
        <div className="absolute right-4 top-4">
          <button onClick={onClose} className="rounded-full p-1 text-gray-400 hover:bg-gray-100 transition-colors">
            <X className="h-5 w-5" />
          </button>
        </div>

        <div className="p-6">
          <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-red-100 text-red-600">
            <AlertTriangle className="h-8 w-8" />
          </div>

          <div className="mt-4 text-center">
            <h3 className="text-xl font-bold text-gray-900">ラベルを削除しますか？</h3>
            <p className="mt-2 text-sm text-gray-500">
              この操作は取り消せません。このラベルはすべてのユーザーから削除されます。
            </p>
          </div>

          <div className="mt-6 rounded-xl bg-gray-50 p-4 border border-gray-100">
            <div className="flex items-center gap-3">
              <div 
                className="h-10 w-10 rounded-lg flex items-center justify-center"
                style={{ backgroundColor: `${label.color}20`, color: label.color }}
              >
                <Tag className="h-5 w-5" />
              </div>
              <div className="flex-1 min-w-0">
                <div className="text-sm font-bold text-gray-900 truncate">{label.name}</div>
                <div className="text-[10px] font-mono text-gray-400 uppercase">{label.color}</div>
              </div>
            </div>
          </div>

          {error && (
            <div className="mt-4 rounded-lg bg-red-50 p-3 text-sm text-red-600">
              {error}
            </div>
          )}

          <div className="mt-8 flex flex-col gap-3">
            <button
              onClick={handleConfirm}
              disabled={isDeleting}
              className="flex w-full items-center justify-center gap-2 rounded-xl bg-red-600 px-6 py-3 text-base font-bold text-white shadow-lg shadow-red-100 hover:bg-red-500 transition-all active:scale-[0.98] disabled:opacity-50"
            >
              {isDeleting ? (
                <Loader2 className="h-5 w-5 animate-spin" />
              ) : (
                "はい、削除します"
              )}
            </button>
            <button
              onClick={onClose}
              disabled={isDeleting}
              className="w-full rounded-xl border border-gray-200 bg-white px-6 py-3 text-base font-semibold text-gray-600 hover:bg-gray-50 transition-colors disabled:opacity-50"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
