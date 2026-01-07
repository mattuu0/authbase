"use client"

import type React from "react"

import { useState, useEffect } from "react"
import { Button } from "../ui/button"
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from "../ui/card"
import { Input } from "../ui/input"
import { Label } from "../ui/label"
import { Switch } from "../ui/switch"
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs"
import { Loader2 } from "lucide-react"
import { getProviders, updateProviderSettings, toggleProvider } from "../../services/provider-service"

export function ProviderSettings() {
  const [providers, setProviders] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)
  const [activeTab, setActiveTab] = useState("google")
  const [formData, setFormData] = useState({
    clientId: "",
    clientSecret: "",
    callbackUrl: "",
  })

  // プロバイダデータの取得
  useEffect(() => {
    const fetchProviders = async () => {
      try {
        setLoading(true)
        const data = await getProviders()
        setProviders(data)

        // 初期タブのデータをフォームにセット
        if (data.length > 0) {
          const initialProvider = data.find((p) => p.ProviderCode === "google") || data[0]
          setActiveTab(initialProvider.ProviderCode)
          setFormData({
            clientId: initialProvider.ClientID,
            clientSecret: initialProvider.ClientSecret,
            callbackUrl: initialProvider.CallbackURL,
          })
        }
      } catch (error) {
        console.error("プロバイダデータの取得に失敗しました:", error)
      } finally {
        setLoading(false)
      }
    }

    fetchProviders()
  }, [])

  // タブ切り替え時にフォームデータを更新
  const handleTabChange = (value: string) => {
    setActiveTab(value)
    const provider = providers.find((p) => p.ProviderCode === value)
    if (provider) {
      setFormData({
        clientId: provider.ClientID,
        clientSecret: provider.ClientSecret,
        callbackUrl: provider.CallbackURL,
      })
    }
  }

  // フォーム入力の処理
  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target
    setFormData((prev) => ({ ...prev, [name]: value }))
  }

  // プロバイダの有効/無効を切り替え
  const handleToggleProvider = async (providerCode: string) => {
    try {
      const updatedProvider = await toggleProvider(providerCode)
      setProviders(
        providers.map((p) => (p.ProviderCode === providerCode ? { ...p, IsEnabled: updatedProvider.IsEnabled } : p)),
      )
    } catch (error) {
      console.error("プロバイダの有効/無効切り替えに失敗しました:", error)
    }
  }

  // 設定の保存
  const handleSave = async () => {
    try {
      setSaving(true)
      await updateProviderSettings(activeTab, {
        name: activeTab, // プロバイダ名
        clientId: formData.clientId,
        clientSecret: formData.clientSecret,
        redirectUri: formData.callbackUrl,
        scopes: "", // スコープは今回省略
      })

      // 保存成功後、プロバイダリストを更新
      const updatedProviders = providers.map((p) =>
        p.ProviderCode === activeTab
          ? {
              ...p,
              ClientID: formData.clientId,
              ClientSecret: formData.clientSecret,
              CallbackURL: formData.callbackUrl,
            }
          : p,
      )
      setProviders(updatedProviders)
    } catch (error) {
      console.error("プロバイダ設定の保存に失敗しました:", error)
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>プロバイダ設定</CardTitle>
          <CardDescription>外部認証プロバイダの設定を管理します</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center justify-center py-10">
            <Loader2 className="h-8 w-8 animate-spin text-primary" />
          </div>
        </CardContent>
      </Card>
    )
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>プロバイダ設定</CardTitle>
        <CardDescription>外部認証プロバイダの設定を管理します</CardDescription>
      </CardHeader>
      <CardContent>
        <Tabs value={activeTab} onValueChange={handleTabChange}>
          <TabsList className="mb-4">
            {providers.map((provider) => (
              <TabsTrigger key={provider.ProviderCode} value={provider.ProviderCode}>
                {provider.ProviderCode.charAt(0).toUpperCase() + provider.ProviderCode.slice(1)}
              </TabsTrigger>
            ))}
          </TabsList>

          {providers.map((provider) => (
            <TabsContent key={provider.ProviderCode} value={provider.ProviderCode} className="space-y-4">
              <div className="flex items-center justify-between">
                <Label htmlFor={`${provider.ProviderCode}-enabled`}>有効</Label>
                <Switch
                  id={`${provider.ProviderCode}-enabled`}
                  checked={provider.IsEnabled === 1}
                  onCheckedChange={() => handleToggleProvider(provider.ProviderCode)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${provider.ProviderCode}-client-id`}>クライアントID</Label>
                <Input
                  id={`${provider.ProviderCode}-client-id`}
                  name="clientId"
                  value={formData.clientId}
                  onChange={handleInputChange}
                  placeholder="クライアントIDを入力"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${provider.ProviderCode}-client-secret`}>クライアントシークレット</Label>
                <Input
                  id={`${provider.ProviderCode}-client-secret`}
                  name="clientSecret"
                  type="password"
                  value={formData.clientSecret}
                  onChange={handleInputChange}
                  placeholder="クライアントシークレットを入力"
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor={`${provider.ProviderCode}-callback-url`}>コールバックURL</Label>
                <Input
                  id={`${provider.ProviderCode}-callback-url`}
                  name="callbackUrl"
                  value={formData.callbackUrl}
                  onChange={handleInputChange}
                  placeholder="コールバックURLを入力"
                />
              </div>
            </TabsContent>
          ))}
        </Tabs>
      </CardContent>
      <CardFooter>
        <Button className="ml-auto" onClick={handleSave} disabled={saving}>
          {saving ? (
            <>
              <Loader2 className="mr-2 h-4 w-4 animate-spin" />
              保存中...
            </>
          ) : (
            "保存"
          )}
        </Button>
      </CardFooter>
    </Card>
  )
}
