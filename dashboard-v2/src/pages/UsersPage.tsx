export default function UsersPage() {
  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900">ユーザー管理</h2>
        <button className="inline-flex items-center rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-blue-600">
          ユーザーを追加
        </button>
      </div>
      <div className="rounded-lg border bg-white shadow-sm">
        <div className="p-6">
          <p className="text-gray-500">ユーザー一覧がここに表示されます。</p>
        </div>
      </div>
    </div>
  );
}
