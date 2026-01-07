"use client"

import { useState } from "react"
import { LabelTable } from "../../components/dashboard/label-table"
import { PageHeader } from "../../components/dashboard/page-header"
import { LabelCreateButton } from "../../components/dashboard/label-create-button"

export default function LabelsPage() {
  const [refreshKey, setRefreshKey] = useState(0)

  const handleLabelCreated = () => {
    // ラベルが作成されたらrefreshKeyを更新して再レンダリングを促す
    setRefreshKey((prevKey) => prevKey + 1)
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <PageHeader title="ラベル管理" description="ユーザーに割り当てるラベルを管理します" />
        <LabelCreateButton onLabelCreated={handleLabelCreated} />
      </div>
      <LabelTable key={refreshKey} />
    </div>
  )
}
