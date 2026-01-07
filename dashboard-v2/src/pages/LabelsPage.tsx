import { useState, useEffect } from "react";
import { Plus, Trash2, Tag, Search } from "lucide-react";
import type { Label } from "../lib/types";
import { getLabels, deleteLabel } from "../services/label-service";
import { LabelCreateModal } from "../components/LabelCreateModal";

export default function LabelsPage() {
  const [labels, setLabels] = useState<Label[]>([]);
  const [loading, setLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);

  useEffect(() => {
    fetchLabels();
  }, []);

  const fetchLabels = async () => {
    setLoading(true);
    try {
      const data = await getLabels();
      setLabels(data);
    } catch (error) {
      console.error("Failed to fetch labels:", error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm("このラベルを削除してもよろしいですか？")) return;
    try {
      await deleteLabel(id);
      setLabels(labels.filter((l) => l.id !== id));
    } catch (error) {
      console.error("Failed to delete label:", error);
    }
  };

  const filteredLabels = labels.filter((label) =>
    label.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-3xl font-bold tracking-tight text-gray-900">ラベル管理</h2>
        <button 
          onClick={() => setIsCreateModalOpen(true)}
          className="inline-flex items-center gap-2 rounded-md bg-blue-600 px-4 py-2 text-sm font-semibold text-white shadow-sm hover:bg-blue-500 transition-colors"
        >
          <Plus className="h-4 w-4" />
          ラベルを作成
        </button>
      </div>

      <div className="flex items-center gap-2 rounded-lg border bg-white px-3 py-2 shadow-sm focus-within:ring-2 focus-within:ring-blue-500 transition-all">
        <Search className="h-5 w-5 text-gray-400" />
        <input
          type="text"
          placeholder="ラベル名で検索..."
          className="flex-1 border-none bg-transparent text-sm outline-none placeholder:text-gray-400"
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
        />
      </div>

      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {loading ? (
          <div className="col-span-full py-10 text-center text-gray-500">読み込み中...</div>
        ) : filteredLabels.length === 0 ? (
          <div className="col-span-full py-10 text-center text-gray-500">ラベルが見つかりませんでした。</div>
        ) : (
          filteredLabels.map((label) => (
            <div
              key={label.id}
              className="group relative flex items-center justify-between rounded-xl border border-gray-200 bg-white p-4 shadow-sm hover:border-blue-200 hover:shadow-md transition-all"
            >
              <div className="flex items-center gap-3">
                <div
                  className="flex h-10 w-10 items-center justify-center rounded-lg"
                  style={{ backgroundColor: `${label.color}20`, color: label.color }}
                >
                  <Tag className="h-5 w-5" />
                </div>
                <div>
                  <h3 className="font-semibold text-gray-900">{label.name}</h3>
                  <p className="text-xs text-gray-500">作成日: {label.createdAt}</p>
                </div>
              </div>
              <button
                onClick={() => handleDelete(label.id)}
                className="rounded-md p-2 text-gray-400 hover:bg-red-50 hover:text-red-600 opacity-0 group-hover:opacity-100 transition-all"
                title="削除"
              >
                <Trash2 className="h-4 w-4" />
              </button>
            </div>
          ))
        )}
      </div>

      <LabelCreateModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onCreated={(newLabel) => {
          setLabels([newLabel, ...labels]);
        }}
      />
    </div>
  );
}
