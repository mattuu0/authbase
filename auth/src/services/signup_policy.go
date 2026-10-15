package services

import (
	"os"
	"strings"
)

// SignupAllowed はユーザーの新規登録(Basic)を許可するかを返す。
// 環境変数 ALLOW_SIGNUP が false / 0 / no / off のときだけ拒否し、未設定や空なら従来どおり許可する。
// Basic プロバイダの有効/無効はログインも止めてしまうため、登録だけを止めたいときはこちらを使う。
func SignupAllowed() bool {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("ALLOW_SIGNUP"))) {
	case "false", "0", "no", "off":
		return false
	}
	return true
}
