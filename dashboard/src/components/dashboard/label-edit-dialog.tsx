"use client"

import { useState, useEffect } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "../ui/dialog"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Badge } from "../ui/badge"
import { updateLabel } from "../../services/label-service"
import { sanitizeInput } from "../../lib/utils"
import { Loader2 } from "lucide-react"

interface LabelEditDialogProps {
  label: {
    id: string
    name: string
    color: string
    createdAt: string
  }
  open: boolean
  onOpenChange: (open: boolean) => void
  onLabelUpdated?: () => void
}

export function LabelEditDialog({ label, open, onOpenChange, onLabelUpdated }: LabelEditDialogProps) {
  const [labelName, setLabelName] = useState(label.name)
  const [selectedColor, setSelectedColor] = useState("#ef4444") // デフォルト色
  const [saving, setSaving] = useState(false)

  // 初期カラーを設定
  useEffect(() => {
    // 既存の色コードを設定
    setSelectedColor(label.color)
    setLabelName(label.name)
  }, [label])

  // 色に基づいてテキスト色を決定（コントラスト確保のため）
  const getTextColor = (bgColor: string) => {
    // 16進数の色コードをRGBに変換
    const r = Number.parseInt(bgColor.slice(1, 3), 16)
    const g = Number.parseInt(bgColor.slice(3, 5), 16)
    const b = Number.parseInt(bgColor.slice(5, 7), 16)

    // 明るさを計算（YIQ方式）
    const yiq = (r * 299 + g * 587 + b * 114) / 1000

    // 明るさに基づいてテキスト色を返す
    return yiq >= 128 ? "#000000" : "#ffffff"
  }

  const handleUpdateLabel = async () => {
    if (labelName.trim()) {
      try {
        setSaving(true)
        // 入力値のサニタイズ
        const sanitizedName = sanitizeInput(labelName)

        // サービスを呼び出してラベルを更新
        await updateLabel({
          id: label.id,
          name: sanitizedName,
          color: selectedColor,
        })

        // 親コンポーネントに通知
        if (onLabelUpdated) {
          onLabelUpdated()
        }

        // 成功したらダイアログを閉じる
        onOpenChange(false)
      } catch (error) {
        console.error("ラベルの更新に失敗しました:", error)
        // エラー処理（実際の実装ではエラーメッセージを表示するなど）
      } finally {
        setSaving(false)
      }
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[425px]">
        <DialogHeader>
          <DialogTitle>ラベル編集</DialogTitle>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid gap-2">
            <Label htmlFor="label-name">ラベル名</Label>
            <Input
              id="label-name"
              value={labelName}
              onChange={(e) => setLabelName(e.target.value)}
              placeholder="ラベル名を入力"
            />
          </div>
          <div className="grid gap-2">
            <Label htmlFor="label-color">ラベル色</Label>
            <div className="flex items-center gap-4">
              <Input
                id="label-color"
                type="color"
                value={selectedColor}
                onChange={(e) => setSelectedColor(e.target.value)}
                className="w-16 h-10 p-1 cursor-pointer"
              />
              <div className="flex-1">
                <p className="text-sm text-muted-foreground mb-1">色を選択してください</p>
                <div className="flex gap-2">
                  {["#ef4444", "#3b82f6", "#22c55e", "#a855f7", "#eab308", "#6b7280"].map((color) => (
                    <div
                      key={color}
                      className="w-6 h-6 rounded-full cursor-pointer border"
                      style={{ backgroundColor: color }}
                      onClick={() => setSelectedColor(color)}
                    />
                  ))}
                </div>
              </div>
            </div>
          </div>
          <div className="mt-2">
            <Label>プレビュー</Label>
            <div className="mt-1">
              <Badge
                style={{
                  backgroundColor: selectedColor,
                  color: getTextColor(selectedColor),
                }}
              >
                {labelName || "ラベル名"}
              </Badge>
            </div>
          </div>
        </div>
        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)} disabled={saving}>
            キャンセル
          </Button>
          <Button onClick={handleUpdateLabel} disabled={!labelName.trim() || saving}>
            {saving ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                更新中...
              </>
            ) : (
              "更新"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
