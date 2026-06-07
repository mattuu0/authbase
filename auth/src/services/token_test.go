// Package services provides unit tests for GetAccessToken.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccessToken(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)
	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("generates token for user without labels", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "notoken@example.com", models.Basic)

		tokenString, err := GetAccessToken(user.UserID)
		require.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)
		assert.Equal(t, user.UserID, info.UserID)
		assert.Equal(t, user.Name, info.Name)
		assert.Equal(t, user.Email, info.Email)
		assert.Empty(t, info.Labels)
	})

	t.Run("token contains labels assigned to user", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "labeled@example.com", models.Basic)

		// ラベルを作成してユーザーに付与
		label1 := testtool.CreateTestLabel(t, db, "premium", "#ff0000")
		label2 := testtool.CreateTestLabel(t, db, "verified", "#00ff00")
		_ = label1
		_ = label2
		require.NoError(t, user.AddLabel("premium"))
		require.NoError(t, user.AddLabel("verified"))

		tokenString, err := GetAccessToken(user.UserID)
		require.NoError(t, err)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{"premium", "verified"}, info.Labels)
	})

	t.Run("token contains correct provider info", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "oauthtoken@example.com", models.Google)

		tokenString, err := GetAccessToken(user.UserID)
		require.NoError(t, err)

		info, err := ParseAccessToken(tokenString)
		require.NoError(t, err)
		assert.Equal(t, models.Google, info.ProvCode)
		assert.Equal(t, user.ProvUID, info.ProvUid)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		_, err := GetAccessToken("non-existent-user-id")
		assert.Error(t, err)
	})

	t.Run("generated token is valid JWT", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "jwtvalid@example.com", models.Basic)

		tokenString, err := GetAccessToken(user.UserID)
		require.NoError(t, err)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return JwtPublicKey, nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodEdDSA.Alg()}))

		require.NoError(t, err)
		assert.True(t, token.Valid)
	})
}
