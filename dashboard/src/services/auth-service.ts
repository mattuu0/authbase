// 認証関連の操作を行うサービス
// 実際の実装ではバックエンドAPIとの通信を行う

export interface AuthUser {
  id: string
  email: string
  role: string
}

// 管理者が存在するかチェック
export async function checkAdminExists(): Promise<boolean> {
  // 実際の実装ではAPIを呼び出して管理者の存在を確認
  console.log("Checking if admin exists")

  // モックの実装（実際の実装では削除）
  // セッションストレージから管理者の存在を確認
  const adminExistsFlag = sessionStorage.getItem("admin_exists")

  // 初回アクセス時は管理者が存在しないと仮定
  if (adminExistsFlag === null) {
    return false
  }

  return adminExistsFlag === "true"
}

// ログイン処理
export async function login(email: string, password: string): Promise<AuthUser> {
  // 実際の実装ではAPIを呼び出して認証を行う
  console.log("Login attempt:", { email, password })

  // モックの認証処理（実際の実装では削除）
  if (email === "admin@example.com" && password === "password") {
    const user = {
      id: "usr_admin",
      email: "admin@example.com",
      role: "admin",
    }

    // 管理者が存在することをマーク
    sessionStorage.setItem("admin_exists", "true")

    // ユーザー情報をセッションに保存
    sessionStorage.setItem("auth_user", JSON.stringify(user))
    return user
  }

  // 認証失敗
  throw new Error("メールアドレスまたはパスワードが正しくありません。")
}

// サインアップ処理
export async function signup(email: string, password: string): Promise<AuthUser> {
  // 実際の実装ではAPIを呼び出してユーザー登録を行う
  console.log("Signup attempt:", { email, password })

  // モックのサインアップ処理（実際の実装では削除）
  // 既存ユーザーチェック（実際の実装ではAPIで行う）
  if (email === "admin@example.com") {
    throw new Error("このメールアドレスは既に登録されています。")
  }

  // 新規ユーザー作成
  const user = {
    id: `usr_${Math.random().toString(36).substring(2, 10)}`,
    email,
    role: "admin",
  }

  // 管理者が存在することをマーク
  sessionStorage.setItem("admin_exists", "true")

  // 成功したらユーザー情報を返す（実際の実装ではAPIからのレスポンスを返す）
  return user
}

// ログアウト処理
export async function logout(): Promise<void> {
  // 実際の実装ではAPIを呼び出してセッションを破棄
  console.log("Logout")

  // セッションからユーザー情報を削除
  sessionStorage.removeItem("auth_user")
}

// 現在のユーザーを取得
export async function getCurrentUser(): Promise<AuthUser | null> {
  // 実際の実装ではAPIを呼び出して現在のユーザー情報を取得

  // セッションからユーザー情報を取得
  const userJson = sessionStorage.getItem("auth_user")
  if (userJson) {
    return JSON.parse(userJson) as AuthUser
  }

  return null
}

// ユーザーが認証済みかチェック
export async function isAuthenticated(): Promise<boolean> {
  const user = await getCurrentUser()
  return user !== null
}
