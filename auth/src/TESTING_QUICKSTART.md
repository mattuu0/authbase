# AuthBase テストガイド

## 前提条件

```bash
cd auth/src
go mod download
```

テストはすべてインメモリ SQLite で完結します。外部DBや環境変数の事前設定は不要です。

---

## テストの実行

### すべてのテストを一括実行

```bash
# プロジェクトルートから（推奨）
task test

# auth/src ディレクトリから直接実行
go test ./...
```

### サマリーのみ確認（高速）

```bash
task test-summary
```

出力例:
```
ok  	auth/controllers  0.4s
ok  	auth/integration  0.6s
ok  	auth/middlewares  0.8s
ok  	auth/models       1.3s
ok  	auth/services    13.1s
```

---

## レイヤー別実行

テストはレイヤーごとにパッケージが分かれています。

### モデル層（DB操作）

```bash
task test-models
```

| テストファイル | 内容 |
|---|---|
| `user_create_test.go` | ユーザー作成・重複メールアドレス検出 |
| `user_read_test.go` | ID/メールアドレスでのユーザー取得・検索 |
| `user_update_delete_test.go` | ユーザー更新・削除 |
| `user_label_test.go` | ラベルの追加・削除・取得 |

### サービス層（ビジネスロジック）

```bash
task test-services
```

| テストファイル | 内容 |
|---|---|
| `jwt_test.go` | アクセストークン生成・クレーム検証・有効期限 |
| `jwt_parse_test.go` | `ParseAccessToken` 正常系・異常系・期限切れ・name/email クレーム |
| `token_test.go` | `GetAccessToken` でのラベル・プロバイダ情報の埋め込み |
| `session_test.go` | セッション作成・検証・削除・BAN ユーザー拒否 |
| `basicuser_test.go` | Basic認証ユーザー作成・ログイン・パスワードハッシュ |
| `basicuser_signup_test.go` | サインアップ（重複・プロバイダ無効・正常系） |
| `basicuser_login_test.go` | ログイン（誤パスワード・存在しないユーザー・複数回試行） |
| `user_test.go` | ユーザーCRUD・BAN切替・公開情報取得 |
| `admin_test.go` | 管理者作成（2人目拒否）・ログイン・ステータス確認 |
| `label_test.go` | ラベルCRUD・重複名エラー |
| `bridge_test.go` | ブリッジトークン発行・交換・使い捨て・改ざん検知 |
| `auth_test.go` | ログアウトでセッション削除・他セッションへの影響なし |

### ミドルウェア層

```bash
task test-middlewares
```

| テストファイル | 内容 |
|---|---|
| `auth_test.go` | 有効トークン・無効トークン・BAN ユーザー・削除済みセッション・コンテキスト値 |

### コントローラー層

```bash
task test-controllers
```

| テストファイル | 内容 |
|---|---|
| `userinfo_test.go` | `/userinfo` エンドポイント・Bearer prefix 処理・セッショントークン拒否・exp フィールド |

### 統合テスト（E2E フロー）

```bash
task test-integration
```

| テストファイル | 内容 |
|---|---|
| `basic_auth_test.go` | サインアップ→/me→ログアウト→無効化→再ログインの完全フロー |
| `token_test.go` | サインアップ→セッショントークン→アクセストークン取得フロー |
| `session_test.go` | 複数セッション管理・片方ログアウト後の独立性確認 |
| `error_test.go` | 不正な認証情報・不正JSON・認証ヘッダーなし |
| `userinfo_test.go` | サインアップ→アクセストークン→`/userinfo` 完全フロー・ラベル検証・exp フィールド |

---

## 特定のテストだけ実行

```bash
cd auth/src

# テスト名を -run で指定（正規表現）
go test -v -run TestCreateBasicUser ./services/...
go test -v -run TestGetUserInfo ./controllers/...
go test -v -run TestUserInfoFlow ./integration/...
go test -v -run TestParseAccessToken ./services/...
go test -v -run TestToggleBan ./services/...

# サブテストまで指定
go test -v -run "TestGetAccessToken/token_contains_labels" ./services/...
```

---

## カバレッジ

```bash
# HTML レポートを生成して coverage.html を開く
task test-coverage
```

---

## その他のオプション

```bash
# レースコンディション検出
task test-race

# テストキャッシュをクリアして再実行
task test-clean
task test
```

---

## テスト設計の方針

| 項目 | 方針 |
|---|---|
| DB | インメモリ SQLite。外部接続不要 |
| 分離 | 各テスト関数で `SetupTestDB` / `CleanupTestDB` を呼び出し独立 |
| 認証キー | テスト用 Ed25519 固定鍵（`testing/constants.go`） |
| セッション秘密 | `TOKEN_SECRET=test-secret-key-for-testing-purposes-only` |

---

## トラブルシューティング

### テストが失敗する

```bash
task test-clean
task test
```

### 依存関係エラー

```bash
cd auth/src
go mod tidy
go mod download
```

### SQLite ドライバエラー

```bash
cd auth/src
go get gorm.io/driver/sqlite
go mod tidy
```
