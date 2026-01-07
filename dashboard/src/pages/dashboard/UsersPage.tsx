import { UserTable } from "../../components/dashboard/user-table"
import { PageHeader } from "../../components/dashboard/page-header"
import { UserCreateButton } from "../../components/dashboard/user-create-button"
import { LoginAsUserBadge } from "../../components/dashboard/login-as-user-badge"

export default function UsersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <PageHeader title="ユーザー管理" description="システムに登録されているユーザーを管理します" />
        <div className="flex items-center gap-4">
          <LoginAsUserBadge />
          <UserCreateButton />
        </div>
      </div>
      <UserTable />
    </div>
  )
}
