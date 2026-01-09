import { useState } from "react";
import { Loader2, X, UserCheck, UserX } from "lucide-react";
import type { User } from "../lib/types";

interface UserBanModalProps {
  user: User | null;
  isOpen: boolean;
  onClose: () => void;
  onConfirm: () => Promise<void>;
}

export function UserBanModal({ user, isOpen, onClose, onConfirm }: UserBanModalProps) {
  const [isProcessing, setIsProcessing] = useState(false);
  const [error, setError] = useState<string | null>(null);

  if (!isOpen || !user) return null;

  const isBanning = !user.banned;

  const handleConfirm = async () => {
    try {
      setIsProcessing(true);
      setError(null);
      await onConfirm();
      onClose();
    } catch (err) {
      console.error("Failed to toggle user ban:", err);
      setError("処理に失敗しました。もう一度お試しください。");
    } finally {
      setIsProcessing(false);
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
          <div className={`mx-auto flex h-14 w-14 items-center justify-center rounded-full ${isBanning ? 'bg-amber-100 text-amber-600' : 'bg-green-100 text-green-600'}`}>
            {isBanning ? <UserX className="h-8 w-8" /> : <UserCheck className="h-8 w-8" />}
          </div>

          <div className="mt-4 text-center">
            <h3 className="text-xl font-bold text-gray-900">
              {isBanning ? "ユーザーをBANしますか？" : "ユーザーのBANを解除しますか？"}
            </h3>
            <p className="mt-2 text-sm text-gray-500">
              {isBanning 
                ? "このユーザーはシステムにログインできなくなります。" 
                : "このユーザーのログイン制限を解除します。"}
            </p>
          </div>

          <div className="mt-6 rounded-xl bg-gray-50 p-4 border border-gray-100 space-y-2">
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">ID</span>
              <span className="font-mono text-gray-900">{user.id}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">名前</span>
              <span className="font-semibold text-gray-900">{user.name}</span>
            </div>
            <div className="flex justify-between text-sm">
              <span className="text-gray-500">メール</span>
              <span className="text-gray-900">{user.email}</span>
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
              disabled={isProcessing}
              className={`flex w-full items-center justify-center gap-2 rounded-xl px-6 py-3 text-base font-bold text-white shadow-lg transition-all active:scale-[0.98] disabled:opacity-50 ${
                isBanning 
                  ? 'bg-amber-600 shadow-amber-100 hover:bg-amber-500' 
                  : 'bg-green-600 shadow-green-100 hover:bg-green-500'
              }`}
            >
              {isProcessing ? (
                <Loader2 className="h-5 w-5 animate-spin" />
              ) : (
                isBanning ? "はい、BANします" : "はい、解除します"
              )}
            </button>
            <button
              onClick={onClose}
              disabled={isProcessing}
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
