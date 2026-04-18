// Package services はビジネスロジックを管理します。
package services

import (
	"auth/models"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// BridgeTokenClaims はブリッジトークン用JWTに含めるデータ構造（クレーム）です
type BridgeTokenClaims struct {
	UserID  string `json:"user_id"`  // ユーザー識別子
	TokenID string `json:"token_id"` // トークン一意識別子
	jwt.RegisteredClaims
}

// bridgeTokenSecret はブリッジトークン署名用の秘密鍵です。
// ※本来は環境変数から取得して管理する必要があります。
var bridgeTokenSecret = []byte("super-secret-bridge-key")

// IssueBridgeToken は指定されたユーザーIDに対して、5分間有効な一時的JWTブリッジトークンを発行します。
func IssueBridgeToken(userID string) (string, error) {
	// トークンを一意に識別するためのUUIDを生成
	tokenID := uuid.New().String()
	expiresAt := time.Now().Add(5 * time.Minute)

	// JWTのクレームを構築
	claims := BridgeTokenClaims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// JWTを署名して生成
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(bridgeTokenSecret)
	if err != nil {
		return "", err
	}

	// 使い捨て管理のため、トークンID（UUID）をDBに保存
	bridgeToken := models.BridgeToken{
		Token:     tokenID,
		UserID:    userID,
		ExpiresAt: expiresAt,
		IsUsed:    false,
	}

	// DBへの保存処理
	if err := models.GetDB().Create(&bridgeToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExchangeBridgeToken はJWT形式のブリッジトークンを検証し、未使用であれば本番用アクセストークンと交換します。
func ExchangeBridgeToken(tokenString string) (string, error) {
	// 1. JWTの署名と形式を検証
	token, err := jwt.ParseWithClaims(tokenString, &BridgeTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return bridgeTokenSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("無効なブリッジトークンです")
	}

	claims, ok := token.Claims.(*BridgeTokenClaims)
	if !ok {
		return "", errors.New("トークンクレームの解析に失敗しました")
	}

	// 2. DB上でトークンが未使用かつ存在することを確認（使い捨ての検証）
	var bridgeToken models.BridgeToken
	err = models.GetDB().Where("token = ? AND is_used = ?", claims.TokenID, false).First(&bridgeToken).Error
	if err != nil {
		return "", errors.New("トークンは既に使用済みか存在しません")
	}

	// 3. トークンを即座に使用済みに更新（二重利用防止）
	bridgeToken.IsUsed = true
	if err := models.GetDB().Save(&bridgeToken).Error; err != nil {
		return "", fmt.Errorf("トークンの更新に失敗しました: %v", err)
	}

	// 4. 検証成功につき、本番用のアクセストークンを生成して返却
	accessToken, err := GetAccessToken(claims.UserID)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
