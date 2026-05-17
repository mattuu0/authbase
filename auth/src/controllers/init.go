package controllers

import "os"

var (
	// アプリにリダイレクトするときのカスタムスキーマ
	CUSTOM_SCHEME = "authbase"

	// 認証完了後のリダイレクト先URL
	LOGIN_REDIRECT_URL = "/statics/home.html"

	// ログイン画面に表示するアプリ名
	APP_NAME = "AuthBase"
)

func init() {
	// 初期化関数
	// 環境変数からカスタムスキームを取得する
	CUSTOM_SCHEME_ENV := os.Getenv("CUSTOM_SCHEME")
	if CUSTOM_SCHEME_ENV != "" {
		CUSTOM_SCHEME = CUSTOM_SCHEME_ENV
	}

	// 環境変数からログイン後リダイレクト先を取得する
	LOGIN_REDIRECT_URL_ENV := os.Getenv("LOGIN_REDIRECT_URL")
	if LOGIN_REDIRECT_URL_ENV != "" {
		LOGIN_REDIRECT_URL = LOGIN_REDIRECT_URL_ENV
	}

	// 環境変数からアプリ名を取得する
	APP_NAME_ENV := os.Getenv("APP_NAME")
	if APP_NAME_ENV != "" {
		APP_NAME = APP_NAME_ENV
	}
}
