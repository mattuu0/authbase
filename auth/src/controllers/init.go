package controllers

import "os"

var (
	// アプリにリダイレクトするときのカスタムスキーマ
	CUSTOM_SCHEME = "authbase"

	// 認証完了後のリダイレクト先URL
	LOGIN_REDIRECT_URL = "/statics/home.html"
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
}
