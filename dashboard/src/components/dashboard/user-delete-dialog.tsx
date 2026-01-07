"use client"

import { useState } from "react"
import { Button } from "../ui/button"
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog"
import { Loader2, AlertCircle } from "lucide-react"
import { Alert, AlertDescription } from "../ui/alert"

interface UserDeleteDialogProps {
  user: {
    id: string
    name: string
    email: string
  }
  open: boolean
  onOpenChange: (open: boolean) => void
  onConfirm: () => Promise<void>
}

export function UserDeleteDialog({ user, open, onOpenChange, onConfirm }: UserDeleteDialogProps) {
  const [isDeleting, setIsDeleting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleConfirm = async () => {
    try {
      setIsDeleting(true)
      setError(null)
      await onConfirm()
      onOpenChange(false)
    } catch (err) {
      console.error("ユーザーの削除に失敗しました:", err)
      setError("ユーザーの削除に失敗しました")
    } finally {
      setIsDeleting(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>ユーザーの削除</DialogTitle>
          <DialogDescription>以下のユーザーを削除してもよろしいですか？この操作は元に戻せません。</DialogDescription>
        </DialogHeader>

        {error && (
          <Alert variant="destructive">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <div className="py-4">
          <div className="space-y-2">
            <div className="flex justify-between">
              <span className="font-medium">ユーザー名:</span>
              <span>{user.name}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">メールアドレス:</span>
              <span>{user.email}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">ユーザーID:</span>
              <span className="font-mono text-sm">{user.id}</span>
            </div>
          </div>
          <p className="mt-4 text-sm text-destructive">
            このユーザーを削除すると、関連するすべてのデータも削除されます。この操作は元に戻せません。
          </p>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={isDeleting}>
            キャンセル
          </Button>
          <Button variant="destructive" onClick={handleConfirm} disabled={isDeleting}>
            {isDeleting ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                削除中...
              </>
            ) : (
              "削除"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
