"use client"

import { useState, useEffect } from "react"
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "../ui/table"
import { Input } from "../ui/input"
import { Button } from "../ui/button"
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuCheckboxItem,
  DropdownMenuTrigger,
} from "../ui/dropdown-menu"
import { Badge } from "../ui/badge"
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"
import { UserEditDialog } from "./user-edit-dialog"
import { Copy, Search, SlidersHorizontal, Ban, CheckCircle, Pencil, LogIn, Loader2, History } from "lucide-react"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../ui/select"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip"
import { getUsers, searchUsers, toggleUserBan, loginAsUser } from "../../services/user-service"
import { getLabels } from "../../services/label-service"
import { useNavigate } from "react-router-dom"
import { Skeleton } from "../ui/skeleton"

// プロバイダーの表示名マッピング
const providerNames: Record<string, string> = {
  google: "Google",
  github: "GitHub",
  microsoft: "Microsoft",
  discord: "Discord",
  basic: "Basic",
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

type Column = {
  id: string
  label: string
  visible: boolean
}

export function UserTable() {
  const navigate = useNavigate()
  const [users, setUsers] = useState<any[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  // selectedProviderとselectedLabelの初期値を"all"に設定
  const [selectedProvider, setSelectedProvider] = useState<string>("all")
  const [selectedLabel, setSelectedLabel] = useState<string>("all")
  const [editingUser, setEditingUser] = useState<any | null>(null)
  const [loading, setLoading] = useState(true)
  const [loadingLabels, setLoadingLabels] = useState(true)
  const [labelColors, setLabelColors] = useState<Record<string, string>>({})
  const [availableLabels, setAvailableLabels] = useState<string[]>([])
  const [actionInProgress, setActionInProgress] = useState<string | null>(null)
  const [tableContainerRef, setTableContainerRef] = useState<HTMLDivElement | null>(null)

  // 表示する列の設定（デフォルトでアイコンも表示）
  const [columns, setColumns] = useState<Column[]>([
    { id: "avatar", label: "アイコン", visible: true },
    { id: "id", label: "ユーザーID", visible: false },
    { id: "name", label: "ユーザー名", visible: true },
    { id: "email", label: "メールアドレス", visible: true },
    { id: "provider", label: "プロバイダ", visible: false },
    { id: "providerId", label: "プロバイダID", visible: false },
    { id: "labels", label: "ラベル", visible: true },
    { id: "createdAt", label: "作成日", visible: true },
    { id: "actions", label: "操作", visible: true },
  ])

  // ユーザーデータとラベルデータの取得
  const fetchData = async () => {
    try {
      setLoading(true)
      setLoadingLabels(true)

      // ユーザーデータを取得
      const userData = await getUsers()
      setUsers(userData)

      // ラベルデータを取得
      const labels = await getLabels()

      // ラベル名のリストを作成
      setAvailableLabels(labels.map((label) => label.name))

      // ラベル名と色のマッピングを作成
      const colorMap: Record<string, string> = {}
      labels.forEach((label) => {
        colorMap[label.name] = label.color
      })

      setLabelColors(colorMap)
    } catch (error) {
      console.error("データの取得に失敗しました:", error)
    } finally {
      setLoading(false)
      setLoadingLabels(false)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  // 検索処理
  useEffect(() => {
    const handleSearch = async () => {
      if (!searchTerm && !selectedProvider && !selectedLabel) {
        // フィルターがなければ全ユーザーを取得
        const data = await getUsers()
        setUsers(data)
        return
      }

      try {
        setLoading(true)
        const data = await searchUsers(searchTerm)

        // クライアント側でさらにフィルタリング
        const filtered = data.filter((user) => {
          const matchesProvider = selectedProvider === "all" ? true : user.provider === selectedProvider
          const matchesLabel = selectedLabel === "all" ? true : user.labels.includes(selectedLabel)
          return matchesProvider && matchesLabel
        })

        setUsers(filtered)
      } catch (error) {
        console.error("ユーザー検索に失敗しました:", error)
      } finally {
        setLoading(false)
      }
    }

    // 検索処理を実行（実際のアプリではデバウンスを実装するとよい）
    const timer = setTimeout(handleSearch, 300)
    return () => clearTimeout(timer)
  }, [searchTerm, selectedProvider, selectedLabel])

  // 列の表示/非表示を切り替える
  const toggleColumn = (columnId: string) => {
    setColumns(columns.map((column) => (column.id === columnId ? { ...column, visible: !column.visible } : column)))
  }

  // テキストをクリップボードにコピー
  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  // BANステータスの切り替え
  const handleToggleBan = async (userId: string) => {
    try {
      setActionInProgress(userId + "_ban")
      // サービスを呼び出してBANステータスを切り替え
      const updatedUser = await toggleUserBan(userId)

      // ユーザーリストを更新
      setUsers(users.map((user) => (user.id === userId ? { ...user, banned: updatedUser.banned } : user)))
    } catch (error) {
      console.error("BANステータスの切り替えに失敗しました:", error)
    } finally {
      setActionInProgress(null)
    }
  }

  // ユーザーとしてログイン
  const handleLoginAsUser = async (userId: string) => {
    try {
      setActionInProgress(userId + "_login")
      // サービスを呼び出してユーザーとしてログイン
      await loginAsUser(userId)

      // カスタムイベントを発行して他のコンポーネントに通知
      window.dispatchEvent(new Event("loginAsUserChanged"))

      // ページをリロードして変更を反映
      window.location.reload()
    } catch (error) {
      console.error("ユーザーとしてのログインに失敗しました:", error)
    } finally {
      setActionInProgress(null)
    }
  }

  // セッション一覧ページに移動
  const navigateToSessions = (userId: string) => {
    navigate(`/dashboard/sessions?userId=${userId}`)
  }

  // ユーザー編集後の更新
  const handleUserUpdated = () => {
    fetchData()
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
        {columns
          .filter((c) => c.visible)
          .map((column, colIndex) => (
            <TableCell key={`skeleton-${index}-${column.id}`}>
              {column.id === "avatar" ? (
                <Skeleton className="h-8 w-8 rounded-full" />
              ) : column.id === "labels" ? (
                <div className="flex gap-1">
                  <Skeleton className="h-6 w-16 rounded-full" />
                  <Skeleton className="h-6 w-12 rounded-full" />
                </div>
              ) : column.id === "actions" ? (
                <div className="flex gap-2">
                  <Skeleton className="h-8 w-8 rounded-md" />
                  <Skeleton className="h-8 w-8 rounded-md" />
                  <Skeleton className="h-8 w-8 rounded-md" />
                  <Skeleton className="h-8 w-8 rounded-md" />
                </div>
              ) : (
                <Skeleton className="h-5 w-full" />
              )}
            </TableCell>
          ))}
      </TableRow>
    )
  }

  return (
    <div className="space-y-4">
      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            placeholder="ユーザー名、メール、IDで検索..."
            className="pl-8"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="flex flex-wrap gap-2">
          <div className={isMobile ? "w-full" : "w-40"}>
            <Select value={selectedProvider} onValueChange={setSelectedProvider}>
              <SelectTrigger>
                <SelectValue placeholder="プロバイダ" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">すべて</SelectItem>
                <SelectItem value="google">Google</SelectItem>
                <SelectItem value="github">GitHub</SelectItem>
                <SelectItem value="microsoft">Microsoft</SelectItem>
                <SelectItem value="discord">Discord</SelectItem>
                <SelectItem value="basic">Basic</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div className={isMobile ? "w-full" : "w-40"}>
            {loadingLabels ? (
              <Skeleton className="h-10 w-full" />
            ) : (
              <Select value={selectedLabel} onValueChange={setSelectedLabel}>
                <SelectTrigger>
                  <SelectValue placeholder="ラベル" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">すべて</SelectItem>
                  {availableLabels.map((label) => (
                    <SelectItem key={label} value={label}>
                      {label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            )}
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="outline" size="icon">
                <SlidersHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" className="w-56">
              <DropdownMenuItem className="font-semibold">表示項目</DropdownMenuItem>
              {columns.map((column) => (
                <DropdownMenuCheckboxItem
                  key={column.id}
                  checked={column.visible}
                  onCheckedChange={() => toggleColumn(column.id)}
                >
                  {column.label}
                </DropdownMenuCheckboxItem>
              ))}
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>

      <div
        className="rounded-md border shadow-sm overflow-x-auto"
        ref={setTableContainerRef}
        style={{ maxWidth: "100%" }}
      >
        <Table>
          <TableHeader>
            <TableRow className="bg-gray-50">
              {columns.find((c) => c.id === "avatar")?.visible && <TableHead className="w-[50px]"></TableHead>}
              {columns.find((c) => c.id === "id")?.visible && <TableHead>ユーザーID</TableHead>}
              {columns.find((c) => c.id === "name")?.visible && <TableHead>ユーザー名</TableHead>}
              {columns.find((c) => c.id === "email")?.visible && <TableHead>メールアドレス</TableHead>}
              {columns.find((c) => c.id === "provider")?.visible && <TableHead>プロバイダ</TableHead>}
              {columns.find((c) => c.id === "providerId")?.visible && <TableHead>プロバイダID</TableHead>}
              {columns.find((c) => c.id === "labels")?.visible && <TableHead>ラベル</TableHead>}
              {columns.find((c) => c.id === "createdAt")?.visible && <TableHead>作成日</TableHead>}
              {columns.find((c) => c.id === "actions")?.visible && <TableHead className="w-[160px]">操作</TableHead>}
            </TableRow>
          </TableHeader>
          <TableBody>
            {loading ? (
              // ローディング状態
              Array(5)
                .fill(0)
                .map((_, index) => renderSkeletonRow(index))
            ) : users.length > 0 ? (
              users.map((user) => (
                <TableRow key={user.id}>
                  {columns.find((c) => c.id === "avatar")?.visible && (
                    <TableCell>
                      <Avatar className="h-8 w-8">
                        <AvatarImage src={user.avatar || "/placeholder.svg"} alt={user.name} />
                        <AvatarFallback>{user.name.substring(0, 2)}</AvatarFallback>
                      </Avatar>
                    </TableCell>
                  )}
                  {columns.find((c) => c.id === "id")?.visible && (
                    <TableCell className="font-mono text-sm">
                      <div className="flex items-center gap-1">
                        {user.id}
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-6 w-6"
                          onClick={() => copyToClipboard(user.id)}
                        >
                          <Copy className="h-3 w-3" />
                        </Button>
                      </div>
                    </TableCell>
                  )}
                  {columns.find((c) => c.id === "name")?.visible && <TableCell>{user.name}</TableCell>}
                  {columns.find((c) => c.id === "email")?.visible && <TableCell>{user.email}</TableCell>}
                  {columns.find((c) => c.id === "provider")?.visible && (
                    <TableCell>{providerNames[user.provider]}</TableCell>
                  )}
                  {columns.find((c) => c.id === "providerId")?.visible && (
                    <TableCell className="font-mono text-sm">
                      <div className="flex items-center gap-1">
                        {user.providerId}
                        <Button
                          variant="ghost"
                          size="icon"
                          className="h-6 w-6"
                          onClick={() => copyToClipboard(user.providerId)}
                        >
                          <Copy className="h-3 w-3" />
                        </Button>
                      </div>
                    </TableCell>
                  )}
                  {columns.find((c) => c.id === "labels")?.visible && (
                    <TableCell>
                      <div className="flex flex-wrap gap-1">
                        {user.labels.map((label: string) => (
                          <Badge
                            key={label}
                            style={{
                              backgroundColor: labelColors[label] || "#6b7280",
                              color: labelColors[label] ? getTextColor(labelColors[label]) : "#ffffff",
                            }}
                          >
                            {label}
                          </Badge>
                        ))}
                      </div>
                    </TableCell>
                  )}
                  {columns.find((c) => c.id === "createdAt")?.visible && <TableCell>{user.createdAt}</TableCell>}
                  {columns.find((c) => c.id === "actions")?.visible && (
                    <TableCell>
                      <div className="flex items-center space-x-2">
                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleToggleBan(user.id)}
                                className={
                                  user.banned
                                    ? "text-red-500 hover:text-red-600"
                                    : "text-green-500 hover:text-green-600"
                                }
                                disabled={actionInProgress === user.id + "_ban"}
                              >
                                {actionInProgress === user.id + "_ban" ? (
                                  <Loader2 className="h-4 w-4 animate-spin" />
                                ) : user.banned ? (
                                  <Ban className="h-4 w-4" />
                                ) : (
                                  <CheckCircle className="h-4 w-4" />
                                )}
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>{user.banned ? "BANを解除" : "BANする"}</TooltipContent>
                          </Tooltip>
                        </TooltipProvider>

                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => setEditingUser(user)}
                                className="text-blue-500 hover:text-blue-600"
                              >
                                <Pencil className="h-4 w-4" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>編集</TooltipContent>
                          </Tooltip>
                        </TooltipProvider>

                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => navigateToSessions(user.id)}
                                className="text-amber-500 hover:text-amber-600"
                              >
                                <History className="h-4 w-4" />
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>セッション一覧</TooltipContent>
                          </Tooltip>
                        </TooltipProvider>

                        <TooltipProvider>
                          <Tooltip>
                            <TooltipTrigger asChild>
                              <Button
                                variant="ghost"
                                size="icon"
                                onClick={() => handleLoginAsUser(user.id)}
                                className="text-purple-500 hover:text-purple-600"
                                disabled={actionInProgress === user.id + "_login"}
                              >
                                {actionInProgress === user.id + "_login" ? (
                                  <Loader2 className="h-4 w-4 animate-spin" />
                                ) : (
                                  <LogIn className="h-4 w-4" />
                                )}
                              </Button>
                            </TooltipTrigger>
                            <TooltipContent>このユーザーとしてログイン</TooltipContent>
                          </Tooltip>
                        </TooltipProvider>
                      </div>
                    </TableCell>
                  )}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.filter((c) => c.visible).length} className="h-24 text-center">
                  該当するユーザーが見つかりませんでした。
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>

      {editingUser && (
        <UserEditDialog
          user={editingUser}
          open={!!editingUser}
          onOpenChange={(open) => {
            if (!open) {
              setEditingUser(null)
              // ユーザー編集後にデータを更新
              handleUserUpdated()
            }
          }}
          onUserDeleted={handleUserUpdated}
        />
      )}
    </div>
  )
}
