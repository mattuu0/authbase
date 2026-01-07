import { Outlet } from "react-router-dom"
import { Sidebar } from "../components/dashboard/sidebar"

export default function DashboardLayout() {
  return (
    <div className="flex min-h-screen bg-gray-50">
      <Sidebar />
      <main className="flex-1 overflow-auto p-6 md:p-6 pt-16 md:pt-6">
        <Outlet />
      </main>
    </div>
  )
}
