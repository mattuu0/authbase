package models_test

import (
	"auth/models"
	"testing"

	testtool "auth/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGetUser tests user retrieval by user ID
func TestGetUser(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("existing user", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Creating user and retrieving by ID")
		user := testtool.CreateTestUser(t, db, "existing@example.com", models.Google)

		retrieved, result := models.GetUser(user.UserID)
		require.NoError(t, result.Error)
		assert.True(t, result.IsExists)
		assert.Equal(t, user.UserID, retrieved.UserID)
		testtool.LogSuccess(t, "User retrieved successfully")
	})

	t.Run("non-existing user", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Attempting to retrieve non-existing user")
		_, result := models.GetUser("non-existing-user-id")
		assert.Error(t, result.Error)
		assert.False(t, result.IsExists)
		testtool.LogSuccess(t, "Correctly failed to retrieve non-existing user")
	})
}

// TestGetUserByEmail tests user retrieval by email address
func TestGetUserByEmail(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("existing email", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Creating user and retrieving by Email")
		email := "findme@example.com"
		user := testtool.CreateTestUser(t, db, email, models.Google)

		retrieved, result := models.GetUserByEmail(email)
		require.NoError(t, result.Error)
		assert.True(t, result.IsExists)
		assert.Equal(t, user.Email, retrieved.Email)
		testtool.LogSuccess(t, "User retrieved by email successfully")
	})

	t.Run("non-existing email", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Attempting to retrieve user by non-existing email")
		_, result := models.GetUserByEmail("notfound@example.com")
		assert.Error(t, result.Error)
		assert.False(t, result.IsExists)
		testtool.LogSuccess(t, "Correctly failed to retrieve user by non-existing email")
	})
}

// TestSearchUser tests user search functionality
func TestSearchUser(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	// テストデータを作成
	testtool.LogStep(t, "Setting up users for search test")
	user1 := testtool.CreateTestUser(t, db, "john.doe@example.com", models.Google)
	user1.Name = "John Doe"
	models.UpdateUser(user1)

	user2 := testtool.CreateTestUser(t, db, "jane.smith@example.com", models.Google)
	user2.Name = "Jane Smith"
	models.UpdateUser(user2)
	testtool.LogSuccess(t, "Test users created")

	t.Run("search by name", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Searching user by name 'John'")
		users, err := models.SearchUserByName("John")
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Equal(t, "John Doe", users[0].Name)
		testtool.LogSuccess(t, "Found expected user by name")
	})

	t.Run("search by email", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Searching user by email 'jane'")
		users, err := models.SearchUserByEmail("jane")
		require.NoError(t, err)
		assert.Len(t, users, 1)
		assert.Contains(t, users[0].Email, "jane")
		testtool.LogSuccess(t, "Found expected user by email")
	})

	t.Run("search with no results", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Searching with no results expected")
		users, err := models.SearchUserByName("NonExistent")
		require.NoError(t, err)
		assert.Len(t, users, 0)
		testtool.LogSuccess(t, "Correctly returned no results")
	})
}
