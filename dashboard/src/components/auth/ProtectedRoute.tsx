"use client"

import { type ReactNode, useEffect, useState } from "react"
import { Navigate, useLocation } from "react-router-dom"
import { isAuthenticated } from "../../services/auth-service"

interface ProtectedRouteProps {
  children: ReactNode
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const location = useLocation()
  const [isChecking, setIsChecking] = useState(true)
  const [isAuth, setIsAuth] = useState(false)

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const authenticated = await isAuthenticated()
        setIsAuth(authenticated)
      } catch (error) {
        console.error("認証チェックに失敗しました:", error)
        setIsAuth(false)
      } finally {
        setIsChecking(false)
      }
    }

    checkAuth()
  }, [])

  // 認証チェック中はローディング表示
  if (isChecking) {
    return (
      <div className="flex h-screen items-center justify-center">
        <div className="h-8 w-8 animate-spin rounded-full border-4 border-blue-500 border-t-transparent"></div>
      </div>
    )
  }

  // 認証されていない場合はログインページにリダイレクト
  if (!isAuth) {
    return <Navigate to={`/login?redirect=${encodeURIComponent(location.pathname)}`} replace />
  }

  // 認証されている場合は子コンポーネントを表示
  return <>{children}</>
}
