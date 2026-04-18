package services

import (
	"auth/models"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"
)

// IssueBridgeToken は指定されたユーザーIDに対して一時的なブリッジトークンを発行します
func IssueBridgeToken(userID string) (string, error) {
	// ランダムなトークンを生成
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := hex.EncodeToString(b)

	// ブリッジトークンを作成しDBに保存
	bridgeToken := models.BridgeToken{
		Token:     token,
		UserID:    userID,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		IsUsed:    false,
	}

	if err := models.GetDB().Create(&bridgeToken).Error; err != nil {
		return "", err
	}

	return token, nil
}

// ExchangeBridgeToken はブリッジトークンを検証し、本番用のアクセストークンと交換します
func ExchangeBridgeToken(token string) (string, error) {
	var bridgeToken models.BridgeToken

	// 未使用かつ有効期限内のトークンを検索
	err := models.GetDB().Where("token = ? AND is_used = ? AND expires_at > ?", token, false, time.Now()).First(&bridgeToken).Error
	if err != nil {
		return "", errors.New("トークンが無効または有効期限切れです")
	}

	// トークンを使用済みに更新（使い捨て）
	bridgeToken.IsUsed = true
	models.GetDB().Save(&bridgeToken)

	// 本番用のアクセストークンを生成
	accessToken, err := GetAccessToken(bridgeToken.UserID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
