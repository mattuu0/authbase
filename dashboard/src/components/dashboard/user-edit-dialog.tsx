"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter } from "../ui/dialog"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"
import { Badge } from "../ui/badge"
import { X, Camera, Copy, Check, Plus, Loader2, Trash2, AlertCircle } from "lucide-react"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "../ui/command"
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover"
import { cn, sanitizeInput } from "../../lib/utils"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip"
import { updateUser, deleteUser } from "../../services/user-service"
import { getLabels } from "../../services/label-service"
import { Skeleton } from "../ui/skeleton"
import { Alert, AlertDescription } from "../ui/alert"
import { UserDeleteDialog } from "./user-delete-dialog"

// プロバイダーの表示名マッピング
const providerNames: Record<string, string> = {
  google: "Google",
  github: "GitHub",
  microsoft: "Microsoft",
  discord: "Discord",
  basic: "Basic",
}

interface UserEditDialogProps {
  user: {
    id: string
    name: string
    email: string
    provider: string
    providerId: string
    avatar: string
    labels: string[]
    createdAt: string
    banned: boolean
  }
  open: boolean
  onOpenChange: (open: boolean) => void
  onUserDeleted?: () => void
}

export function UserEditDialog({ user, open, onOpenChange, onUserDeleted }: UserEditDialogProps) {
  const [name, setName] = useState(user.name)
  const [avatar, setAvatar] = useState(user.avatar)
  const [selectedLabels, setSelectedLabels] = useState<string[]>(user.labels)
  const [previewImage, setPreviewImage] = useState<string | null>(user.avatar)
  const [commandOpen, setCommandOpen] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [copied, setCopied] = useState<Record<string, boolean>>({})
  const [availableLabels, setAvailableLabels] = useState<Array<{ id: string; name: string; color: string }>>([])
  const [loadingLabels, setLoadingLabels] = useState(true)
  const [saving, setSaving] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  // ラベル一覧を取得
  useEffect(() => {
    const fetchLabels = async () => {
      try {
        setLoadingLabels(true)
        const labels = await getLabels()
        setAvailableLabels(labels)
      } catch (error) {
        console.error("ラベルの取得に失敗しました:", error)
      } finally {
        setLoadingLabels(false)
      }
    }

    fetchLabels()
  }, [])

  const handleRemoveLabel = (label: string) => {
    setSelectedLabels(selectedLabels.filter((l) => l !== label))
  }

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0]
    if (file) {
      // ファイルをBase64に変換してプレビュー表示
      const reader = new FileReader()
      reader.onload = (event) => {
        if (event.target?.result) {
          setPreviewImage(event.target.result as string)
          // 実際の実装では、ここでアップロードサービスを呼び出し、
          // 返されたURLをavatarにセットします
          setAvatar(event.target.result as string)
        }
      }
      reader.readAsDataURL(file)
    }
  }

  const triggerFileInput = () => {
    fileInputRef.current?.click()
  }

  const copyToClipboard = (text: string, field: string) => {
    navigator.clipboard.writeText(text)
    setCopied({ ...copied, [field]: true })
    setTimeout(() => {
      setCopied({ ...copied, [field]: false })
    }, 2000)
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

  const handleSave = async () => {
    try {
      setSaving(true)
      setError(null)
      // 入力値のサニタイズ
      const sanitizedName = sanitizeInput(name)

      // サービスを呼び出してユーザー情報を更新
      await updateUser({
        id: user.id,
        name: sanitizedName,
        avatar,
        labels: selectedLabels,
      })
      onOpenChange(false)
    } catch (error) {
      console.error("ユーザー情報の更新に失敗しました:", error)
      setError("ユーザー情報の更新に失敗しました")
    } finally {
      setSaving(false)
    }
  }

  const handleDeleteUser = async () => {
    try {
      await deleteUser(user.id)
      onOpenChange(false)
      if (onUserDeleted) {
        onUserDeleted()
      }
    } catch (error) {
      console.error("ユーザーの削除に失敗しました:", error)
      setError("ユーザーの削除に失敗しました")
    }
  }

  return (
    <>
      <Dialog open={open} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-[700px]">
          <DialogHeader>
            <DialogTitle>ユーザー編集</DialogTitle>
          </DialogHeader>

          {error && (
            <Alert variant="destructive" className="mt-2">
              <AlertCircle className="h-4 w-4" />
              <AlertDescription>{error}</AlertDescription>
            </Alert>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 左側カラム */}
            <div className="space-y-4">
              <div className="flex flex-col items-center gap-4">
                <div className="relative">
                  <Avatar className="h-24 w-24">
                    <AvatarImage src={previewImage || "/placeholder.svg"} alt={name} />
                    <AvatarFallback>{name.substring(0, 2)}</AvatarFallback>
                  </Avatar>
                  <Button
                    type="button"
                    size="icon"
                    className="absolute bottom-0 right-0 rounded-full bg-primary text-primary-foreground"
                    onClick={triggerFileInput}
                  >
                    <Camera className="h-4 w-4" />
                  </Button>
                  <input
                    type="file"
                    ref={fileInputRef}
                    className="hidden"
                    accept="image/*"
                    onChange={handleFileChange}
                  />
                </div>
                <div className="text-center text-sm text-muted-foreground">クリックしてアイコン画像をアップロード</div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="user-name">ユーザー名</Label>
                <Input id="user-name" value={name} onChange={(e) => setName(e.target.value)} />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="user-id">ユーザーID</Label>
                <div className="flex">
                  <Input id="user-id" value={user.id} disabled className="rounded-r-none font-mono text-sm" />
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="secondary"
                          className="rounded-l-none"
                          onClick={() => copyToClipboard(user.id, "id")}
                        >
                          {copied["id"] ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{copied["id"] ? "コピーしました" : "コピー"}</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
              </div>

              <div className="grid gap-2">
                <Label htmlFor="user-email">メールアドレス</Label>
                <div className="flex">
                  <Input id="user-email" value={user.email} disabled className="rounded-r-none" />
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="secondary"
                          className="rounded-l-none"
                          onClick={() => copyToClipboard(user.email, "email")}
                        >
                          {copied["email"] ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{copied["email"] ? "コピーしました" : "コピー"}</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
              </div>
            </div>

            {/* 右側カラム */}
            <div className="space-y-4">
              <div className="grid gap-2">
                <Label htmlFor="provider">認証プロバイダ</Label>
                <Input id="provider" value={providerNames[user.provider]} disabled />
              </div>

              <div className="grid gap-2">
                <Label htmlFor="provider-id">プロバイダID</Label>
                <div className="flex">
                  <Input
                    id="provider-id"
                    value={user.providerId}
                    disabled
                    className="rounded-r-none font-mono text-sm"
                  />
                  <TooltipProvider>
                    <Tooltip>
                      <TooltipTrigger asChild>
                        <Button
                          variant="secondary"
                          className="rounded-l-none"
                          onClick={() => copyToClipboard(user.providerId, "providerId")}
                        >
                          {copied["providerId"] ? <Check className="h-4 w-4" /> : <Copy className="h-4 w-4" />}
                        </Button>
                      </TooltipTrigger>
                      <TooltipContent>{copied["providerId"] ? "コピーしました" : "コピー"}</TooltipContent>
                    </Tooltip>
                  </TooltipProvider>
                </div>
              </div>

              <div className="grid gap-2">
                <Label>ラベル</Label>
                {loadingLabels ? (
                  <div className="space-y-2">
                    <div className="flex items-center gap-2">
                      <Skeleton className="h-6 w-16 rounded-full" />
                      <Skeleton className="h-6 w-20 rounded-full" />
                    </div>
                    <Skeleton className="h-10 w-full" />
                  </div>
                ) : (
                  <>
                    <div className="flex flex-wrap gap-2 mb-2 min-h-[36px]">
                      {selectedLabels.map((labelName) => {
                        const labelInfo = availableLabels.find((label) => label.name === labelName)
                        return (
                          <Badge
                            key={labelName}
                            className="pl-2 pr-1 py-1"
                            style={{
                              backgroundColor: labelInfo?.color || "#6b7280",
                              color: labelInfo?.color ? getTextColor(labelInfo.color) : "#ffffff",
                            }}
                          >
                            {labelName}
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-4 w-4 ml-1"
                              onClick={() => handleRemoveLabel(labelName)}
                            >
                              <X className="h-3 w-3" />
                            </Button>
                          </Badge>
                        )
                      })}
                    </div>
                    <Popover open={commandOpen} onOpenChange={setCommandOpen}>
                      <PopoverTrigger asChild>
                        <Button variant="outline" className="w-full justify-between">
                          <span>ラベルを追加</span>
                          <Plus className="h-4 w-4 opacity-50" />
                        </Button>
                      </PopoverTrigger>
                      <PopoverContent className="p-0" align="start" side="bottom" sideOffset={5}>
                        <Command>
                          <CommandInput placeholder="ラベルを検索..." />
                          <CommandList>
                            <CommandEmpty>ラベルが見つかりません</CommandEmpty>
                            <CommandGroup>
                              {availableLabels
                                .filter((label) => !selectedLabels.includes(label.name))
                                .map((label) => (
                                  <CommandItem
                                    key={label.id}
                                    onSelect={() => {
                                      setSelectedLabels([...selectedLabels, label.name])
                                      setCommandOpen(false)
                                    }}
                                  >
                                    <Badge
                                      className={cn("mr-2")}
                                      style={{
                                        backgroundColor: label.color,
                                        color: getTextColor(label.color),
                                      }}
                                    >
                                      {label.name}
                                    </Badge>
                                    <span>{label.name}</span>
                                  </CommandItem>
                                ))}
                            </CommandGroup>
                          </CommandList>
                        </Command>
                      </PopoverContent>
                    </Popover>
                  </>
                )}
              </div>

              <div className="grid gap-2">
                <Label>作成日</Label>
                <Input value={user.createdAt} disabled />
              </div>
            </div>
          </div>
          <DialogFooter className="mt-4 flex justify-between">
            <Button
              variant="destructive"
              onClick={() => setDeleteDialogOpen(true)}
              disabled={saving}
              className="mr-auto"
            >
              <Trash2 className="mr-2 h-4 w-4" />
              ユーザーを削除
            </Button>
            <div className="flex gap-2">
              <Button variant="outline" onClick={() => onOpenChange(false)} disabled={saving}>
                キャンセル
              </Button>
              <Button onClick={handleSave} disabled={saving}>
                {saving ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    保存中...
                  </>
                ) : (
                  "保存"
                )}
              </Button>
            </div>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <UserDeleteDialog
        user={user}
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        onConfirm={handleDeleteUser}
      />
    </>
  )
}
