import { useState, useEffect } from "react";
import type { User } from "../lib/types";
import { getUsers, toggleUserBan, deleteUser } from "../services/user-service";

export function useUsers() {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await getUsers();
      setUsers(data);
    } catch (err) {
      setError("Failed to fetch users");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchUsers();
  }, []);

  const handleToggleBan = async (userId: string) => {
    try {
      await toggleUserBan(userId);
      setUsers(users.map(u => u.id === userId ? { ...u, banned: !u.banned } : u));
    } catch (err) {
      console.error("Failed to toggle ban:", err);
    }
  };

  const confirmDelete = async (userId: string) => {
    try {
      await deleteUser(userId);
      setUsers(users.filter(u => u.id !== userId));
    } catch (err) {
      console.error("Failed to delete user:", err);
      throw err;
    }
  };

  const addUser = (user: User) => {
    setUsers([user, ...users]);
  };

  const updateUserInfo = (updated: User) => {
    setUsers(users.map(u => u.id === updated.id ? updated : u));
  };

  return {
    users,
    loading,
    error,
    handleToggleBan,
    confirmDelete,
    addUser,
    updateUserInfo,
    refresh: fetchUsers
  };
}
