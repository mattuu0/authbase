# AuthBase
認証方法のテンプレートをまとめたリポジトリ

## 使用技術
- 言語: Golang
- 環境: Docker
- DB: MySQL / PostgreSQL (環境変数で切替可能)
- ダッシュボード: React
- リバースプロキシ: nginx

## データベース設定
環境変数を使用して、使用するデータベースの種類と接続情報を設定できます。

### auth サービス
- `DB_TYPE`: `mysql` または `postgres` (デフォルト: `mysql`)
- `DB_DSN`: 接続文字列
  - MySQL例: `user:password@tcp(db:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local`
  - Postgres例: `host=db user=user password=password dbname=dbname port=5432 sslmode=disable TimeZone=Asia/Tokyo`

### app サービス
- `DATABASE_TYPE`: `mysql` または `postgres` (デフォルト: `mysql`)
- `DATABASE_DSN`: 接続文字列

## 認証後のリダイレクト設定
ログイン完了後のリダイレクト先URLを環境変数で変更できます。

### auth サービス
- `LOGIN_REDIRECT_URL`: ログイン完了後のリダイレクト先URL (デフォルト: `/statics/home.html`)
- `APP_NAME`: ログイン画面に表示するアプリ名 (デフォルト: `AuthBase`)

`config/auth.env` に以下を追加して設定します:
```env
LOGIN_REDIRECT_URL = "/your/path/here"
APP_NAME = "MyApp"
```

## セキュリティ関連の環境変数

auth サービスの正常動作に必要なシークレットです。`task setup`（`config/config.py`）を実行すると自動生成されます。手動で設定する場合は十分な長さのランダム文字列を使用してください。

| 変数名 | 説明 | 自動生成 |
|--------|------|---------|
| `TOKEN_SECRET` | セッショントークンの署名鍵 | ✅ |
| `ADMIN_SESSION_KEY` | 管理者セッションの署名鍵 | ✅ |
| `BRIDGE_TOKEN_SECRET` | ブリッジトークン（アプリ間トークン交換）の署名鍵 | ✅ |
| `JWT_PRIVATE_KEY` | アクセストークン署名用 Ed25519 秘密鍵（PEM形式） | ✅ (openssl) |

> **注意**: 既存環境でセットアップスクリプトを実行済みの場合、`BRIDGE_TOKEN_SECRET` を `auth.env` に手動で追加する必要があります。

## セットアップ方法
1. [taskfile](https://taskfile.dev/installation/) をインストールする
2. ```task --version``` を実行して taskfile を確認する
3. configs 配下にある *env_temolate をコピーして *.env を作成する
   各種 Secret などは [こちら](https://www.graviness.com/app/pwg/?l=64&n=1&m=1&r=3&s=1&c=0-9A-Za-z!%22%23%24%25%26'()*%2B%2C%5C-.%2F%3A%3B%3C%3D%3E%3F%40%5B%5C%5D%5E_%60%7B%7C%7D~) などで作成しておく
4. ```task setup``` を実行する 
    注意: すでに設置アップ済みの場合データベースの中身が削除されます
5. mysqlコンテナの起動が完了したら データベースコンテナを再起動する
    ```task restart```
6. https://localhost:8947/auth/_/ にアクセスして管理ユーザーを作成する
7. ダッシュボードで各種プロバイダの設定をする
8. https://localhost:8947/statics/ で確かめてみる
9.  終わり
   
## セットアップの動作
- alpine と openssl コンテナが起動します
  - nginx 用の自己証明書が発行されます
  - jwt 用の秘密鍵が発行されます
  - jwt 用の公開鍵が発行されます
- mysql コンテナが起動します
- auth コンテナが起動します
- app コンテナが起動します
- nginx コンテナが起動します

## ディレクトリ構成
- configs : 設定ファイル .env が格納されています
- database : 設定ファイル my.cnf が格納されています
- openssl : nginx 周りの jwt 秘密鍵や公開鍵が格納されています
- nginx : nginx 周りの設定ファイルが格納されています

## 各種コマンド
- ```task setup``` : セットアップ
- ```task clean``` : コンテナ落として全て削除
- ```task down``` : コンテナ落とす

## テスト

### 前提条件

```bash
cd auth/src
go mod download
```

テストはすべてインメモリ SQLite で完結します。外部DBや環境変数の事前設定は不要です。

### テストの実行

```bash
# すべてのテストを実行（推奨）
task test

# サマリーのみ表示（高速確認）
task test-summary
```

### レイヤー別に実行する

| コマンド | 対象 |
|---|---|
| `task test-models` | モデル層（DB操作） |
| `task test-services` | サービス層（ビジネスロジック） |
| `task test-middlewares` | ミドルウェア層（認証） |
| `task test-controllers` | コントローラー層 |
| `task test-integration` | 統合テスト（E2Eフロー） |

### カバレッジレポート

```bash
# HTML レポートを生成（coverage.html が作成される）
task test-coverage
```

### 特定のテスト関数だけ実行

```bash
cd auth/src
go test -v -run TestCreateBasicUser ./services/...
go test -v -run TestGetUserInfo ./controllers/...
go test -v -run TestUserInfoFlow ./integration/...
```

詳細は [`auth/src/TESTING_QUICKSTART.md`](auth/src/TESTING_QUICKSTART.md) を参照してください。

## CI/CD

GitHub Actions を使用して、コードの品質管理を行っています。

| ワークフロー | トリガー | 内容 |
|---|---|---|
| `auth-test.yml` | `auth/src/**` への push / PR | ビルド確認・全テスト実行・カバレッジ計測 |
| `docker-publish.yml` | `v*.*.*` タグ push | Docker イメージをビルドして GHCR に push |

### テストワークフローの流れ

1. Go セットアップ（`go.mod` からバージョンを自動取得）
2. `go build ./...` でビルド確認
3. `go test -v -count=1 -coverprofile=coverage.out ./...` で全レイヤーのテストを実行
4. カバレッジサマリーをログに表示
5. `coverage.out` / `coverage.html` をアーティファクトとして保存
