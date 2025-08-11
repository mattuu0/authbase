import secrets
import os

def generate_random_key(length=64):
    """
    暗号学的に安全なランダムキーを、指定された長さ（デフォルト64文字）で生成します。
    """
    return secrets.token_urlsafe(length)

def get_oauth_credentials(provider_name):
    """
    指定されたプロバイダーのOAuthクライアントIDとシークレットをユーザーに入力させます。
    設定をスキップするオプションも提供します。
    """
    print(f"\n--- {provider_name} OAuth 設定 ---")
    response = input(f"{provider_name} の設定をしますか？ (y/n): ")
    if response.lower() == 'y':
        client_id = input(f"{provider_name} のクライアントIDを入力してください: ")
        client_secret = input(f"{provider_name} のクライアントシークレットを入力してください: ")
        return client_id, client_secret
    else:
        # 'n'が入力された場合は空の文字列を返す
        return "", ""

def create_env_file(file_path, content):
    """
    指定されたファイルパスに、指定された内容で設定ファイルを生成します。
    """
    with open(file_path, "w", encoding="utf-8") as file:
        file.write(content.strip())
    print(f"✅ ファイル '{file_path}' を生成しました。")

def create_auth_env():
    """
    auth.env ファイルを生成するための設定情報を対話形式で取得し、ファイルに書き出します。
    """
    # 各OAuthプロバイダーの認証情報を対話形式で取得
    discord_client_id, discord_client_secret = get_oauth_credentials("Discord")
    google_client_id, google_client_secret = get_oauth_credentials("Google")
    github_client_id, github_client_secret = get_oauth_credentials("Github")
    microsoft_client_id, microsoft_client_secret = get_oauth_credentials("Microsoft")

    # 認証とセッション用のランダムキーを自動生成（長さ64文字）
    token_secret_key = generate_random_key()
    admin_session_key = generate_random_key()

    # auth.env のテンプレート
    auth_env_template = f"""
DiscordClientID = {discord_client_id}
DiscordClientSecret = {discord_client_secret}
DiscordCallback = https://localhost:8947/auth/oauth/discord/callback

GoogleClientID = {google_client_id}
GoogleClientSecret = {google_client_secret}
GoogleCallback = https://localhost:8947/auth/oauth/google/callback

GithubClientID = {github_client_id}
GithubClientSecret = {github_client_secret}
GithubCallback = https://localhost:8947/auth/oauth/github/callback

MicrosoftClientID = {microsoft_client_id}
MicrosoftClientSecret = {microsoft_client_secret}
MicrosoftCallback = https://localhost:8947/auth/oauth/microsoftonline/callback

DB_DSN = "main:main@tcp(db:3306)/authdb?charset=utf8mb4&parseTime=True&loc=Local"

TOKEN_SECRET = {token_secret_key}
ADMIN_SESSION_KEY = {admin_session_key}

GRPC_ADDR = ":9000"
"""
    create_env_file("auth.env", auth_env_template)

def main():
    # フォルダを移動する
    os.chdir("./data")

    """
    メイン処理：複数の設定ファイル生成関数を呼び出します。
    """
    print("--- OAuth およびアプリケーション設定の開始 ---")
    
    # auth.env ファイルを生成
    create_auth_env()

    # app.env のテンプレート
    session_secret_key = generate_random_key()
    app_env_template = f"""
SessionSecret = "{session_secret_key}"
GRPC_SERVER = auth:9000
DATABASE_DSN = "main:main@tcp(db:3306)/maindb?charset=utf8mb4&parseTime=True&loc=Local"
"""

    # app.env ファイルを生成
    create_env_file("app.env", app_env_template)

    print(f"\n--- 設定完了！ ---")
    print(f"設定ファイルがすべて生成されました。")

if __name__ == "__main__":
    main()