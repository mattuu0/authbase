# 🚀 AuthBase テスト クイックスタートガイド

このガイドではテストの実行方法を簡潔に説明します。

## 📋 前提条件

```bash
# 必要な依存関係をインストール
cd auth/src
go mod download
```

## 🎯 基本的な使い方

### 1. すべてのテストを実行

```bash
# Makefileを使用（推奨）
make test

# または直接実行
go test -v ./...
```

**出力例:**
```
╔════════════════════════════════════════════════╗
║                                                ║
║         AuthBase Testing Suite                ║
║                                                ║
╚════════════════════════════════════════════════╝

🚀 [START] CreateUser - Success Case
📍 [STEP] CreateUser - Success Case: Creating test provider
📍 [STEP] CreateUser - Success Case: Creating new user
✅ [PASS] CreateUser - Success Case: User creation: no error
📍 [STEP] CreateUser - Success Case: Verifying user was created
✅ [PASS] CreateUser - Success Case: User retrieval: no error
🏁 [END] CreateUser - Success Case: completed in 45ms
─────────────────────────────────────────────────────
```

### 2. 特定のテストカテゴリを実行

```bash
# モデル層のみ
make test-models

# サービス層のみ
make test-services

# 統合テストのみ
make test-integration
```

### 3. 機能別テスト

```bash
# ユーザー機能のテスト
make test-user

# 認証機能のテスト
make test-auth

# セッション機能のテスト
make test-session
```

## 🎨 テスト出力の見方

### ✅ 成功時の出力

```
🚀 [START] TestName
📍 [STEP] TestName: Doing something
✅ [PASS] TestName: Operation succeeded
🏁 [END] TestName: completed in 23ms
```

### ❌ 失敗時の出力

```
🚀 [START] TestName
📍 [STEP] TestName: Doing something
❌ [FAIL] TestName: Expected X but got Y
```

### ℹ️ 情報出力

```
ℹ️ [INFO] TestName: Additional information
⚠️ [WARN] TestName: Warning message
```

## 🔧 便利なコマンド

### サマリーのみ表示（高速）

```bash
make test-summary
```

**出力:**
```
✓ ok    auth/models      2.345s
✓ ok    auth/services    3.123s
✓ ok    auth/middlewares 1.234s
```

### カバレッジレポート生成

```bash
make test-coverage
```

**出力:**
```
▶ Generating coverage report...
✓ Coverage report generated: coverage.html
Total Coverage: 85.4%
```

ブラウザで `coverage.html` を開くとビジュアルなカバレッジを確認できます。

### テストをリアルタイムで監視

```bash
# シェルスクリプト版
chmod +x test.sh
./test.sh watch

# またはMakefile版（entr コマンドが必要）
make watch
```

## 📊 テストの構成

```
auth/src/
├── models/
│   ├── user_create_test.go      # ユーザー作成
│   ├── user_retrieve_test.go    # ユーザー取得
│   ├── user_update_test.go      # ユーザー更新
│   ├── user_delete_test.go      # ユーザー削除
│   └── user_label_test.go       # ラベル管理
├── services/
│   ├── basicuser_signup_test.go # サインアップ
│   └── basicuser_login_test.go  # ログイン
└── testing/
    ├── setup.go                  # テストセットアップ
    └── helper.go                 # テストヘルパー
```

## 🎭 シェルスクリプトを使用

```bash
# すべてのテスト（見やすい出力）
./test.sh

# モデルのみ
./test.sh models

# サービスのみ
./test.sh services

# クイックテスト
./test.sh quick
```

## 💡 Tips

### 特定のテスト関数だけ実行

```bash
go test -v -run TestCreateUser_Success ./models/
```

### 失敗したテストのみ再実行

```bash
go test -v -count=1 ./... | grep FAIL
```

### テストを並列実行

```bash
go test -v -parallel 4 ./...
```

### レースコンディションをチェック

```bash
make test-race
```

## 🐛 トラブルシューティング

### テストが失敗する場合

1. **キャッシュをクリア**
   ```bash
   make test-clean
   make test
   ```

2. **依存関係を再インストール**
   ```bash
   go mod tidy
   go mod download
   ```

3. **詳細ログを表示**
   ```bash
   make test-verbose
   ```

### よくあるエラー

#### "failed to connect to test database"

```bash
# SQLiteドライバを確認
go get gorm.io/driver/sqlite
go mod tidy
```

#### "package testing is not in GOROOT"

```bash
# auth/src ディレクトリで実行していることを確認
cd auth/src
make test
```

## 📈 継続的インテグレーション

GitHub Actionsでの使用例:

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.23'
      - name: Run tests
        run: |
          cd auth/src
          make test-coverage
```

## 🎓 次のステップ

- 📖 詳細なドキュメント: `TEST_README.md` を参照
- 🔍 テストコードを読む: `models/*_test.go` から始める
- ✏️ テストを追加: 既存のテストをテンプレートとして使用

## 📞 サポート

問題が発生した場合:
1. `make help` でコマンド一覧を確認
2. `TEST_README.md` で詳細を確認
3. GitHubでIssueを作成

---

**Happy Testing! 🎉**