"use client"

import { createContext, useContext, useState, useEffect, type ReactNode } from "react"
import { isAuthenticated, getCurrentUser, logout, type AuthUser } from "../services/auth-service"

interface AuthContextType {
  user: AuthUser | null
  loading: boolean
  error: string | null
  login: (email: string, password: string) => Promise<void>
  logout: () => Promise<void>
  signup: (email: string, password: string) => Promise<void>
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

export function useAuth() {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider")
  }
  return context
}

interface AuthProviderProps {
  children: ReactNode
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<AuthUser | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const checkAuth = async () => {
      try {
        setLoading(true)
        const authenticated = await isAuthenticated()
        if (authenticated) {
          const currentUser = await getCurrentUser()
          setUser(currentUser)
        }
      } catch (err) {
        console.error("認証チェックに失敗しました:", err)
        setError("認証に失敗しました")
      } finally {
        setLoading(false)
      }
    }

    checkAuth()
  }, [])

  const loginUser = async (email: string, password: string) => {
    try {
      setLoading(true)
      setError(null)
      const user = await import("../services/auth-service").then((module) => module.login(email, password))
      setUser(user)
    } catch (err: any) {
      setError(err.message || "ログインに失敗しました")
      throw err
    } finally {
      setLoading(false)
    }
  }

  const signupUser = async (email: string, password: string) => {
    try {
      setLoading(true)
      setError(null)
      await import("../services/auth-service").then((module) => module.signup(email, password))
    } catch (err: any) {
      setError(err.message || "サインアップに失敗しました")
      throw err
    } finally {
      setLoading(false)
    }
  }

  const logoutUser = async () => {
    try {
      setLoading(true)
      await logout()
      setUser(null)
    } catch (err: any) {
      setError(err.message || "ログアウトに失敗しました")
    } finally {
      setLoading(false)
    }
  }

  const value = {
    user,
    loading,
    error,
    login: loginUser,
    logout: logoutUser,
    signup: signupUser,
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}
