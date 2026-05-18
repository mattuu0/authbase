// Package services provides unit tests for bridge token services.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupBridgeTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db := testtool.SetupTestDB(t)
	// BridgeToken テーブルを追加でマイグレーション
	require.NoError(t, db.AutoMigrate(&models.BridgeToken{}))
	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()
	return db
}

func TestIssueBridgeToken(t *testing.T) {
	db := setupBridgeTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	t.Run("issues a valid JWT string", func(t *testing.T) {
		tokenString, err := IssueBridgeToken("user-bridge-001", "dummy-refresh-token")
		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// JWT形式（3パーツ）であることを確認
		parts := 0
		for _, c := range tokenString {
			if c == '.' {
				parts++
			}
		}
		assert.Equal(t, 2, parts)
	})

	t.Run("bridge token is persisted in DB", func(t *testing.T) {
		_, err := IssueBridgeToken("user-bridge-002", "refresh-for-db")
		require.NoError(t, err)

		var count int64
		db.Model(&models.BridgeToken{}).Where("user_id = ?", "user-bridge-002").Count(&count)
		assert.Equal(t, int64(1), count)
	})

	t.Run("each call produces a different token", func(t *testing.T) {
		token1, err1 := IssueBridgeToken("user-bridge-003", "refresh-1")
		token2, err2 := IssueBridgeToken("user-bridge-003", "refresh-2")
		require.NoError(t, err1)
		require.NoError(t, err2)
		assert.NotEqual(t, token1, token2)
	})
}

func TestExchangeBridgeToken(t *testing.T) {
	db := setupBridgeTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	t.Run("valid bridge token returns refresh token", func(t *testing.T) {
		refreshToken := "original-refresh-token"
		bridgeToken, err := IssueBridgeToken("user-exchange-001", refreshToken)
		require.NoError(t, err)

		result, err := ExchangeBridgeToken(bridgeToken)
		require.NoError(t, err)
		assert.Equal(t, refreshToken, result["refresh_token"])
	})

	t.Run("bridge token is deleted after exchange (one-time use)", func(t *testing.T) {
		bridgeToken, err := IssueBridgeToken("user-exchange-002", "one-time-refresh")
		require.NoError(t, err)

		_, err = ExchangeBridgeToken(bridgeToken)
		require.NoError(t, err)

		// 2回目の交換は失敗する
		_, err = ExchangeBridgeToken(bridgeToken)
		assert.Error(t, err)
	})

	t.Run("invalid JWT string is rejected", func(t *testing.T) {
		_, err := ExchangeBridgeToken("not.a.valid.bridge.token")
		assert.Error(t, err)
	})

	t.Run("empty string is rejected", func(t *testing.T) {
		_, err := ExchangeBridgeToken("")
		assert.Error(t, err)
	})

	t.Run("tampered token is rejected", func(t *testing.T) {
		_, err := ExchangeBridgeToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiZXZpbCJ9.tampered")
		assert.Error(t, err)
	})
}
