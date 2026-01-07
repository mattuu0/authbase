// user_label_test.go - ユーザーラベル管理に関するテスト
package models_test

import (
	"auth/models"
	testtool "auth/testing"
	"fmt"
	"testing"
)

// TestAddLabel_Single tests adding a single label to user
func TestAddLabel_Single(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "AddLabel - Single")
	defer log.Finish()

	log.LogStep("Creating test user and label")
	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "label@example.com", models.Google)
	label := testtool.CreateTestLabel(t, db, "premium", "#FFD700")

	log.LogStep("Adding label to user")
	err := user.AddLabel(label.Name)
	log.AssertNoError(err, "Add label")

	log.LogStep("Verifying label was added")
	labels, err := user.GetLabels()
	log.AssertNoError(err, "Get labels")
	log.AssertEqual(1, len(labels), "Should have 1 label")
	log.AssertEqual("premium", labels[0].Name, "Label name should match")
}

// TestAddLabel_Multiple tests adding multiple labels
func TestAddLabel_Multiple(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "AddLabel - Multiple")
	defer log.Finish()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "multilabel@example.com", models.Google)

	log.LogStep("Creating multiple labels")
	labels := []string{"admin", "moderator", "premium"}
	for _, name := range labels {
		testtool.CreateTestLabel(t, db, name, "#000000")
	}

	log.LogStep("Adding all labels to user")
	for _, name := range labels {
		err := user.AddLabel(name)
		log.AssertNoError(err, fmt.Sprintf("Add label: %s", name))
	}

	log.LogStep("Verifying all labels")
	userLabels, err := user.GetLabels()
	log.AssertNoError(err, "Get labels")
	log.AssertEqual(len(labels), len(userLabels), "Label count should match")
	log.LogInfo(fmt.Sprintf("User has %d labels", len(userLabels)))
}

// TestRemoveLabel tests removing a label from user
func TestRemoveLabel(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "RemoveLabel")
	defer log.Finish()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "removelabel@example.com", models.Google)
	label := testtool.CreateTestLabel(t, db, "temporary", "#FF0000")

	log.LogStep("Adding label")
	err := user.AddLabel(label.Name)
	log.AssertNoError(err, "Add label")

	log.LogStep("Verifying label exists")
	labels, _ := user.GetLabels()
	log.AssertEqual(1, len(labels), "Should have 1 label")

	log.LogStep("Removing label")
	err = user.RemoveLabel(label.Name)
	log.AssertNoError(err, "Remove label")

	log.LogStep("Verifying label removed")
	labels, _ = user.GetLabels()
	log.AssertEqual(0, len(labels), "Should have 0 labels")
}

// TestRemoveAllLabels tests removing all labels from user
func TestRemoveAllLabels(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "RemoveAllLabels")
	defer log.Finish()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "removeall@example.com", models.Google)

	log.LogStep("Adding multiple labels")
	testtool.CreateTestLabel(t, db, "label1", "#111111")
	testtool.CreateTestLabel(t, db, "label2", "#222222")
	testtool.CreateTestLabel(t, db, "label3", "#333333")
	
	user.AddLabel("label1")
	user.AddLabel("label2")
	user.AddLabel("label3")

	labels, _ := user.GetLabels()
	log.LogInfo(fmt.Sprintf("Added %d labels", len(labels)))

	log.LogStep("Removing all labels")
	err := user.RemoveAllLabels()
	log.AssertNoError(err, "Remove all labels")

	log.LogStep("Verifying all labels removed")
	labels, _ = user.GetLabels()
	log.AssertEqual(0, len(labels), "Should have 0 labels")
}

// TestGetLabelNames tests getting label names
func TestGetLabelNames(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "GetLabelNames")
	defer log.Finish()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "labelnames@example.com", models.Google)

	log.LogStep("Creating and adding labels")
	expectedNames := []string{"admin", "moderator", "premium"}
	for _, name := range expectedNames {
		testtool.CreateTestLabel(t, db, name, "#000000")
		user.AddLabel(name)
	}

	log.LogStep("Getting label names")
	names, err := user.GetLabelNames()
	log.AssertNoError(err, "Get label names")
	log.AssertEqual(len(expectedNames), len(names), "Name count should match")
	
	log.LogStep("Verifying each label name")
	for i, name := range names {
		log.LogInfo(fmt.Sprintf("Label %d: %s", i+1, name))
	}
}
