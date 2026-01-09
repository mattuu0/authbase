import { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { 
  Search, 
  Filter, 
  Smartphone, 
  Monitor, 
  Trash2, 
  Clock, 
  MapPin, 
  User as UserIcon,
  ChevronDown,
  History,
  Users
} from "lucide-react";
import type { Session, User } from "../lib/types";
import { getSessions, deleteSession } from "../services/session-service";
import { getUsers } from "../services/user-service";
import { cn } from "../lib/utils";
import { SessionDeleteModal } from "../components/SessionDeleteModal";

export default function SessionsPage() {
  const [searchParams] = useSearchParams();
  const userIdParam = searchParams.get("userId");

  const [sessions, setSessions] = useState<Session[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [statusFilter, setStatusFilter] = useState<"all" | "active" | "inactive">("all");
  const [selectedUserId, setSelectedUserId] = useState<string>(userIdParam || "all");
  const [deletingSession, setDeletingSession] = useState<Session | null>(null);

  useEffect(() => {
    fetchSessions();
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    try {
      const data = await getUsers();
      setUsers(data);
    } catch (error) {
      console.error("Failed to fetch users:", error);
    }
  };

  const fetchSessions = async () => {
    setLoading(true);
    try {
      const data = await getSessions();
      setSessions(data);
    } catch (error) {
      console.error("Failed to fetch sessions:", error);
    } finally {
      setLoading(false);
    }
  };

  const confirmDelete = async (id: string) => {
    try {
      await deleteSession(id);
      setSessions(sessions.filter((s) => s.id !== id));
    } catch (error) {
      console.error("Failed to delete session:", error);
      throw error;
    }
  };

  const filteredSessions = sessions.filter((session) => {
// ... (omitting middle part for clarity in search, will replace the whole block if needed)
    const matchesStatus = 
      statusFilter === "all" || 
      (statusFilter === "active" && session.isActive) || 
      (statusFilter === "inactive" && !session.isActive);
    
    const matchesUser = 
      selectedUserId === "all" || 
      session.userId === selectedUserId;
    
    const searchLower = searchQuery.toLowerCase();
    const matchesSearch = 
      session.userId.toLowerCase().includes(searchLower) ||
      session.userName.toLowerCase().includes(searchLower) ||
      session.userEmail.toLowerCase().includes(searchLower) ||
      session.ipAddress.toLowerCase().includes(searchLower);

    return matchesStatus && matchesUser && matchesSearch;
  });

  const getDeviceIcon = (userAgent: string) => {
    if (userAgent.toLowerCase().includes("mobile") || userAgent.toLowerCase().includes("android") || userAgent.toLowerCase().includes("iphone")) {
      return <Smartphone className="h-5 w-5" />;
    }
    return <Monitor className="h-5 w-5" />;
  };

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900">セッション管理</h2>
        <div className="flex items-center gap-2 text-sm text-gray-500">
          <History className="h-4 w-4" />
          <span>全 {sessions.length} セッション</span>
        </div>
      </div>

      <div className="grid gap-4 md:grid-cols-12">
        <div className="md:col-span-6 flex items-center gap-2 rounded-lg border bg-white px-3 py-2 shadow-sm focus-within:ring-2 focus-within:ring-blue-500 transition-all">
          <Search className="h-5 w-5 text-gray-400" />
          <input
            type="text"
            placeholder="メール、IPアドレスで検索..."
            className="flex-1 border-none bg-transparent text-sm outline-none placeholder:text-gray-400"
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        
        <div className="md:col-span-3">
          <div className="relative">
            <select
              value={selectedUserId}
              onChange={(e) => setSelectedUserId(e.target.value)}
              className="w-full appearance-none rounded-lg border bg-white pl-10 pr-10 py-2 text-sm font-medium text-gray-700 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 hover:bg-gray-50 transition-all cursor-pointer"
            >
              <option value="all">すべてのユーザー</option>
              {users.map((user) => (
                <option key={user.id} value={user.id}>
                  {user.name}
                </option>
              ))}
            </select>
            <Users className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 pointer-events-none" />
            <ChevronDown className="absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 pointer-events-none" />
          </div>
        </div>

        <div className="md:col-span-3">
          <div className="relative">
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value as any)}
              className="w-full appearance-none rounded-lg border bg-white pl-10 pr-10 py-2 text-sm font-medium text-gray-700 shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 hover:bg-gray-50 transition-all cursor-pointer"
            >
              <option value="all">すべてのステータス</option>
              <option value="active">オンラインのみ</option>
              <option value="inactive">オフラインのみ</option>
            </select>
            <Filter className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 pointer-events-none" />
            <ChevronDown className="absolute right-3 top-1/2 h-4 w-4 -translate-y-1/2 text-gray-400 pointer-events-none" />
          </div>
        </div>
      </div>

      <div className="grid gap-4">
        {loading ? (
          <div className="py-20 text-center">
            <div className="inline-block h-8 w-8 animate-spin rounded-full border-4 border-solid border-blue-600 border-r-transparent"></div>
            <p className="mt-4 text-gray-500 font-medium">読み込み中...</p>
          </div>
        ) : filteredSessions.length === 0 ? (
          <div className="rounded-2xl border-2 border-dashed border-gray-200 bg-white py-20 text-center">
            <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-full bg-gray-50 text-gray-400">
              <History className="h-8 w-8" />
            </div>
            <h3 className="mt-4 text-lg font-semibold text-gray-900">セッションが見つかりません</h3>
            <p className="mt-2 text-sm text-gray-500">検索条件やフィルターを変更してみてください。</p>
          </div>
        ) : (
          filteredSessions.map((session) => (
            <div
              key={session.id}
              className="group flex flex-col sm:flex-row sm:items-center justify-between gap-4 rounded-xl border border-gray-200 bg-white p-5 shadow-sm hover:border-blue-200 hover:shadow-md transition-all"
            >
              <div className="flex items-start gap-4">
                <div className={cn(
                  "flex h-12 w-12 shrink-0 items-center justify-center rounded-xl transition-colors",
                  session.isActive ? "bg-green-50 text-green-600 group-hover:bg-green-100" : "bg-gray-100 text-gray-500 group-hover:bg-gray-200"
                )}>
                  {getDeviceIcon(session.userAgent)}
                </div>
                <div className="min-w-0 space-y-1">
                  <div className="flex flex-wrap items-center gap-2">
                    <h3 className="font-bold text-gray-900">{session.ipAddress}</h3>
                    {session.isActive ? (
                      <span className="inline-flex items-center gap-1 rounded-full bg-green-100 px-2.5 py-0.5 text-xs font-bold text-green-700">
                        <span className="h-1.5 w-1.5 rounded-full bg-green-500 animate-pulse"></span>
                        オンライン
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-bold text-gray-500">
                        オフライン
                      </span>
                    )}
                  </div>
                  
                  <div className="flex flex-wrap gap-x-4 gap-y-1.5">
                    <div className="flex items-center gap-2 text-sm text-gray-600">
                      <UserIcon className="h-4 w-4 text-gray-400" />
                      <span className="font-bold text-gray-900">{session.userName}</span>
                      <span className="text-gray-500 text-xs">{session.userEmail}</span>
                      <span className="text-gray-300 text-[10px] font-mono select-all" title="User ID">UID: {session.userId}</span>
                    </div>
                  </div>

                  <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-gray-400">
                    <span className="flex items-center gap-1">
                      <Clock className="h-3.5 w-3.5" />
                      作成: {new Date(session.createdAt).toLocaleString("ja-JP")}
                    </span>
                    <span className="flex items-center gap-1">
                      <Clock className="h-3.5 w-3.5 text-amber-400" />
                      期限: {new Date(session.expiresAt).toLocaleString("ja-JP")}
                    </span>
                    <span className="flex items-center gap-1">
                      <MapPin className="h-3.5 w-3.5" />
                      IP: {session.ipAddress}
                    </span>
                  </div>
                  
                  <div className="text-[10px] font-mono text-gray-400 bg-gray-50 px-2 py-0.5 rounded border border-gray-100 inline-block">
                    SID: {session.id}
                  </div>
                  
                  <p className="text-xs text-gray-400 break-all" title={session.userAgent}>
                    {session.userAgent}
                  </p>
                </div>
              </div>
              
              <div className="flex shrink-0 items-center gap-3">
                <button
                  onClick={() => setDeletingSession(session)}
                  className="flex flex-1 sm:flex-none items-center justify-center gap-2 rounded-lg border border-red-100 px-4 py-2 text-sm font-bold text-red-600 hover:bg-red-50 hover:border-red-200 transition-all active:scale-[0.98]"
                >
                  <Trash2 className="h-4 w-4" />
                  終了
                </button>
              </div>
            </div>
          ))
        )}
      </div>

      <SessionDeleteModal
        session={deletingSession}
        isOpen={!!deletingSession}
        onClose={() => setDeletingSession(null)}
        onConfirm={() => confirmDelete(deletingSession!.id)}
      />
    </div>
  );
}
