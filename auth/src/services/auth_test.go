// Package services provides unit tests for Logout service.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogout(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("logout deletes the session", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "logout@example.com", models.Google)

		token, err := NewSession(SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "10.0.0.1",
			UserAgent: "Test Browser",
		})
		require.NoError(t, err)

		session, err := GetSession(token)
		require.NoError(t, err)

		err = Logout(session)
		require.NoError(t, err)

		// セッションが削除されていることを確認
		_, err = GetSession(token)
		assert.Error(t, err)
	})

	t.Run("other sessions for same user remain after logout", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "multilogout@example.com", models.Google)

		token1, err := NewSession(SessionArgs{UserID: user.UserID, RemoteIP: "10.0.0.2", UserAgent: "Browser A"})
		require.NoError(t, err)
		token2, err := NewSession(SessionArgs{UserID: user.UserID, RemoteIP: "10.0.0.3", UserAgent: "Browser B"})
		require.NoError(t, err)

		session1, err := GetSession(token1)
		require.NoError(t, err)

		// session1 だけログアウト
		err = Logout(session1)
		require.NoError(t, err)

		// session2 は引き続き有効
		session2, err := GetSession(token2)
		require.NoError(t, err)
		assert.Equal(t, user.UserID, session2.UserID)
	})
}
