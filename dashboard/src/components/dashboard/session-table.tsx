"use client"

import { useState } from "react"
import { MoreHorizontal, Trash2 } from "lucide-react"
import { Button } from "../ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/table"
import { SessionDeleteDialog } from "./session-delete-dialog"
import { deleteSession } from "../../services/session-service"
import { useToast } from "../ui/use-toast"

export interface Session {
  id: string
  userId: string
  ipAddress: string
  userAgent: string
  createdAt: string
  expiresAt: string
  isActive: boolean
}

interface SessionTableProps {
  sessions: Session[]
  loading: boolean
  users: any[]
}

export function SessionTable({ sessions, loading, users }: SessionTableProps) {
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [selectedSession, setSelectedSession] = useState<Session | null>(null)
  const [isDeleting, setIsDeleting] = useState(false)
  const { toast } = useToast()

  const handleDeleteClick = (session: Session) => {
    setSelectedSession(session)
    setDeleteDialogOpen(true)
  }

  const handleDeleteClose = () => {
    setDeleteDialogOpen(false)
    setSelectedSession(null)
  }

  const handleDeleteConfirm = async () => {
    if (!selectedSession) return

    try {
      setIsDeleting(true)
      await deleteSession(selectedSession.id)
      toast({
        title: "セッションを削除しました",
        description: `セッション ${selectedSession.id} を削除しました`,
      })
      // ここでセッションリストを更新する処理を追加（実際の実装では親コンポーネントから渡された関数を呼び出すなど）
    } catch (error) {
      console.error("セッションの削除に失敗しました:", error)
      toast({
        title: "エラー",
        description: "セッションの削除に失敗しました",
        variant: "destructive",
      })
    } finally {
      setIsDeleting(false)
      handleDeleteClose()
    }
  }

  // ユーザー名を取得
  const getUserName = (userId: string) => {
    const user = users.find((u) => u.id === userId)
    return user ? user.name : "不明なユーザー"
  }

  // 日付をフォーマット
  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString("ja-JP", {
      year: "numeric",
      month: "2-digit",
      day: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    })
  }

  // ローディング中のスケルトン表示
  if (loading) {
    return (
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ユーザー</TableHead>
              <TableHead>IPアドレス</TableHead>
              <TableHead>ブラウザ/デバイス</TableHead>
              <TableHead>作成日時</TableHead>
              <TableHead>有効期限</TableHead>
              <TableHead>状態</TableHead>
              <TableHead className="w-[80px]">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {Array(5)
              .fill(0)
              .map((_, index) => (
                <TableRow key={index}>
                  <TableCell>
                    <div className="h-5 w-32 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-5 w-24 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-5 w-48 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-5 w-32 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-5 w-32 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-5 w-16 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                  <TableCell>
                    <div className="h-8 w-8 animate-pulse rounded-md bg-gray-200"></div>
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </div>
    )
  }

  return (
    <>
      <div className="rounded-md border">
        <Table>
          <TableHeader>
            <TableRow>
              <TableHead>ユーザー</TableHead>
              <TableHead>IPアドレス</TableHead>
              <TableHead>ブラウザ/デバイス</TableHead>
              <TableHead>作成日時</TableHead>
              <TableHead>有効期限</TableHead>
              <TableHead>状態</TableHead>
              <TableHead className="w-[80px]">操作</TableHead>
            </TableRow>
          </TableHeader>
          <TableBody>
            {sessions.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} className="h-24 text-center">
                  セッションがありません
                </TableCell>
              </TableRow>
            ) : (
              sessions.map((session) => (
                <TableRow key={session.id}>
                  <TableCell>{getUserName(session.userId)}</TableCell>
                  <TableCell>{session.ipAddress}</TableCell>
                  <TableCell className="max-w-[200px] truncate">{session.userAgent}</TableCell>
                  <TableCell>{formatDate(session.createdAt)}</TableCell>
                  <TableCell>{formatDate(session.expiresAt)}</TableCell>
                  <TableCell>
                    <span
                      className={`inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium ${
                        session.isActive ? "bg-green-100 text-green-800" : "bg-gray-100 text-gray-800"
                      }`}
                    >
                      {session.isActive ? "アクティブ" : "期限切れ"}
                    </span>
                  </TableCell>
                  <TableCell>
                    <DropdownMenu>
                      <DropdownMenuTrigger asChild>
                        <Button variant="ghost" className="h-8 w-8 p-0">
                          <span className="sr-only">メニューを開く</span>
                          <MoreHorizontal className="h-4 w-4" />
                        </Button>
                      </DropdownMenuTrigger>
                      <DropdownMenuContent align="end">
                        <DropdownMenuLabel>アクション</DropdownMenuLabel>
                        <DropdownMenuSeparator />
                        <DropdownMenuItem
                          className="text-destructive focus:text-destructive"
                          onClick={() => handleDeleteClick(session)}
                        >
                          <Trash2 className="mr-2 h-4 w-4" />
                          削除
                        </DropdownMenuItem>
                      </DropdownMenuContent>
                    </DropdownMenu>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>

      {selectedSession && (
        <SessionDeleteDialog
          open={deleteDialogOpen}
          onOpenChange={setDeleteDialogOpen}
          onConfirm={handleDeleteConfirm}
          isDeleting={isDeleting}
        />
      )}
    </>
  )
}
