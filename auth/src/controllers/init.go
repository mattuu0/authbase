package controllers

import "os"

var (
	// アプリにリダイレクトするときのカスタムスキーマ
	CUSTOM_SCHEME = "authbase"
)

func init() {
	// 初期化関数
	// 環境変数からカスタムスキームを取得する
	CUSTOM_SCHEME_ENV := os.Getenv("CUSTOM_SCHEME")

	if CUSTOM_SCHEME_ENV != "" {
		// 環境変数からカスタムスキームを取得
		CUSTOM_SCHEME = CUSTOM_SCHEME_ENV
	}
}
