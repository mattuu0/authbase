"use client"

import { useState, useEffect } from "react"
import { useSearchParams, useNavigate } from "react-router-dom"
import { PageHeader } from "../../components/dashboard/page-header"
import { SessionTable } from "../../components/dashboard/session-table"
import { Input } from "../../components/ui/input"
import { Button } from "../../components/ui/button"
import { Search, ArrowLeft } from "lucide-react"
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "../../components/ui/select"
import { getSessions, searchSessions, getSessionsByUserId } from "../../services/session-service"
import { getUsers } from "../../services/user-service"

export default function SessionsPage() {
  const [searchParams] = useSearchParams()
  const navigate = useNavigate()
  const userId = searchParams.get("userId")

  const [sessions, setSessions] = useState<any[]>([])
  const [searchTerm, setSearchTerm] = useState("")
  const [selectedUser, setSelectedUser] = useState<string>(userId || "all")
  const [users, setUsers] = useState<any[]>([])
  const [loading, setLoading] = useState(true)

  // ユーザーとセッションデータの取得
  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true)

        // ユーザーデータを取得
        const userData = await getUsers()
        setUsers(userData)

        // セッションデータを取得
        let sessionData
        if (userId) {
          sessionData = await getSessionsByUserId(userId)
          setSelectedUser(userId)
        } else {
          sessionData = await getSessions()
        }

        setSessions(sessionData)
      } catch (error) {
        console.error("データの取得に失敗しました:", error)
      } finally {
        setLoading(false)
      }
    }

    fetchData()
  }, [userId])

  // 検索処理
  useEffect(() => {
    const handleSearch = async () => {
      try {
        setLoading(true)

        let filteredSessions

        // ユーザーフィルターと検索語の両方が指定されている場合
        if (selectedUser !== "all" && searchTerm) {
          const allSessions = await searchSessions(searchTerm)
          filteredSessions = allSessions.filter((session) => session.userId === selectedUser)
        }
        // ユーザーフィルターのみ指定されている場合
        else if (selectedUser !== "all") {
          filteredSessions = await getSessionsByUserId(selectedUser)
        }
        // 検索語のみ指定されている場合
        else if (searchTerm) {
          filteredSessions = await searchSessions(searchTerm)
        }
        // どちらも指定されていない場合
        else {
          filteredSessions = await getSessions()
        }

        setSessions(filteredSessions)
      } catch (error) {
        console.error("セッション検索に失敗しました:", error)
      } finally {
        setLoading(false)
      }
    }

    // 検索処理を実行（実際のアプリではデバウンスを実装するとよい）
    const timer = setTimeout(handleSearch, 300)
    return () => clearTimeout(timer)
  }, [searchTerm, selectedUser])

  // 前のページに戻る
  const handleGoBack = () => {
    navigate(-1)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <PageHeader title="セッション管理" description="ユーザーのセッション情報を管理します" />
        <Button variant="outline" onClick={handleGoBack} className="flex items-center gap-2">
          <ArrowLeft className="h-4 w-4" />
          戻る
        </Button>
      </div>

      <div className="flex flex-col sm:flex-row gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-2.5 top-2.5 h-4 w-4 text-gray-500" />
          <Input
            placeholder="セッションID、IPアドレスで検索..."
            className="pl-8"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
          />
        </div>
        <div className="sm:w-64">
          <Select value={selectedUser} onValueChange={setSelectedUser}>
            <SelectTrigger>
              <SelectValue placeholder="ユーザーで絞り込み" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="all">すべてのユーザー</SelectItem>
              {users.map((user) => (
                <SelectItem key={user.id} value={user.id}>
                  {user.name} ({user.email})
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      <SessionTable sessions={sessions} loading={loading} users={users} />
    </div>
  )
}
