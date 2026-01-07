"use client"

import { useState, useEffect } from "react"
import { Link, useLocation, useNavigate } from "react-router-dom"
import { cn } from "../../lib/utils"
import { Button } from "../ui/button"
import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip"
import { LogOut, Users, Tag, Settings, Menu, X, History } from "lucide-react"
import { logout, getCurrentUser } from "../../services/auth-service"
import { Sheet, SheetContent, SheetTrigger } from "../ui/sheet"

const navItems = [
  {
    title: "ユーザー管理",
    href: "/dashboard/users",
    icon: Users,
  },
  {
    title: "ラベル管理",
    href: "/dashboard/labels",
    icon: Tag,
  },
  {
    title: "プロバイダ設定",
    href: "/dashboard/providers",
    icon: Settings,
  },
  {
    title: "セッション管理",
    href: "/dashboard/sessions",
    icon: History,
  },
]

export function Sidebar() {
  const location = useLocation()
  const navigate = useNavigate()
  const [isCollapsed, setIsCollapsed] = useState(false)
  const [isMobile, setIsMobile] = useState(false)
  const [currentUser, setCurrentUser] = useState<any>(null)
  const [mobileOpen, setMobileOpen] = useState(false)

  // 現在のユーザー情報を取得
  useEffect(() => {
    const fetchUserInfo = async () => {
      try {
        const user = await getCurrentUser()
        setCurrentUser(user)
      } catch (error) {
        console.error("ユーザー情報の取得に失敗しました:", error)
      }
    }

    fetchUserInfo()
  }, [])

  // 画面サイズに応じてモバイル表示を切り替え
  useEffect(() => {
    const checkIsMobile = () => {
      const mobile = window.innerWidth < 768
      setIsMobile(mobile)
      if (mobile && !isCollapsed) {
        setIsCollapsed(true)
      }
    }

    // 初期チェック
    checkIsMobile()

    // リサイズイベントのリスナーを追加
    window.addEventListener("resize", checkIsMobile)

    // クリーンアップ
    return () => window.removeEventListener("resize", checkIsMobile)
  }, [isCollapsed])

  // ログアウト処理
  const handleLogout = async () => {
    try {
      await logout()
      navigate("/login")
    } catch (error) {
      console.error("ログアウトに失敗しました:", error)
    }
  }

  // モバイル用のサイドバー
  const MobileSidebar = () => (
    <Sheet open={mobileOpen} onOpenChange={setMobileOpen}>
      <SheetTrigger asChild>
        <Button variant="ghost" size="icon" className="md:hidden">
          <Menu className="h-5 w-5" />
          <span className="sr-only">メニューを開く</span>
        </Button>
      </SheetTrigger>
      <SheetContent side="left" className="p-0 w-64">
        <div className="flex h-full flex-col">
          <div className="flex h-14 items-center border-b px-4">
            <h1 className="font-semibold text-sm">authkit</h1>
            <Button variant="ghost" size="icon" className="ml-auto" onClick={() => setMobileOpen(false)}>
              <X className="h-5 w-5" />
            </Button>
          </div>
          <nav className="flex-1 overflow-auto py-4">
            <ul className="grid gap-1 px-2">
              {navItems.map((item) => {
                const isActive = location.pathname === item.href || location.pathname.startsWith(`${item.href}/`)
                return (
                  <li key={item.href}>
                    <Link
                      to={item.href}
                      className={cn(
                        "flex h-10 items-center rounded-md px-3 text-sm font-medium transition-colors",
                        isActive ? "bg-blue-50 text-blue-700" : "text-gray-500 hover:bg-gray-100 hover:text-gray-900",
                      )}
                      onClick={() => setMobileOpen(false)}
                    >
                      <item.icon className={cn("h-5 w-5 mr-3", isActive ? "text-blue-700" : "text-gray-500")} />
                      <span>{item.title}</span>
                    </Link>
                  </li>
                )
              })}
            </ul>
          </nav>
          <div className="border-t p-2">
            <Button
              variant="ghost"
              className="w-full justify-start text-gray-500 hover:bg-gray-100 hover:text-gray-900"
              onClick={() => {
                handleLogout()
                setMobileOpen(false)
              }}
            >
              <LogOut className="h-5 w-5 mr-3" />
              <span>ログアウト</span>
            </Button>
          </div>
        </div>
      </SheetContent>
    </Sheet>
  )

  // デスクトップ用のサイドバー
  return (
    <>
      {/* モバイル用ハンバーガーメニュー */}
      <div className="md:hidden fixed top-0 left-0 z-20 p-4">
        <MobileSidebar />
      </div>

      {/* デスクトップ用サイドバー */}
      <TooltipProvider>
        <div
          className={cn(
            "hidden md:flex h-full flex-col border-r bg-white shadow-sm transition-all duration-300 ease-in-out",
            isCollapsed ? "w-16" : "w-56",
          )}
          style={{ minWidth: isCollapsed ? "4rem" : "14rem" }}
        >
          <div className="flex h-14 items-center border-b px-3">
            <h1
              className={cn(
                "font-semibold text-sm transition-opacity duration-200 whitespace-nowrap overflow-hidden",
                isCollapsed ? "opacity-0 w-0" : "opacity-100 w-auto",
              )}
            >
              authkit
            </h1>
            <Button
              variant="ghost"
              size="icon"
              className={cn("transition-all duration-300", isCollapsed ? "ml-auto" : "ml-auto")}
              onClick={() => setIsCollapsed(!isCollapsed)}
              aria-label={isCollapsed ? "展開" : "折りたたむ"}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                strokeWidth="2"
                strokeLinecap="round"
                strokeLinejoin="round"
                className={cn("transition-transform duration-200", isCollapsed ? "rotate-180" : "")}
              >
                <path d="m15 6-6 6 6 6" />
              </svg>
            </Button>
          </div>
          <nav className="flex-1 overflow-auto py-4">
            <ul className="grid gap-1 px-2">
              {navItems.map((item) => {
                const isActive = location.pathname === item.href || location.pathname.startsWith(`${item.href}/`)
                return (
                  <li key={item.href}>
                    <Tooltip delayDuration={0}>
                      <TooltipTrigger asChild>
                        <Link
                          to={item.href}
                          className={cn(
                            "flex h-10 items-center rounded-md px-3 text-sm font-medium transition-colors",
                            isActive
                              ? "bg-blue-50 text-blue-700"
                              : "text-gray-500 hover:bg-gray-100 hover:text-gray-900",
                          )}
                        >
                          <item.icon className={cn("h-5 w-5", isActive ? "text-blue-700" : "text-gray-500")} />
                          <span
                            className={cn(
                              "ml-3 transition-all duration-200 whitespace-nowrap overflow-hidden",
                              isCollapsed ? "w-0 opacity-0" : "w-auto opacity-100",
                            )}
                          >
                            {item.title}
                          </span>
                        </Link>
                      </TooltipTrigger>
                      {isCollapsed && <TooltipContent side="right">{item.title}</TooltipContent>}
                    </Tooltip>
                  </li>
                )
              })}
            </ul>
          </nav>

          <div className="border-t p-2">
            <Tooltip delayDuration={0}>
              <TooltipTrigger asChild>
                <Button
                  variant="ghost"
                  className={cn(
                    "w-full justify-start text-gray-500 hover:bg-gray-100 hover:text-gray-900",
                    isCollapsed ? "px-3" : "",
                  )}
                  onClick={handleLogout}
                >
                  <LogOut className="h-5 w-5" />
                  <span
                    className={cn(
                      "ml-3 transition-all duration-200 whitespace-nowrap overflow-hidden",
                      isCollapsed ? "w-0 opacity-0" : "w-auto opacity-100",
                    )}
                  >
                    ログアウト
                  </span>
                </Button>
              </TooltipTrigger>
              {isCollapsed && <TooltipContent side="right">ログアウト</TooltipContent>}
            </Tooltip>
          </div>
        </div>
      </TooltipProvider>
    </>
  )
}
