package models_test

import (
	"auth/models"
	"testing"

	testtool "auth/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUserLabelManagement tests adding and removing labels from users
func TestUserLabelManagement(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.CreateTestProvider(t, db, models.Google)

	t.Run("add label to user", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Adding 'premium' label to user")
		user := testtool.CreateTestUser(t, db, "label@example.com", models.Google)
		label := testtool.CreateTestLabel(t, db, "premium", "#FFD700")

		err := user.AddLabel(label.Name)
		require.NoError(t, err)

		// ラベルが追加されたことを確認
		labels, err := user.GetLabels()
		require.NoError(t, err)
		assert.Len(t, labels, 1)
		assert.Equal(t, "premium", labels[0].Name)
		testtool.LogSuccess(t, "Label added successfully")
	})

	t.Run("remove label from user", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Adding and removing 'temporary' label")
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
		testtool.LogSuccess(t, "Label removed successfully")
	})

	t.Run("get label names", func(t *testing.T) {
		testtool.RecordResult(t)
		testtool.LogStep(t, "Adding multiple labels and retrieving names")
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
		testtool.LogSuccess(t, "Label names retrieved successfully")
	})
}
