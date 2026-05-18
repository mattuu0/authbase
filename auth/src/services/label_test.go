// Package services provides unit tests for label management services.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateLabel(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("creates label with name and color", func(t *testing.T) {
		err := CreateLabel(CreateLabelArgs{Name: "vip", Color: "#gold"})
		require.NoError(t, err)

		label, err := GetLabel("vip")
		require.NoError(t, err)
		assert.Equal(t, "vip", label.Name)
		assert.Equal(t, "#gold", label.Color)
	})

	t.Run("duplicate name returns error", func(t *testing.T) {
		err := CreateLabel(CreateLabelArgs{Name: "dup", Color: "#111"})
		require.NoError(t, err)

		err = CreateLabel(CreateLabelArgs{Name: "dup", Color: "#222"})
		assert.Error(t, err)
	})
}

func TestGetLabels(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("returns all created labels", func(t *testing.T) {
		require.NoError(t, CreateLabel(CreateLabelArgs{Name: "tag-a", Color: "#aaa"}))
		require.NoError(t, CreateLabel(CreateLabelArgs{Name: "tag-b", Color: "#bbb"}))

		labels, err := GetLabels()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(labels), 2)

		names := make([]string, len(labels))
		for i, l := range labels {
			names[i] = l.Name
		}
		assert.Contains(t, names, "tag-a")
		assert.Contains(t, names, "tag-b")
	})

	t.Run("each label has required fields", func(t *testing.T) {
		labels, err := GetLabels()
		require.NoError(t, err)

		for _, l := range labels {
			assert.NotEmpty(t, l.Name)
			assert.NotEmpty(t, l.ID)
		}
	})
}

func TestGetLabel(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("existing label is returned", func(t *testing.T) {
		require.NoError(t, CreateLabel(CreateLabelArgs{Name: "findme", Color: "#fff"}))

		label, err := GetLabel("findme")
		require.NoError(t, err)
		assert.Equal(t, "findme", label.Name)
	})

	t.Run("non-existent label returns error", func(t *testing.T) {
		_, err := GetLabel("ghost-label")
		assert.Error(t, err)
	})
}

func TestUpdateLabel(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("name and color are updated", func(t *testing.T) {
		require.NoError(t, CreateLabel(CreateLabelArgs{Name: "old-name", Color: "#000"}))

		err := UpdateLabel(LabelUpdateArgs{ID: "old-name", Name: "old-name", Color: "#ffffff"})
		require.NoError(t, err)

		label, err := GetLabel("old-name")
		require.NoError(t, err)
		assert.Equal(t, "#ffffff", label.Color)
	})

	t.Run("non-existent label returns error", func(t *testing.T) {
		err := UpdateLabel(LabelUpdateArgs{ID: "ghost", Name: "ghost", Color: "#fff"})
		assert.Error(t, err)
	})
}

func TestDeleteLabel(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)

	t.Run("existing label is deleted", func(t *testing.T) {
		require.NoError(t, CreateLabel(CreateLabelArgs{Name: "remove-me", Color: "#eee"}))

		err := DeleteLabel(DeleteLabelArgs{ID: "remove-me"})
		require.NoError(t, err)

		_, err = GetLabel("remove-me")
		assert.Error(t, err)
	})

	t.Run("non-existent label returns error", func(t *testing.T) {
		err := DeleteLabel(DeleteLabelArgs{ID: "ghost-delete"})
		assert.Error(t, err)
	})
}
