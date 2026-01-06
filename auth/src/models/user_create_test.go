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
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)

	testtool.LogStep(t, "Creating Test Provider (Google)")
	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("valid user creation", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Attempting to create a valid user")
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
		testtool.LogSuccess(t, "User created and retrieved successfully")
	})

	t.Run("duplicate email", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Attempting to create users with duplicate email")
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
		testtool.LogSuccess(t, "First user created")

		// 同じメールアドレスで2人目を作成しようとする
		err = models.CreateUser(user2, models.Google)
		assert.Error(t, err, "Should fail with duplicate email")
		testtool.LogSuccess(t, "Duplicate email check passed (creation failed as expected)")
	})
}
