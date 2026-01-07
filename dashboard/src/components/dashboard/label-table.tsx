"use client"

import { useState, useEffect } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/table"
import { Input } from "../ui/input"
import { Button } from "../ui/button"
import { Badge } from "../ui/badge"
import { Pencil, Search, SlidersHorizontal, Trash2 } from "lucide-react"
import { LabelEditDialog } from "./label-edit-dialog"
import { LabelDeleteDialog } from "./label-delete-dialog"
import { getLabels, searchLabels, deleteLabel } from "../../services/label-service"
import { Skeleton } from "../ui/skeleton"

export function LabelTable() {
  const [labels, setLabels] = useState<any[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [editingLabel, setEditingLabel] = useState<any | null>(null)
  const [deletingLabel, setDeletingLabel] = useState<any | null>(null)
  const [loading, setLoading] = useState(true)
  const [isDeleting, setIsDeleting] = useState(false)

  // ラベルデータの取得
  const fetchLabels = async () => {
    try {
      setLoading(true)
      const data = await getLabels()
      setLabels(data)
    } catch (error) {
      console.error("ラベルデータの取得に失敗しました:", error)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    fetchLabels()
  }, [])

  // 検索処理
  useEffect(() => {
    const handleSearch = async () => {
      if (!searchTerm) {
        // 検索語がなければ全ラベルを取得
        fetchLabels()
        return
      }

      try {
        setLoading(true)
        const data = await searchLabels(searchTerm)
        setLabels(data)
      } catch (error) {
        console.error("ラベル検索に失敗しました:", error)
      } finally {
        setLoading(false)
      }
    }

    // 検索処理を実行（実際のアプリではデバウンスを実装するとよい）
    const timer = setTimeout(handleSearch, 300)
    return () => clearTimeout(timer)
  }, [searchTerm])

  // ラベルの削除確認ダイアログを表示
  const handleDeleteClick = (label: any) => {
    setDeletingLabel(label)
  }

  // ラベルの削除実行
  const handleConfirmDelete = async () => {
    if (!deletingLabel) return

    try {
      setIsDeleting(true)
      // サービスを呼び出してラベルを削除
      await deleteLabel(deletingLabel.id)

      // 成功したらラベルリストを更新
      await fetchLabels()

      // 削除ダイアログを閉じる
      setDeletingLabel(null)
    } catch (error) {
      console.error("ラベルの削除に失敗しました:", error)
    } finally {
      setIsDeleting(false)
    }
  }

  // ラベルの編集
  const handleEditLabel = (label: any) => {
    setEditingLabel(label)
  }

  // ラベル編集後の更新
  const handleLabelUpdated = () => {
    fetchLabels()
  }

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

  // レスポンシブ対応
  const [isMobile, setIsMobile] = useState(false)

  useEffect(() => {
    const checkIsMobile = () => {
      setIsMobile(window.innerWidth < 768)
    }

    // 初期チェック
    checkIsMobile()

    // リサイズイベントのリスナーを追加
    window.addEventListener("resize", checkIsMobile)

    // クリーンアップ
    return () => window.removeEventListener("resize", checkIsMobile)
  }, [])

  // ローディング中のスケルトン表示
  const renderSkeletonRow = (index: number) => {
    return (
      <TableRow key={`skeleton-row-${index}`}>
        <TableCell key={`skeleton-name-${index}`}>
          <Skeleton className="h-5 w-32" />
        </TableCell>
        <TableCell key={`skeleton-color-${index}`}>
          <Skeleton className="h-6 w-24 rounded-full" />
        </TableCell>
        <TableCell key={`skeleton-date-${index}`}>
          <Skeleton className="h-5 w-24" />
        </TableCell>
        <TableCell key={`skeleton-actions-${index}`}>
          <div className="flex gap-2">
            <Skeleton className="h-8 w-8 rounded-md" />
            <Skeleton className="h-8 w-8 rounded-md" />
          </div>
        </TableCell>
      </TableRow>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            placeholder="ラベル名で検索..."
            className="pl-8"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <Button variant="outline" size="icon">
          <SlidersHorizontal className="h-4 w-4" />
        </Button>
      </div>

      <div className="rounded-md border shadow-sm overflow-x-auto">
        <Table>
          <TableHeader>
            <TableRow className="bg-gray-50">
              <TableHead>ラベル名</TableHead>
              <TableHead>ラベル色</TableHead>
              <TableHead>作成日</TableHead>
              <TableHead className="w-[100px]">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              // ローディング状態
              Array(5)
                .fill(0)
                .map((_, index) => renderSkeletonRow(index))
            ) : labels.length > 0 ? (
              labels.map((label) => (
                <TableRow key={label.id}>
                  <TableCell>{label.name}</TableCell>
                  <TableCell>
                    <Badge
                      style={{
                        backgroundColor: label.color,
                        color: getTextColor(label.color),
                      }}
                    >
                      {label.name}
                    </Badge>
                  </TableCell>
                  <TableCell>{label.createdAt}</TableCell>
                  <TableCell>
                    <div className="flex items-center space-x-2">
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleEditLabel(label)}
                        className="text-blue-500 hover:text-blue-600"
                      >
                        <Pencil className="h-4 w-4" />
                      </Button>
                      <Button
                        variant="ghost"
                        size="icon"
                        onClick={() => handleDeleteClick(label)}
                        className="text-red-500 hover:text-red-600"
                      >
                        <Trash2 className="h-4 w-4" />
                      </Button>
                    </div>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} className="h-24 text-center">
                  該当するラベルが見つかりませんでした。
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {editingLabel && (
        <LabelEditDialog
          label={editingLabel}
          open={!!editingLabel}
          onOpenChange={(open) => !open && setEditingLabel(null)}
          onLabelUpdated={handleLabelUpdated}
        />
      )}

      {deletingLabel && (
        <LabelDeleteDialog
          label={deletingLabel}
          open={!!deletingLabel}
          onOpenChange={(open) => !open && setDeletingLabel(null)}
          onConfirm={handleConfirmDelete}
          isDeleting={isDeleting}
        />
      )}
    </div>
  )
}
