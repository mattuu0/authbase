// Package services provides unit tests for Logout service.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"有効: 英数字8文字", "abcd1234", false},
		{"有効: 英数字混在で長い", "password123", false},
		{"有効: 記号を含んでいても通過", "pass1!@#", false},
		{"無効: 7文字", "abcd123", true},
		{"無効: 数字のみ", "12345678", true},
		{"無効: 英字のみ", "abcdefgh", true},
		{"無効: 空文字", "", true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePassword(tc.input)
			if tc.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCreateBasicUser_PasswordValidation(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)

	base := CreateBasicUserArgs{Name: "テストユーザー", Email: "pwtest@example.com"}

	t.Run("短すぎるパスワードは拒否される", func(t *testing.T) {
		_, result := CreateBasicUser(CreateBasicUserArgs{Name: base.Name, Email: "short@example.com", Password: "abc1"})
		assert.Error(t, result.Error)
		assert.Equal(t, 400, result.Code)
	})

	t.Run("数字を含まないパスワードは拒否される", func(t *testing.T) {
		_, result := CreateBasicUser(CreateBasicUserArgs{Name: base.Name, Email: "nodigit@example.com", Password: "abcdefgh"})
		assert.Error(t, result.Error)
		assert.Equal(t, 400, result.Code)
	})

	t.Run("英字を含まないパスワードは拒否される", func(t *testing.T) {
		_, result := CreateBasicUser(CreateBasicUserArgs{Name: base.Name, Email: "noletter@example.com", Password: "12345678"})
		assert.Error(t, result.Error)
		assert.Equal(t, 400, result.Code)
	})

	t.Run("要件を満たすパスワードは登録される", func(t *testing.T) {
		token, result := CreateBasicUser(CreateBasicUserArgs{Name: base.Name, Email: "valid@example.com", Password: "pass1234"})
		require.NoError(t, result.Error)
		assert.NotEmpty(t, token)
	})
}

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
