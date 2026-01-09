import { Link, useLocation, useNavigate } from "react-router-dom";
import { Users, Tag, Settings, History, LogOut, Menu } from "lucide-react";
import { cn } from "../lib/utils";
import { logout } from "../services/auth-service";

const navItems = [
  { title: "ユーザー管理", href: "/dashboard/users", icon: Users },
  { title: "ラベル管理", href: "/dashboard/labels", icon: Tag },
  { title: "プロバイダ設定", href: "/dashboard/providers", icon: Settings },
  { title: "セッション管理", href: "/dashboard/sessions", icon: History },
];

export function Sidebar() {
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = async () => {
    await logout();
    navigate("/login");
  };

  return (
    <div className="flex h-screen w-64 flex-col border-r bg-white">
      <div className="flex h-16 items-center border-b px-6">
        <h1 className="text-xl font-bold text-gray-800">AuthBase</h1>
      </div>
      <nav className="flex-1 space-y-1 px-3 py-4">
        {navItems.map((item) => {
          const isActive = location.pathname.startsWith(item.href);
          return (
            <Link
              key={item.href}
              to={item.href}
              className={cn(
                "flex items-center rounded-lg px-3 py-2 text-sm font-medium transition-colors",
                isActive
                  ? "bg-blue-50 text-blue-700"
                  : "text-gray-600 hover:bg-gray-100 hover:text-gray-900"
              )}
            >
              <item.icon className={cn("mr-3 h-5 w-5", isActive ? "text-blue-700" : "text-gray-400")} />
              {item.title}
            </Link>
          );
        })}
      </nav>
      <div className="border-t p-4">
        <button 
          onClick={handleLogout}
          className="flex w-full items-center rounded-lg px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-100 hover:text-gray-900 transition-colors"
        >
          <LogOut className="mr-3 h-5 w-5 text-gray-400" />
          ログアウト
        </button>
      </div>
    </div>
  );
}
