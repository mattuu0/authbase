"use client"

import { useState, useEffect } from "react"
import { Badge } from "../ui/badge"
import { Button } from "../ui/button"
import { Users } from "lucide-react"
import { getCurrentLoginAs, clearLoginAs } from "../../services/user-service"
import { useNavigate } from "react-router-dom"

export function LoginAsUserBadge() {
  const [loginAsUser, setLoginAsUser] = useState<any>(null)
  const [loading, setLoading] = useState(true)
  const navigate = useNavigate()

  // ログイン中のユーザー情報を取得
  useEffect(() => {
    const fetchLoginAsUser = async () => {
      try {
        setLoading(true)
        const user = await getCurrentLoginAs()
        setLoginAsUser(user)
      } catch (error) {
        console.error("ログイン中のユーザー情報の取得に失敗しました:", error)
      } finally {
        setLoading(false)
      }
    }

    fetchLoginAsUser()

    // ストレージの変更を監視して状態を更新
    const handleStorageChange = () => {
      fetchLoginAsUser()
    }

    window.addEventListener("storage", handleStorageChange)

    // カスタムイベントを監視
    const handleLoginAsChange = () => {
      fetchLoginAsUser()
    }

    window.addEventListener("loginAsUserChanged", handleLoginAsChange)

    return () => {
      window.removeEventListener("storage", handleStorageChange)
      window.removeEventListener("loginAsUserChanged", handleLoginAsChange)
    }
  }, [])

  // ログイン中のユーザーをクリア
  const handleClearLoginAs = async () => {
    try {
      await clearLoginAs()
      setLoginAsUser(null)
      // 現在のページをリロード
      window.location.reload()
    } catch (error) {
      console.error("ログイン中のユーザーのクリアに失敗しました:", error)
    }
  }

  if (loading || !loginAsUser) {
    return null
  }

  return (
    <div className="flex items-center">
      <Badge
        variant="outline"
        className="bg-amber-50 text-amber-600 border-amber-200 px-3 py-1.5 flex items-center gap-2"
      >
        <Users className="h-3.5 w-3.5" />
        <span>{loginAsUser.name}としてログイン中</span>
        <Button variant="ghost" size="sm" className="h-6 px-2 ml-1 hover:bg-amber-100" onClick={handleClearLoginAs}>
          解除
        </Button>
      </Badge>
    </div>
  )
}
