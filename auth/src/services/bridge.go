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
func IssueBridgeToken(userID string, refreshToken string) (string, error) {
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
	// 使い捨て管理のため、トークンID（UUID）とリフレッシュトークンをDBに保存
	bridgeToken := models.BridgeToken{
		Token:        tokenID,
		UserID:       userID,
		ExpiresAt:    expiresAt,
		RefreshToken: refreshToken, // リフレッシュトークンを保存
	}

	// DBへの保存処理
	if err := models.GetDB().Create(&bridgeToken).Error; err != nil {
		return "", err
	}

	return tokenString, nil
}

// ExchangeBridgeToken はJWT形式のブリッジトークンを検証し、未使用であればアクセストークンと保存されていたリフレッシュトークンを返却します。
func ExchangeBridgeToken(tokenString string) (map[string]string, error) {
	// 1. JWTの署名と形式を検証
	token, err := jwt.ParseWithClaims(tokenString, &BridgeTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return bridgeTokenSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("無効なブリッジトークンです")
	}

	claims, ok := token.Claims.(*BridgeTokenClaims)
	if !ok {
		return nil, errors.New("トークンクレームの解析に失敗しました")
	}

	// 2. DB上でトークンが存在することを確認（存在すれば有効）
	var bridgeToken models.BridgeToken
	err = models.GetDB().Where("token = ?", claims.TokenID).First(&bridgeToken).Error
	if err != nil {
		return nil, errors.New("トークンは無効か、既に使用済みです")
	}

	// 3. トークンをDBから物理削除（使い捨て・存在自体を無効化）
	if err := models.GetDB().Delete(&bridgeToken).Error; err != nil {
		return nil, fmt.Errorf("トークンの削除に失敗しました: %v", err)
	}

	// 4. 新しいアクセストークンを生成
	accessToken, err := GetAccessToken(claims.UserID)
	if err != nil {
		return nil, err
	}

	// 5. 保存されていたリフレッシュトークンと新しいアクセストークンを返却
	return map[string]string{
		"id_token":      accessToken,
		"refresh_token": bridgeToken.RefreshToken,
	}, nil
}
