"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { Link, useNavigate } from "react-router-dom"
import { Button } from "../components/ui/button"
import { Input } from "../components/ui/input"
import { Label } from "../components/ui/label"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../components/ui/card"
import { Alert, AlertDescription } from "../components/ui/alert"
import { AlertCircle, Loader2 } from "lucide-react"
import { signup, checkAdminExists } from "../services/auth-service"

export default function SignupPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [confirmPassword, setConfirmPassword] = useState("")
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  const [isCheckingAdmin, setIsCheckingAdmin] = useState(true)

  // 管理者が既に存在するかチェック
  useEffect(() => {
    const checkAdmin = async () => {
      try {
        setIsCheckingAdmin(true)
        const adminExists = await checkAdminExists()

        // 管理者が既に存在する場合はログインページにリダイレクト
        if (adminExists) {
          navigate("/login?adminExists=true")
        }
      } catch (err) {
        console.error("管理者チェックに失敗しました:", err)
        setError("システムエラーが発生しました。しばらく経ってからお試しください。")
      } finally {
        setIsCheckingAdmin(false)
      }
    }

    checkAdmin()
  }, [navigate])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError(null)

    // パスワード確認
    if (password !== confirmPassword) {
      setError("パスワードが一致しません。")
      return
    }

    setIsLoading(true)

    try {
      // 認証サービスを呼び出してサインアップ処理
      await signup(email, password)

      // サインアップ成功後、ログインページにリダイレクト
      navigate("/login?registered=true")
    } catch (err: any) {
      // エラーメッセージを表示
      setError(err.message || "アカウント作成に失敗しました。")
    } finally {
      setIsLoading(false)
    }
  }

  // 管理者チェック中はローディング表示
  if (isCheckingAdmin) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-gray-50">
        <div className="text-center">
          <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
          <p className="mt-4 text-gray-600">管理者情報を確認中...</p>
        </div>
      </div>
    )
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gray-50 px-4 py-12 sm:px-6 lg:px-8">
      <div className="w-full max-w-md">
        <div className="mb-8 text-center">
          <h1 className="text-3xl font-bold">認証基盤管理</h1>
          <p className="mt-2 text-gray-600">管理者アカウントを作成してください</p>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>サインアップ</CardTitle>
            <CardDescription>新しい管理者アカウントを作成します</CardDescription>
          </CardHeader>
          <form onSubmit={handleSubmit}>
            <CardContent className="space-y-4">
              {error && (
                <Alert variant="destructive">
                  <AlertCircle className="h-4 w-4" />
                  <AlertDescription>{error}</AlertDescription>
                </Alert>
              )}

              <div className="space-y-2">
                <Label htmlFor="email">メールアドレス</Label>
                <Input
                  id="email"
                  type="email"
                  placeholder="your@email.com"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">パスワード</Label>
                <Input
                  id="password"
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
                <p className="text-xs text-gray-500">8文字以上で、英字、数字、記号を含めてください</p>
              </div>

              <div className="space-y-2">
                <Label htmlFor="confirm-password">パスワード（確認）</Label>
                <Input
                  id="confirm-password"
                  type="password"
                  value={confirmPassword}
                  onChange={(e) => setConfirmPassword(e.target.value)}
                  required
                />
              </div>
            </CardContent>
            <CardFooter className="flex flex-col space-y-4">
              <Button type="submit" className="w-full" disabled={isLoading}>
                {isLoading ? (
                  <>
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                    アカウント作成中...
                  </>
                ) : (
                  "アカウント作成"
                )}
              </Button>
              <div className="text-center text-sm">
                <Link to="/login" className="text-blue-600 hover:text-blue-500">
                  ログインページに戻る
                </Link>
              </div>
            </CardFooter>
          </form>
        </Card>
      </div>
    </div>
  )
}
