import { ProviderSettings } from "../../components/dashboard/provider-settings"
import { PageHeader } from "../../components/dashboard/page-header"

export default function ProvidersPage() {
  return (
    <div className="space-y-6">
      <PageHeader title="プロバイダ設定" description="認証プロバイダの設定を管理します" />
      <ProviderSettings />
    </div>
  )
}
