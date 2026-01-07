import { useState, useEffect } from "react";
import { History, Smartphone, Monitor, Trash2, Clock, MapPin } from "lucide-react";
import type { Session } from "../lib/types";
import { getSessions, deleteSession } from "../services/session-service";
import { cn } from "../lib/utils";

export default function SessionsPage() {
  const [sessions, setSessions] = useState<Session[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchSessions();
  }, []);

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

  const handleDelete = async (id: string) => {
    if (!confirm("このセッションを終了してもよろしいですか？")) return;
    try {
      await deleteSession(id);
      setSessions(sessions.filter((s) => s.id !== id));
    } catch (error) {
      console.error("Failed to delete session:", error);
    }
  };

  const getDeviceIcon = (userAgent: string) => {
    if (userAgent.toLowerCase().includes("mobile") || userAgent.toLowerCase().includes("android") || userAgent.toLowerCase().includes("iphone")) {
      return <Smartphone className="h-5 w-5" />;
    }
    return <Monitor className="h-5 w-5" />;
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900">セッション管理</h2>
      </div>

      <div className="grid gap-4">
        {loading ? (
          <div className="py-10 text-center text-gray-500">読み込み中...</div>
        ) : sessions.length === 0 ? (
          <div className="py-10 text-center text-gray-500">アクティブなセッションはありません。</div>
        ) : (
          sessions.map((session) => (
            <div
              key={session.id}
              className="flex items-center justify-between rounded-xl border border-gray-200 bg-white p-6 shadow-sm hover:border-blue-100 transition-all"
            >
              <div className="flex items-center gap-4">
                <div className="flex h-12 w-12 items-center justify-center rounded-full bg-blue-50 text-blue-600">
                  {getDeviceIcon(session.userAgent)}
                </div>
                <div>
                  <div className="flex items-center gap-2">
                    <h3 className="font-semibold text-gray-900">{session.ipAddress}</h3>
                    {session.isActive && (
                      <span className="inline-flex items-center rounded-full bg-green-100 px-2 py-0.5 text-xs font-medium text-green-700">
                        オンライン
                      </span>
                    )}
                  </div>
                  <div className="mt-1 flex flex-wrap gap-x-4 gap-y-1 text-sm text-gray-500">
                    <span className="flex items-center gap-1">
                      <Clock className="h-3 w-3" />
                      作成: {new Date(session.createdAt).toLocaleString()}
                    </span>
                    <span className="flex items-center gap-1">
                      <MapPin className="h-3 w-3" />
                      ユーザーID: {session.userId}
                    </span>
                  </div>
                  <p className="mt-1 text-xs text-gray-400 truncate max-w-md">
                    {session.userAgent}
                  </p>
                </div>
              </div>
              <button
                onClick={() => handleDelete(session.id)}
                className="inline-flex items-center gap-2 rounded-md border border-red-200 px-3 py-1.5 text-sm font-medium text-red-600 hover:bg-red-50 transition-colors"
              >
                <Trash2 className="h-4 w-4" />
                セッションを終了
              </button>
            </div>
          ))
        )}
      </div>
    </div>
  );
}
