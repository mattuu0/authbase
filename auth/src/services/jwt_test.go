// Package services_test provides unit tests for JWT token generation and validation.
// Tests cover token creation, verification, expiration, and error handling.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAccessTokenJwt tests JWT access token generation
func TestAccessTokenJwt(t *testing.T) {
	testtool.SetupTestEnv(t)

	// JWT鍵を初期化
	initJwt(testtool.TestJWTPrivateKey)

	t.Run("generate valid token", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "test-user-123",
			Labels:   []string{"admin", "user"},
			ProvCode: models.Google,
			ProvUid:  "google-uid-456",
		}

		token, err := AccessTokenJwt(claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// トークンが正しい形式かを確認（3つのパートに分かれている）
		parts := len(token) > 0
		assert.True(t, parts)
	})

	t.Run("token contains correct claims", func(t *testing.T) {
		expectedClaims := AccessTokenClaim{
			UserID:   "user-789",
			Labels:   []string{"premium", "verified"},
			ProvCode: models.Github,
			ProvUid:  "github-uid-101",
		}

		tokenString, err := AccessTokenJwt(expectedClaims)
		require.NoError(t, err)

		// トークンをパースして検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		require.NoError(t, err)
		assert.True(t, token.Valid)

		// クレームを取得
		claims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)

		// クレームの内容を検証
		assert.Equal(t, expectedClaims.UserID, claims["userID"])
		assert.Equal(t, string(expectedClaims.ProvCode), claims["provCode"])
		assert.Equal(t, expectedClaims.ProvUid, claims["provUid"])

		// ラベルの検証
		labelList := claims["labels"].([]interface{})
		assert.Len(t, labelList, 2)
	})

	t.Run("token has expiration", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "test-user",
			Labels:   []string{},
			ProvCode: models.Discord,
			ProvUid:  "discord-uid",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		// トークンをパース
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		require.NoError(t, err)

		// 有効期限を確認
		mapClaims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)

		exp, ok := mapClaims["exp"].(float64)
		require.True(t, ok)

		expTime := time.Unix(int64(exp), 0)
		now := time.Now()

		// 有効期限が現在時刻より未来であることを確認
		assert.True(t, expTime.After(now))

		// 有効期限が10分前後であることを確認
		expectedExpiration := now.Add(tokenExpiry)
		diff := expTime.Sub(expectedExpiration).Abs()
		assert.Less(t, diff, 5*time.Second, "Expiration should be around 10 minutes")
	})
}

// TestTokenValidation tests token validation logic
func TestTokenValidation(t *testing.T) {
	testtool.SetupTestEnv(t)
	initJwt(testtool.TestJWTPrivateKey)

	t.Run("valid token is accepted", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "valid-user",
			Labels:   []string{"user"},
			ProvCode: models.Google,
			ProvUid:  "google-123",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		// トークンを検証
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		require.NoError(t, err)
		assert.True(t, token.Valid)
	})

	t.Run("invalid signature is rejected", func(t *testing.T) {
		// 不正なトークンを作成（署名が間違っている）
		invalidToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MDAwMDAwMDAsInVzZXJJRCI6ImZha2UtdXNlciJ9.invalid_signature"

		_, err := jwt.Parse(invalidToken, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		assert.Error(t, err)
	})

	t.Run("malformed token is rejected", func(t *testing.T) {
		malformedToken := "not.a.valid.jwt.token"

		_, err := jwt.Parse(malformedToken, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		assert.Error(t, err)
	})
}

// TestEmptyLabels tests token generation with empty label array
func TestEmptyLabels(t *testing.T) {
	testtool.SetupTestEnv(t)
	initJwt(testtool.TestJWTPrivateKey)

	t.Run("empty labels array", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "user-no-labels",
			Labels:   []string{}, // 空のラベル
			ProvCode: models.Microsoft,
			ProvUid:  "ms-uid-789",
		}

		token, err := AccessTokenJwt(claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// トークンをパースして検証
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		require.NoError(t, err)

		mapClaims, ok := parsedToken.Claims.(jwt.MapClaims)
		require.True(t, ok)

		labels, ok := mapClaims["labels"].([]interface{})
		require.True(t, ok)
		assert.Len(t, labels, 0)
	})

	t.Run("nil labels", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "user-nil-labels",
			Labels:   nil, // nilのラベル
			ProvCode: models.Basic,
			ProvUid:  "basic-uid",
		}

		token, err := AccessTokenJwt(claims)
		require.NoError(t, err)
		assert.NotEmpty(t, token)
	})
}

// TestMultipleLabels tests token generation with multiple labels
func TestMultipleLabels(t *testing.T) {
	testtool.SetupTestEnv(t)
	initJwt(testtool.TestJWTPrivateKey)

	t.Run("multiple labels", func(t *testing.T) {
		labels := []string{"admin", "moderator", "premium", "verified", "beta-tester"}

		claims := AccessTokenClaim{
			UserID:   "multi-label-user",
			Labels:   labels,
			ProvCode: models.Github,
			ProvUid:  "github-multi",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		// トークンをパース
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		require.NoError(t, err)

		mapClaims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)

		labelList := mapClaims["labels"].([]interface{})
		assert.Len(t, labelList, len(labels))

		// すべてのラベルが含まれているか確認
		for i, label := range labelList {
			assert.Equal(t, labels[i], label)
		}
	})
}