// Package services provides unit tests for ParseAccessToken.
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

func TestParseAccessToken(t *testing.T) {
	testtool.SetupTestEnv(t)
	initJwt(testtool.TestJWTPrivateKey)

	t.Run("valid token returns correct fields", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "parse-user-001",
			Name:     "Parse User",
			Email:    "parse@example.com",
			Labels:   []string{"admin", "premium"},
			ProvCode: models.Google,
			ProvUid:  "google-parse-uid",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)
		require.NotNil(t, info)

		assert.Equal(t, claims.UserID, info.UserID)
		assert.Equal(t, claims.Name, info.Name)
		assert.Equal(t, claims.Email, info.Email)
		assert.Equal(t, claims.ProvCode, info.ProvCode)
		assert.Equal(t, claims.ProvUid, info.ProvUid)
		assert.ElementsMatch(t, claims.Labels, info.Labels)
		assert.Greater(t, info.Exp, time.Now().Unix())
	})

	t.Run("empty labels are preserved", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "parse-user-002",
			Name:     "No Labels",
			Email:    "nolabels@example.com",
			Labels:   []string{},
			ProvCode: models.Basic,
			ProvUid:  "",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)

		assert.Empty(t, info.Labels)
		assert.Equal(t, claims.Email, info.Email)
	})

	t.Run("nil labels result in empty slice", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "parse-user-003",
			Name:     "Nil Labels",
			Email:    "nil@example.com",
			Labels:   nil,
			ProvCode: models.Github,
			ProvUid:  "gh-uid",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)

		assert.Nil(t, info.Labels)
	})

	t.Run("invalid signature is rejected", func(t *testing.T) {
		invalidToken := "eyJhbGciOiJFZERTQSIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiJ4In0.invalidsig"

		_, err := ParseAccessToken(invalidToken)
		assert.Error(t, err)
	})

	t.Run("malformed token is rejected", func(t *testing.T) {
		_, err := ParseAccessToken("not.a.jwt")
		assert.Error(t, err)
	})

	t.Run("empty string is rejected", func(t *testing.T) {
		_, err := ParseAccessToken("")
		assert.Error(t, err)
	})

	t.Run("expired token is rejected", func(t *testing.T) {
		// 有効期限が過去のトークンを手動で生成
		token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
			"exp":      time.Now().Add(-1 * time.Minute).Unix(),
			"userID":   "expired-user",
			"name":     "Expired",
			"email":    "expired@example.com",
			"labels":   []string{},
			"provCode": "basic",
			"provUid":  "",
		})
		tokenString, err := token.SignedString(JwtPrivateKey)
		require.NoError(t, err)

		_, err = ParseAccessToken(tokenString)
		assert.Error(t, err)
	})

	t.Run("exp field is populated correctly", func(t *testing.T) {
		before := time.Now().Unix()

		claims := AccessTokenClaim{
			UserID:   "exp-user",
			Name:     "Exp User",
			Email:    "exp@example.com",
			Labels:   []string{},
			ProvCode: models.Discord,
			ProvUid:  "disc-uid",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)

		after := time.Now().Add(tokenExpiry).Unix()
		assert.GreaterOrEqual(t, info.Exp, before)
		assert.LessOrEqual(t, info.Exp, after)
	})
}

func TestAccessTokenClaimsIncludeNameAndEmail(t *testing.T) {
	testtool.SetupTestEnv(t)
	initJwt(testtool.TestJWTPrivateKey)

	t.Run("name and email are embedded in JWT claims", func(t *testing.T) {
		claims := AccessTokenClaim{
			UserID:   "claim-check-user",
			Name:     "Claim Checker",
			Email:    "checker@example.com",
			Labels:   []string{"verified"},
			ProvCode: models.Microsoft,
			ProvUid:  "ms-uid-claim",
		}

		tokenString, err := AccessTokenJwt(claims)
		require.NoError(t, err)

		// 公開鍵で直接パースしてクレームを検査
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))
		require.NoError(t, err)

		mapClaims, ok := token.Claims.(jwt.MapClaims)
		require.True(t, ok)

		assert.Equal(t, claims.Name, mapClaims["name"])
		assert.Equal(t, claims.Email, mapClaims["email"])
		assert.Equal(t, claims.UserID, mapClaims["userID"])
		assert.Equal(t, string(claims.ProvCode), mapClaims["provCode"])
	})
}
