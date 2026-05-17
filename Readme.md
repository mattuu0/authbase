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

`config/auth.env` に以下を追加して設定します:
```env
LOGIN_REDIRECT_URL = "/your/path/here"
```

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

## CI/CD
GitHub Actions を使用して、コードの品質管理を行っています。

- **自動テスト**: `auth/src` 配下のコードが変更（push または pull request）された際に自動で実行されます。
  - **実行内容**: ユニットテスト、統合テストの実行およびカバレッジの測定。
  - **カバレッジ**: 実行結果はアーティファクト（`coverage-report`）として保存されます。
- **詳細**: テストの実行方法や構成については `auth/src/TESTING_QUICKSTART.md` を参照してください。
