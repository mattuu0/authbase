// Package services provides unit tests for user management services.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)
	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("existing user is returned with correct fields", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "getuser@example.com", models.Basic)

		user, err := GetUser(m.UserID)
		require.NoError(t, err)

		assert.Equal(t, m.UserID, user.ID)
		assert.Equal(t, m.Name, user.Name)
		assert.Equal(t, m.Email, user.Email)
		assert.Equal(t, string(models.Basic), user.Provider)
		assert.False(t, user.Banned)
	})

	t.Run("labels are included in response", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "userlabels@example.com", models.Basic)
		testtool.CreateTestLabel(t, db, "beta", "#0000ff")
		require.NoError(t, m.AddLabel("beta"))

		user, err := GetUser(m.UserID)
		require.NoError(t, err)

		assert.Contains(t, user.Labels, "beta")
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		_, err := GetUser("ghost-user-id")
		assert.Error(t, err)
	})
}

func TestGetUsers(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)

	t.Run("returns all users", func(t *testing.T) {
		testtool.CreateTestUser(t, db, "allusers1@example.com", models.Basic)
		testtool.CreateTestUser(t, db, "allusers2@example.com", models.Basic)

		users, err := GetUsers()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(users), 2)
	})

	t.Run("each user has required fields", func(t *testing.T) {
		users, err := GetUsers()
		require.NoError(t, err)

		for _, u := range users {
			assert.NotEmpty(t, u.ID)
			assert.NotEmpty(t, u.Name)
			assert.NotEmpty(t, u.Email)
			assert.NotEmpty(t, u.CreatedAt)
		}
	})
}

func TestUpdateUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)

	t.Run("name is updated correctly", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "updatename@example.com", models.Basic)

		err := UpdateUser(UpdateUserData{
			ID:   m.UserID,
			Name: "Updated Name",
		})
		require.NoError(t, err)

		updated, err := GetUser(m.UserID)
		require.NoError(t, err)
		assert.Equal(t, "Updated Name", updated.Name)
	})

	t.Run("labels are replaced on update", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "updatelabels@example.com", models.Basic)
		testtool.CreateTestLabel(t, db, "old-label", "#aaaaaa")
		testtool.CreateTestLabel(t, db, "new-label", "#bbbbbb")
		require.NoError(t, m.AddLabel("old-label"))

		err := UpdateUser(UpdateUserData{
			ID:     m.UserID,
			Name:   m.Name,
			Labels: []string{"new-label"},
		})
		require.NoError(t, err)

		updated, err := GetUser(m.UserID)
		require.NoError(t, err)
		assert.Contains(t, updated.Labels, "new-label")
		assert.NotContains(t, updated.Labels, "old-label")
	})

	t.Run("updating non-existent user returns error", func(t *testing.T) {
		err := UpdateUser(UpdateUserData{
			ID:   "ghost-user",
			Name: "Ghost",
		})
		assert.Error(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)

	t.Run("existing user is deleted", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "deleteuser@example.com", models.Basic)

		err := DeleteUser(m.UserID)
		require.NoError(t, err)

		_, result := models.GetUser(m.UserID)
		assert.False(t, result.IsExists)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		err := DeleteUser("ghost-user-delete")
		assert.Error(t, err)
	})
}

func TestToggleBan(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)

	t.Run("user is banned correctly", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "banme@example.com", models.Basic)

		user, err := ToggleBan(BanArgs{IsBanned: true, UserID: m.UserID})
		require.NoError(t, err)
		assert.True(t, user.Banned)
	})

	t.Run("user is unbanned correctly", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "unbanme@example.com", models.Basic)
		m.IsBanned = 1
		models.UpdateUser(m)

		user, err := ToggleBan(BanArgs{IsBanned: false, UserID: m.UserID})
		require.NoError(t, err)
		assert.False(t, user.Banned)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		_, err := ToggleBan(BanArgs{IsBanned: true, UserID: "ghost-ban"})
		assert.Error(t, err)
	})
}

func TestGetInfo(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Basic)

	t.Run("returns public info", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "pubinfo@example.com", models.Basic)

		info, err := GetInfo(m.UserID)
		require.NoError(t, err)
		assert.Equal(t, m.UserID, info.UserID)
		assert.Equal(t, m.Name, info.Name)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		_, err := GetInfo("ghost-info")
		assert.Error(t, err)
	})
}

func TestGetMe(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	Init()

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("returns full user info including email and provider", func(t *testing.T) {
		m := testtool.CreateTestUser(t, db, "getme@example.com", models.Google)

		info, err := GetMe(m.UserID)
		require.NoError(t, err)
		assert.Equal(t, m.UserID, info.UserID)
		assert.Equal(t, m.Email, info.Email)
		assert.Equal(t, string(models.Google), info.ProvCode)
	})

	t.Run("non-existent user returns error", func(t *testing.T) {
		_, err := GetMe("ghost-me")
		assert.Error(t, err)
	})
}
