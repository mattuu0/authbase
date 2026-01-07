import { useState, useEffect } from "react";
import { 
  Search, 
  MoreHorizontal, 
  UserX, 
  UserCheck, 
  Trash2, 
  Edit2,
  ExternalLink,
  Shield,
  ShieldAlert
} from "lucide-react";
import type { User } from "../lib/types";
import { getUsers, toggleUserBan, deleteUser } from "../services/user-service";
import { cn } from "../lib/utils";
import { UserEditModal } from "../components/UserEditModal";
import { UserDeleteModal } from "../components/UserDeleteModal";
import { UserCreateModal } from "../components/UserCreateModal";

export default function UsersPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [editingUser, setEditingUser] = useState<User | null>(null);
  const [deletingUser, setDeletingUser] = useState<User | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  useEffect(() => {
    fetchUsers();
  }, []);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await getUsers();
      setUsers(data);
    } catch (error) {
      console.error("Failed to fetch users:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleToggleBan = async (userId: string) => {
    try {
      await toggleUserBan(userId);
      setUsers(users.map(u => u.id === userId ? { ...u, banned: !u.banned } : u));
    } catch (error) {
      console.error("Failed to toggle ban:", error);
    }
  };

  const confirmDelete = async (userId: string) => {
    try {
      await deleteUser(userId);
      setUsers(users.filter(u => u.id !== userId));
    } catch (error) {
      console.error("Failed to delete user:", error);
      throw error; // Let the modal handle error UI
    }
  };

  const filteredUsers = users.filter(user => 
    user.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    user.email.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900">ユーザー管理</h2>
        <button 
          onClick={() => setIsCreateModalOpen(true)}
          className="inline-flex items-center justify-center rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 transition-colors"
        >
          ユーザーを追加
        </button>
      </div>

      <div className="flex items-center gap-2 rounded-lg border bg-white px-3 py-2 shadow-sm focus-within:ring-2 focus-within:ring-blue-500 transition-all">
        <Search className="h-5 w-5 text-gray-400" />
        <input
          type="text"
          placeholder="名前、メールアドレスで検索..."
          className="flex-1 border-none bg-transparent text-sm outline-none placeholder:text-gray-400"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <div className="overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
        <div className="overflow-x-auto">
          <table className="w-full text-left text-sm">
            <thead className="border-b bg-gray-50 text-xs font-semibold uppercase text-gray-500">
              <tr>
                <th className="px-6 py-4">ユーザー</th>
                <th className="px-6 py-4">プロバイダ</th>
                <th className="px-6 py-4">ラベル</th>
                <th className="px-6 py-4">作成日</th>
                <th className="px-6 py-4">ステータス</th>
                <th className="px-6 py-4 text-right">アクション</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-gray-200">
              {loading ? (
                <tr>
                  <td colSpan={6} className="px-6 py-10 text-center text-gray-500">
                    読み込み中...
                  </td>
                </tr>
              ) : filteredUsers.length === 0 ? (
                <tr>
                  <td colSpan={6} className="px-6 py-10 text-center text-gray-500">
                    ユーザーが見つかりませんでした。
                  </td>
                </tr>
              ) : (
                filteredUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50 transition-colors">
                    <td className="px-6 py-4">
                      <div className="flex items-center gap-3">
                        <img
                          src={user.avatar}
                          alt={user.name}
                          className="h-10 w-10 rounded-full border border-gray-100"
                        />
                        <div>
                          <div className="font-medium text-gray-900">{user.name}</div>
                          <div className="text-xs text-gray-500">{user.email}</div>
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4">
                      <span className="inline-flex items-center rounded-full bg-gray-100 px-2.5 py-0.5 text-xs font-medium text-gray-800">
                        {user.provider}
                      </span>
                    </td>
                    <td className="px-6 py-4">
                      <div className="flex flex-wrap gap-1">
                        {user.labels.map((label) => (
                          <span
                            key={label}
                            className="inline-flex items-center rounded-full bg-blue-50 px-2 py-0.5 text-xs font-medium text-blue-700"
                          >
                            {label}
                          </span>
                        ))}
                      </div>
                    </td>
                    <td className="px-6 py-4 text-gray-500">{user.createdAt}</td>
                    <td className="px-6 py-4">
                      {user.banned ? (
                        <span className="inline-flex items-center gap-1 text-red-600 font-medium">
                          <ShieldAlert className="h-4 w-4" />
                          BAN済み
                        </span>
                      ) : (
                        <span className="inline-flex items-center gap-1 text-green-600 font-medium">
                          <Shield className="h-4 w-4" />
                          アクティブ
                        </span>
                      )}
                    </td>
                    <td className="px-6 py-4 text-right">
                      <div className="flex justify-end gap-2">
                        <button
                          onClick={() => setEditingUser(user)}
                          className="rounded-md p-2 text-blue-600 hover:bg-blue-50 transition-colors"
                          title="編集"
                        >
                          <Edit2 className="h-4 w-4" />
                        </button>
                        <button
                          onClick={() => handleToggleBan(user.id)}
                          className={cn(
                            "rounded-md p-2 transition-colors",
                            user.banned ? "text-green-600 hover:bg-green-50" : "text-amber-600 hover:bg-amber-50"
                          )}
                          title={user.banned ? "制限解除" : "BAN"}
                        >
                          {user.banned ? <UserCheck className="h-4 w-4" /> : <UserX className="h-4 w-4" />}
                        </button>
                        <button
                          onClick={() => setDeletingUser(user)}
                          className="rounded-md p-2 text-red-600 hover:bg-red-50 transition-colors"
                          title="削除"
                        >
                          <Trash2 className="h-4 w-4" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>

      <UserEditModal
        user={editingUser}
        isOpen={!!editingUser}
        onClose={() => setEditingUser(null)}
        onUpdate={(updated) => {
          setUsers(users.map((u) => (u.id === updated.id ? updated : u)));
        }}
        onDelete={(userId) => {
          const userToDelete = users.find(u => u.id === userId);
          if (userToDelete) {
            setEditingUser(null);
            setDeletingUser(userToDelete);
          }
        }}
      />

      <UserDeleteModal
        user={deletingUser}
        isOpen={!!deletingUser}
        onClose={() => setDeletingUser(null)}
        onConfirm={() => confirmDelete(deletingUser!.id)}
      />

      <UserCreateModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCreated={(newUser) => {
          setUsers([newUser, ...users]);
        }}
      />
    </div>
  );
}