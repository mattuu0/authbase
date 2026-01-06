// Package models_test provides unit tests for user-related models.
// Tests cover user creation, retrieval, updates, deletion, and label management.
package models_test

import (
	"auth/models"
	"testing"

	testtool "auth/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateUser tests the creation of a new user
func TestCreateUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	// 既存のデータベース接続を差し替え
	models.ReplaceDB(db)

	// プロバイダーを作成
	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("valid user creation", func(t *testing.T) {
		user := &models.User{
			UserID:   "test-user-1",
			Name:     "Test User",
			Email:    "test@example.com",
			ProvCode: models.Google,
			ProvUID:  "google-uid-123",
		}

		err := models.CreateUser(user, models.Google)
		require.NoError(t, err)

		// ユーザーが作成されたことを確認
		retrieved, result := models.GetUser(user.UserID)
		require.NoError(t, result.Error)
		assert.True(t, result.IsExists)
		assert.Equal(t, user.Email, retrieved.Email)
		assert.Equal(t, user.Name, retrieved.Name)
	})

	t.Run("duplicate email", func(t *testing.T) {
		user1 := &models.User{
			UserID:   "test-user-2",
			Name:     "User 1",
			Email:    "duplicate@example.com",
			ProvCode: models.Google,
			ProvUID:  "uid-1",
		}

		user2 := &models.User{
			UserID:   "test-user-3",
			Name:     "User 2",
			Email:    "duplicate@example.com",
			ProvCode: models.Google,
			ProvUID:  "uid-2",
		}

		err := models.CreateUser(user1, models.Google)
		require.NoError(t, err)

		// 同じメールアドレスで2人目を作成しようとする
		err = models.CreateUser(user2, models.Google)
		assert.Error(t, err, "Should fail with duplicate email")
	})
}

// TestGetUser tests user retrieval by user ID
func TestGetUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("existing user", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "existing@example.com", models.Google)

		retrieved, result := models.GetUser(user.UserID)
		require.NoError(t, result.Error)
		assert.True(t, result.IsExists)
		assert.Equal(t, user.UserID, retrieved.UserID)
	})

	t.Run("non-existing user", func(t *testing.T) {
		_, result := models.GetUser("non-existing-user-id")
		assert.Error(t, result.Error)
		assert.False(t, result.IsExists)
	})
}

// TestGetUserByEmail tests user retrieval by email address
func TestGetUserByEmail(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("existing email", func(t *testing.T) {
		email := "findme@example.com"
		user := testtool.CreateTestUser(t, db, email, models.Google)

		retrieved, result := models.GetUserByEmail(email)
		require.NoError(t, result.Error)
		assert.True(t, result.IsExists)
		assert.Equal(t, user.Email, retrieved.Email)
	})

	t.Run("non-existing email", func(t *testing.T) {
		_, result := models.GetUserByEmail("notfound@example.com")
		assert.Error(t, result.Error)
		assert.False(t, result.IsExists)
	})
}

// TestUpdateUser tests user information updates
func TestUpdateUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("update user name", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "update@example.com", models.Google)

		// 名前を更新
		user.Name = "Updated Name"
		err := models.UpdateUser(user)
		require.NoError(t, err)

		// 更新が反映されたことを確認
		retrieved, result := models.GetUser(user.UserID)
		require.NoError(t, result.Error)
		assert.Equal(t, "Updated Name", retrieved.Name)
	})
}

// TestDeleteUser tests user deletion
func TestDeleteUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("delete existing user", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "delete@example.com", models.Google)

		err := models.DeleteUser(user.UserID)
		require.NoError(t, err)

		// ユーザーが削除されたことを確認
		_, result := models.GetUser(user.UserID)
		assert.Error(t, result.Error)
		assert.False(t, result.IsExists)
	})

	t.Run("delete non-existing user", func(t *testing.T) {
		err := models.DeleteUser("non-existing-user")
		assert.Error(t, err)
	})
}

// TestUserLabelManagement tests adding and removing labels from users
func TestUserLabelManagement(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("add label to user", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "label@example.com", models.Google)
		label := testtool.CreateTestLabel(t, db, "premium", "#FFD700")

		err := user.AddLabel(label.Name)
		require.NoError(t, err)

		// ラベルが追加されたことを確認
		labels, err := user.GetLabels()
		require.NoError(t, err)
		assert.Len(t, labels, 1)
		assert.Equal(t, "premium", labels[0].Name)
	})

	t.Run("remove label from user", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "removelabel@example.com", models.Google)
		label := testtool.CreateTestLabel(t, db, "temporary", "#FF0000")

		// ラベルを追加
		err := user.AddLabel(label.Name)
		require.NoError(t, err)

		// ラベルを削除
		err = user.RemoveLabel(label.Name)
		require.NoError(t, err)

		// ラベルが削除されたことを確認
		labels, err := user.GetLabels()
		require.NoError(t, err)
		assert.Len(t, labels, 0)
	})

	t.Run("get label names", func(t *testing.T) {
		user := testtool.CreateTestUser(t, db, "labelnames@example.com", models.Google)
		testtool.CreateTestLabel(t, db, "admin", "#0000FF")
		testtool.CreateTestLabel(t, db, "moderator", "#00FF00")

		user.AddLabel("admin")
		user.AddLabel("moderator")

		names, err := user.GetLabelNames()
		require.NoError(t, err)
		assert.Len(t, names, 2)
		assert.Contains(t, names, "admin")
		assert.Contains(t, names, "moderator")
	})
}

// TestSearchUser tests user search functionality
func TestSearchUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	// テストデータを作成
	user1 := testtool.CreateTestUser(t, db, "john.doe@example.com", models.Google)
	user1.Name = "John Doe"
	models.UpdateUser(user1)

	user2 := testtool.CreateTestUser(t, db, "jane.smith@example.com", models.Google)
	user2.Name = "Jane Smith"
	models.UpdateUser(user2)

	t.Run("search by name", func(t *testing.T) {
		users, err := models.SearchUserByName("John")
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
	})

	t.Run("search by email", func(t *testing.T) {
		users, err := models.SearchUserByEmail("jane")
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Contains(t, users[0].Email, "jane")
	})

	t.Run("search with no results", func(t *testing.T) {
		users, err := models.SearchUserByName("NonExistent")
		require.NoError(t, err)
		assert.Len(t, users, 0)
	})
}
