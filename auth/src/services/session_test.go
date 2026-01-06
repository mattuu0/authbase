// Package services_test provides unit tests for session management services.
// Tests cover session creation, token generation, validation, and deletion.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewSession tests session creation and token generation
func TestNewSession(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	// DBを置き換え
	models.ReplaceDB(db)

	testtool.SetupTestEnv(t)

	Init()

	// テストユーザーとプロバイダーを作成
	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "session@example.com", models.Google)

	t.Run("create new session", func(t *testing.T) {
		args := SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "192.168.1.1",
			UserAgent: "Mozilla/5.0 Test Browser",
		}

		token, err := NewSession(args)
		require.NoError(t, err)
		assert.NotEmpty(t, token)

		// セッションがデータベースに保存されたことを確認
		sessions, err := models.GetAllSessions()
		require.NoError(t, err)
		assert.Len(t, sessions, 1)
		assert.Equal(t, user.UserID, sessions[0].UserID)
		assert.Equal(t, args.RemoteIP, sessions[0].RemoteIP)
		assert.Equal(t, args.UserAgent, sessions[0].UserAgent)
	})

	t.Run("multiple sessions for same user", func(t *testing.T) {
		args1 := SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "192.168.1.2",
			UserAgent: "Browser 1",
		}
		args2 := SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "192.168.1.3",
			UserAgent: "Browser 2",
		}

		token1, err := NewSession(args1)
		require.NoError(t, err)
		assert.NotEmpty(t, token1)

		token2, err := NewSession(args2)
		require.NoError(t, err)
		assert.NotEmpty(t, token2)

		// トークンが異なることを確認
		assert.NotEqual(t, token1, token2)
	})

	t.Run("banned user cannot create session", func(t *testing.T) {
		// BANされたユーザーを作成
		bannedUser := testtool.CreateTestUser(t, db, "banned@example.com", models.Google)
		bannedUser.IsBanned = 1
		models.UpdateUser(bannedUser)

		args := SessionArgs{
			UserID:    bannedUser.UserID,
			RemoteIP:  "192.168.1.4",
			UserAgent: "Test Browser",
		}

		_, err := NewSession(args)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "banned")
	})
}

// TestSessionTokenValidation tests session token validation
func TestSessionTokenValidation(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	Init()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "token@example.com", models.Google)

	t.Run("validate valid token", func(t *testing.T) {
		args := SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "192.168.1.5",
			UserAgent: "Valid Browser",
		}

		token, err := NewSession(args)
		require.NoError(t, err)

		// トークンを検証
		sessionID, err := ValidateSessionToken(token)
		require.NoError(t, err)
		assert.NotEmpty(t, sessionID)

		// セッションを取得
		session, err := models.GetSession(sessionID)
		require.NoError(t, err)
		assert.Equal(t, user.UserID, session.UserID)
	})

	t.Run("invalid token", func(t *testing.T) {
		invalidToken := "invalid.token.string"

		_, err := ValidateSessionToken(invalidToken)
		assert.Error(t, err)
	})

	t.Run("malformed token", func(t *testing.T) {
		malformedToken := "not-a-jwt-token"

		_, err := ValidateSessionToken(malformedToken)
		assert.Error(t, err)
	})
}

// TestGetSession tests session retrieval by token
func TestGetSession(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	Init()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "getsession@example.com", models.Google)

	t.Run("get existing session", func(t *testing.T) {
		args := SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "192.168.1.6",
			UserAgent: "Get Browser",
		}

		token, err := NewSession(args)
		require.NoError(t, err)

		// セッションを取得
		session, err := GetSession(token)
		require.NoError(t, err)
		assert.NotNil(t, session)
		assert.Equal(t, user.UserID, session.UserID)
		assert.Equal(t, args.RemoteIP, session.RemoteIP)
	})

	t.Run("get non-existent session", func(t *testing.T) {
		// 存在しないセッションIDでトークンを生成
		fakeToken, _ := GenSessionToken("non-existent-session-id")

		_, err := GetSession(fakeToken)
		assert.Error(t, err)
	})
}

// TestDeleteSession tests session deletion
func TestDeleteSession(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	Init()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "deletesession@example.com", models.Google)

	t.Run("delete existing session", func(t *testing.T) {
		args := SessionArgs{
			UserID:    user.UserID,
			RemoteIP:  "192.168.1.7",
			UserAgent: "Delete Browser",
		}

		token, err := NewSession(args)
		require.NoError(t, err)

		// セッションIDを取得
		sessionID, err := ValidateSessionToken(token)
		require.NoError(t, err)

		// セッションを削除
		err = DeleteSession(sessionID)
		require.NoError(t, err)

		// セッションが削除されたことを確認
		_, err = models.GetSession(sessionID)
		assert.Error(t, err)
	})

	t.Run("delete non-existent session", func(t *testing.T) {
		err := DeleteSession("non-existent-session-id")
		// エラーにならない（削除済みまたは存在しない）
		assert.NoError(t, err)
	})
}

// TestGetAllSessions tests retrieving all sessions
func TestGetAllSessions(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	Init()

	testtool.CreateTestProvider(t, db, models.Google)
	user1 := testtool.CreateTestUser(t, db, "user1@example.com", models.Google)
	user2 := testtool.CreateTestUser(t, db, "user2@example.com", models.Google)

	t.Run("get all sessions", func(t *testing.T) {
		// 複数のセッションを作成
		NewSession(SessionArgs{UserID: user1.UserID, RemoteIP: "192.168.1.8", UserAgent: "Browser 1"})
		NewSession(SessionArgs{UserID: user1.UserID, RemoteIP: "192.168.1.9", UserAgent: "Browser 2"})
		NewSession(SessionArgs{UserID: user2.UserID, RemoteIP: "192.168.1.10", UserAgent: "Browser 3"})

		sessions, err := GetAllSessions()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(sessions), 3)

		// セッションの構造を確認
		for _, session := range sessions {
			assert.NotEmpty(t, session.ID)
			assert.NotEmpty(t, session.UserID)
			assert.NotEmpty(t, session.IPAddress)
			assert.NotEmpty(t, session.UserAgent)
			assert.Greater(t, session.CreatedAt, int64(0))
			assert.Greater(t, session.ExpiresAt, int64(0))
			assert.True(t, session.IsActive)
		}
	})

	t.Run("empty session list", func(t *testing.T) {
		// すべてのセッションを削除
		allSessions, _ := models.GetAllSessions()
		for _, s := range allSessions {
			DeleteSession(s.SessionID)
		}

		sessions, err := GetAllSessions()
		require.NoError(t, err)
		assert.Len(t, sessions, 0)
	})
}