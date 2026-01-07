"use client"

import { useState } from "react"
import { Button } from "../ui/button"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "../ui/dialog"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Plus, Loader2 } from "lucide-react"
import { createLabel } from "../../services/label-service"
import { sanitizeInput } from "../../lib/utils"
import { Badge } from "../ui/badge"

interface LabelCreateButtonProps {
  onLabelCreated?: () => void
}

export function LabelCreateButton({ onLabelCreated }: LabelCreateButtonProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [labelName, setLabelName] = useState("")
  const [selectedColor, setSelectedColor] = useState("#ef4444") // デフォルト色

  const handleOpen = () => setIsOpen(true)
  const handleClose = () => setIsOpen(false)

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

  const handleCreateLabel = async () => {
    if (labelName.trim()) {
      try {
        setIsLoading(true)
        // 入力値のサニタイズ
        const sanitizedName = sanitizeInput(labelName)

        // サービスを呼び出してラベルを作成
        await createLabel({
          name: sanitizedName,
          color: selectedColor,
        })

        // 親コンポーネントに通知
        if (onLabelCreated) {
          onLabelCreated()
        }

        // 成功したらフォームをリセットしてダイアログを閉じる
        setLabelName("")
        setSelectedColor("#ef4444")
        setIsOpen(false)
      } catch (error) {
        console.error("ラベルの作成に失敗しました:", error)
        // エラー処理（実際の実装ではエラーメッセージを表示するなど）
      } finally {
        setIsLoading(false)
      }
    }
  }

  return (
    <div>
      <Button onClick={handleOpen} className="bg-blue-600 hover:bg-blue-700">
        <Plus className="mr-2 h-4 w-4" />
        ラベル作成
      </Button>
      <Dialog open={isOpen} onOpenChange={setIsOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>新規ラベル作成</DialogTitle>
          </DialogHeader>
          <div className="space-y-4 py-4">
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
            <Button variant="outline" onClick={handleClose} disabled={isLoading}>
              キャンセル
            </Button>
            <Button onClick={handleCreateLabel} disabled={!labelName.trim() || isLoading}>
              {isLoading ? (
                <>
                  <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  作成中...
                </>
              ) : (
                "作成"
              )}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}
