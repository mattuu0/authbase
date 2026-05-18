// Package services provides unit tests for admin user management services.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAdminStatus(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("no system user returns HasSystemUser=false", func(t *testing.T) {
		// GetSystemAdminUser は存在しない場合 ErrRecordNotFound を返す。
		// GetAdminStatus はその error を呼び出し元に返すが HasSystemUser=false も同時に返す。
		status, _ := GetAdminStatus()
		assert.False(t, status.HasSystemUser)
	})

	t.Run("after creating admin returns HasSystemUser=true", func(t *testing.T) {
		err := CreateAdminUser(CreateAdminUserArgs{
			Username: "admin-status",
			Password: "AdminPass123!",
		})
		require.NoError(t, err)

		status, err := GetAdminStatus()
		require.NoError(t, err)
		assert.True(t, status.HasSystemUser)
	})
}

func TestCreateAdminUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("first admin is created successfully", func(t *testing.T) {
		err := CreateAdminUser(CreateAdminUserArgs{
			Username: "first-admin",
			Password: "SecurePass123!",
		})
		require.NoError(t, err)

		user, result := models.GetAdminUser("first-admin")
		require.NoError(t, result.Error)
		assert.Equal(t, "first-admin", user.Username)
		assert.Equal(t, 1, user.IsSystem)
		// パスワードが平文で保存されていない
		assert.NotEqual(t, "SecurePass123!", user.PasswordHash)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("second admin creation is rejected", func(t *testing.T) {
		// 1人目は既に作成済み（前のサブテストで作成）
		err := CreateAdminUser(CreateAdminUserArgs{
			Username: "second-admin",
			Password: "AnotherPass456!",
		})
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestLoginAdminUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	// テスト用管理者を作成
	err := CreateAdminUser(CreateAdminUserArgs{
		Username: "login-admin",
		Password: "LoginPass789!",
	})
	require.NoError(t, err)

	t.Run("correct credentials return user ID", func(t *testing.T) {
		userID, err := LoginAdminUser(LoginAdminUserArgs{
			Username: "login-admin",
			Password: "LoginPass789!",
		})
		require.NoError(t, err)
		assert.NotEmpty(t, userID)
	})

	t.Run("wrong password is rejected", func(t *testing.T) {
		_, err := LoginAdminUser(LoginAdminUserArgs{
			Username: "login-admin",
			Password: "WrongPass!",
		})
		assert.Error(t, err)
	})

	t.Run("non-existent username is rejected", func(t *testing.T) {
		_, err := LoginAdminUser(LoginAdminUserArgs{
			Username: "ghost-admin",
			Password: "SomePass123!",
		})
		assert.Error(t, err)
	})

	t.Run("empty credentials are rejected", func(t *testing.T) {
		_, err := LoginAdminUser(LoginAdminUserArgs{
			Username: "",
			Password: "",
		})
		assert.Error(t, err)
	})
}
