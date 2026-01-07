"use client"

import type React from "react"

import { useState, useRef, useEffect } from "react"
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogFooter, DialogTrigger } from "../ui/dialog"
import { Button } from "../ui/button"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Avatar, AvatarFallback, AvatarImage } from "../ui/avatar"
import { Badge } from "../ui/badge"
import { X, Camera, Plus, Loader2, AlertCircle, Eye, EyeOff, RefreshCw } from "lucide-react"
import { Command, CommandEmpty, CommandGroup, CommandInput, CommandItem, CommandList } from "../ui/command"
import { Popover, PopoverContent, PopoverTrigger } from "../ui/popover"
import { cn, sanitizeInput } from "../../lib/utils"
import { createUser } from "../../services/user-service"
import { getLabels } from "../../services/label-service"
import { Alert, AlertDescription } from "../ui/alert"
import { Card, CardContent } from "../ui/card"
import { Separator } from "../ui/separator"

interface UserCreateButtonProps {
  onUserCreated?: () => void
}

export function UserCreateButton({ onUserCreated }: UserCreateButtonProps) {
  const [open, setOpen] = useState(false)
  const [name, setName] = useState("")
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [userId, setUserId] = useState(`usr_${Math.random().toString(36).substring(2, 8)}`)
  const [avatar, setAvatar] = useState("/placeholder.svg?height=40&width=40")
  const [selectedLabels, setSelectedLabels] = useState<string[]>([])
  const [previewImage, setPreviewImage] = useState<string | null>("/placeholder.svg?height=40&width=40")
  const [commandOpen, setCommandOpen] = useState(false)
  const fileInputRef = useRef<HTMLInputElement>(null)
  const [availableLabels, setAvailableLabels] = useState<Array<{ id: string; name: string; color: string }>>([])
  const [loadingLabels, setLoadingLabels] = useState(true)
  const [creating, setCreating] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [showPassword, setShowPassword] = useState(false)
  const [showConfirmPassword, setShowConfirmPassword] = useState(false)

  // ランダムなユーザーIDを生成
  const generateRandomUserId = () => {
    setUserId(`usr_${Math.random().toString(36).substring(2, 8)}`)
  }

  // ラベル一覧を取得
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

  useEffect(() => {
    if (open) {
      fetchLabels()
      generateRandomUserId() // モーダルを開いたときに新しいIDを生成
    }
  }, [open])

  const handleOpenChange = (newOpen: boolean) => {
    setOpen(newOpen)
    if (!newOpen) {
      resetForm()
    }
  }

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

  const resetForm = () => {
    setName("")
    setEmail("")
    setPassword("")
    setConfirmPassword("")
    setAvatar("/placeholder.svg?height=40&width=40")
    setSelectedLabels([])
    setPreviewImage("/placeholder.svg?height=40&width=40")
    setError(null)
  }

  const validateForm = () => {
    if (!name.trim()) {
      setError("ユーザー名を入力してください")
      return false
    }

    if (!email.trim()) {
      setError("メールアドレスを入力してください")
      return false
    }

    // メールアドレスの簡易バリデーション
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(email)) {
      setError("有効なメールアドレスを入力してください")
      return false
    }

    if (!password) {
      setError("パスワードを入力してください")
      return false
    }

    if (password.length < 8) {
      setError("パスワードは8文字以上で入力してください")
      return false
    }

    if (password !== confirmPassword) {
      setError("パスワードと確認用パスワードが一致しません")
      return false
    }

    return true
  }

  const handleCreate = async () => {
    setError(null)

    if (!validateForm()) {
      return
    }

    try {
      setCreating(true)
      // 入力値のサニタイズ
      const sanitizedName = sanitizeInput(name)
      const sanitizedEmail = sanitizeInput(email)

      // サービスを呼び出してユーザーを作成
      await createUser({
        id: userId, // 生成したIDを使用
        name: sanitizedName,
        email: sanitizedEmail,
        password: password,
        provider: "basic",
        providerId: `basic_${Date.now()}`,
        avatar,
        labels: selectedLabels,
      })

      // 成功したらフォームをリセットしてダイアログを閉じる
      resetForm()
      setOpen(false)

      // 親コンポーネントに通知
      if (onUserCreated) {
        onUserCreated()
      }
    } catch (error) {
      console.error("ユーザーの作成に失敗しました:", error)
      setError("ユーザーの作成に失敗しました")
    } finally {
      setCreating(false)
    }
  }

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogTrigger asChild>
        <Button className="bg-blue-600 hover:bg-blue-700">
          <Plus className="mr-2 h-4 w-4" />
          ユーザー作成
        </Button>
      </DialogTrigger>
      <DialogContent className="sm:max-w-[700px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle className="text-xl">新規ユーザー作成</DialogTitle>
        </DialogHeader>

        {error && (
          <Alert variant="destructive" className="mt-2">
            <AlertCircle className="h-4 w-4" />
            <AlertDescription>{error}</AlertDescription>
          </Alert>
        )}

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* 左側カラム */}
          <div className="space-y-6">
            <Card>
              <CardContent className="pt-6">
                <div className="flex flex-col items-center gap-4">
                  <div className="relative">
                    <Avatar className="h-24 w-24">
                      <AvatarImage src={previewImage || "/placeholder.svg"} alt={name} />
                      <AvatarFallback>{name.substring(0, 2) || "U"}</AvatarFallback>
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
                  <div className="text-center text-sm text-muted-foreground">
                    クリックしてアイコン画像をアップロード
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6 space-y-4">
                <h3 className="text-sm font-medium mb-2">基本情報</h3>

                <div className="grid gap-3">
                  <div className="grid gap-2">
                    <Label htmlFor="user-id">ユーザーID</Label>
                    <div className="flex">
                      <Input id="user-id" value={userId} readOnly className="rounded-r-none font-mono text-sm" />
                      <Button
                        type="button"
                        variant="secondary"
                        className="rounded-l-none"
                        onClick={generateRandomUserId}
                        title="新しいIDを生成"
                      >
                        <RefreshCw className="h-4 w-4" />
                      </Button>
                    </div>
                    <p className="text-xs text-muted-foreground">ユーザーIDはシステム内で一意の識別子です</p>
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="user-name">
                      ユーザー名 <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="user-name"
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      placeholder="ユーザー名を入力"
                      required
                    />
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="user-email">
                      メールアドレス <span className="text-red-500">*</span>
                    </Label>
                    <Input
                      id="user-email"
                      type="email"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      placeholder="example@mail.com"
                      required
                    />
                  </div>
                </div>
              </CardContent>
            </Card>
          </div>

          {/* 右側カラム */}
          <div className="space-y-6">
            <Card>
              <CardContent className="pt-6 space-y-4">
                <h3 className="text-sm font-medium mb-2">認証情報</h3>

                <div className="grid gap-3">
                  <div className="grid gap-2">
                    <Label htmlFor="user-password">
                      パスワード <span className="text-red-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        id="user-password"
                        type={showPassword ? "text" : "password"}
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        placeholder="パスワードを入力"
                        required
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute right-0 top-0 h-full px-3"
                        onClick={() => setShowPassword(!showPassword)}
                      >
                        {showPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </Button>
                    </div>
                    <p className="text-xs text-muted-foreground">8文字以上で入力してください</p>
                  </div>

                  <div className="grid gap-2">
                    <Label htmlFor="user-confirm-password">
                      パスワード（確認） <span className="text-red-500">*</span>
                    </Label>
                    <div className="relative">
                      <Input
                        id="user-confirm-password"
                        type={showConfirmPassword ? "text" : "password"}
                        value={confirmPassword}
                        onChange={(e) => setConfirmPassword(e.target.value)}
                        placeholder="パスワードを再入力"
                        required
                      />
                      <Button
                        type="button"
                        variant="ghost"
                        size="icon"
                        className="absolute right-0 top-0 h-full px-3"
                        onClick={() => setShowConfirmPassword(!showConfirmPassword)}
                      >
                        {showConfirmPassword ? <EyeOff className="h-4 w-4" /> : <Eye className="h-4 w-4" />}
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardContent className="pt-6 space-y-4">
                <h3 className="text-sm font-medium mb-2">ラベル設定</h3>

                <div className="grid gap-3">
                  <div className="grid gap-2">
                    <Label>ラベル</Label>
                    {loadingLabels ? (
                      <div className="space-y-2">
                        <div className="flex items-center gap-2">
                          <div className="h-6 w-16 rounded-full bg-gray-200 animate-pulse"></div>
                          <div className="h-6 w-20 rounded-full bg-gray-200 animate-pulse"></div>
                        </div>
                        <div className="h-10 w-full bg-gray-200 rounded animate-pulse"></div>
                      </div>
                    ) : (
                      <>
                        <div className="flex flex-wrap gap-2 mb-2 min-h-[36px] border rounded-md p-2 bg-gray-50">
                          {selectedLabels.length === 0 ? (
                            <span className="text-sm text-muted-foreground">選択されたラベルはありません</span>
                          ) : (
                            selectedLabels.map((labelId) => {
                              const labelInfo = availableLabels.find((label) => label.id === labelId)
                              return labelInfo ? (
                                <Badge
                                  key={labelId}
                                  className="pl-2 pr-1 py-1"
                                  style={{
                                    backgroundColor: labelInfo.color || "#6b7280",
                                    color: labelInfo.color ? getTextColor(labelInfo.color) : "#ffffff",
                                  }}
                                >
                                  {labelInfo.name}
                                  <Button
                                    variant="ghost"
                                    size="icon"
                                    className="h-4 w-4 ml-1"
                                    onClick={() => handleRemoveLabel(labelId)}
                                  >
                                    <X className="h-3 w-3" />
                                  </Button>
                                </Badge>
                              ) : null
                            })
                          )}
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
                                    .filter((label) => !selectedLabels.includes(label.id))
                                    .map((label) => (
                                      <CommandItem
                                        key={label.id}
                                        onSelect={() => {
                                          setSelectedLabels([...selectedLabels, label.id])
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
                </div>
              </CardContent>
            </Card>
          </div>
        </div>

        <Separator className="my-4" />

        <DialogFooter className="flex justify-between sm:justify-end gap-2">
          <Button variant="outline" onClick={() => setOpen(false)} disabled={creating}>
            キャンセル
          </Button>
          <Button onClick={handleCreate} disabled={creating} className="bg-blue-600 hover:bg-blue-700">
            {creating ? (
              <>
                <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                作成中...
              </>
            ) : (
              "ユーザーを作成"
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}
