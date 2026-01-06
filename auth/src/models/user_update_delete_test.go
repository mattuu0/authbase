package models_test

import (
	"auth/models"
	"testing"

	testtool "auth/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUpdateUser tests user information updates
func TestUpdateUser(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("update user name", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Creating user and updating name")
		user := testtool.CreateTestUser(t, db, "update@example.com", models.Google)

		// 名前を更新
		user.Name = "Updated Name"
		err := models.UpdateUser(user)
		require.NoError(t, err)

		// 更新が反映されたことを確認
		retrieved, result := models.GetUser(user.UserID)
		require.NoError(t, result.Error)
		assert.Equal(t, "Updated Name", retrieved.Name)
		testtool.LogSuccess(t, "User updated successfully")
	})
}

// TestDeleteUser tests user deletion
func TestDeleteUser(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("delete existing user", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Creating and deleting user")
		user := testtool.CreateTestUser(t, db, "delete@example.com", models.Google)

		err := models.DeleteUser(user.UserID)
		require.NoError(t, err)

		// ユーザーが削除されたことを確認
		_, result := models.GetUser(user.UserID)
		assert.Error(t, result.Error)
		assert.False(t, result.IsExists)
		testtool.LogSuccess(t, "User deleted successfully")
	})

	t.Run("delete non-existing user", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Attempting to delete non-existing user")
		err := models.DeleteUser("non-existing-user")
		assert.Error(t, err)
		testtool.LogSuccess(t, "Correctly failed to delete non-existing user")
	})
}
