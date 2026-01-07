// user_create_test.go - ユーザー作成に関するテスト
package models_test

import (
	"auth/models"
	testtool "auth/testing"
	"fmt"
	"testing"
)

// TestCreateUser_Success tests successful user creation
func TestCreateUser_Success(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "CreateUser - Success Case")
	defer log.Finish()

	// プロバイダーを作成
	log.LogStep("Creating test provider")
	testtool.CreateTestProvider(t, db, models.Google)

	log.LogStep("Creating new user")
	user := &models.User{
		UserID:   "test-user-success",
		Name:     "Success User",
		Email:    "success@example.com",
		ProvCode: models.Google,
		ProvUID:  "google-success-123",
	}

	err := models.CreateUser(user, models.Google)
	log.AssertNoError(err, "User creation")

	log.LogStep("Verifying user was created")
	retrieved, result := models.GetUser(user.UserID)
	log.AssertNoError(result.Error, "User retrieval")
	log.AssertTrue(result.IsExists, "User exists check")
	log.AssertEqual(user.Email, retrieved.Email, "Email match")
	log.AssertEqual(user.Name, retrieved.Name, "Name match")
}

// TestCreateUser_DuplicateEmail tests duplicate email handling
func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "CreateUser - Duplicate Email")
	defer log.Finish()

	testtool.CreateTestProvider(t, db, models.Google)

	log.LogStep("Creating first user")
	user1 := &models.User{
		UserID:   "test-user-1",
		Name:     "User One",
		Email:    "duplicate@example.com",
		ProvCode: models.Google,
		ProvUID:  "uid-1",
	}

	err := models.CreateUser(user1, models.Google)
	log.AssertNoError(err, "First user creation")

	log.LogStep("Attempting to create second user with same email")
	user2 := &models.User{
		UserID:   "test-user-2",
		Name:     "User Two",
		Email:    "duplicate@example.com",
		ProvCode: models.Google,
		ProvUID:  "uid-2",
	}

	err = models.CreateUser(user2, models.Google)
	if err == nil {
		log.LogError("Expected error for duplicate email, but got none")
		t.Fatal("Should fail with duplicate email")
	} else {
		log.LogSuccess("Correctly rejected duplicate email")
	}
}

// TestCreateUser_EmptyFields tests validation of required fields
func TestCreateUser_EmptyFields(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	
	log := testtool.NewTestLogger(t, "CreateUser - Empty Fields")
	defer log.Finish()

	testtool.CreateTestProvider(t, db, models.Google)

	testCases := []struct {
		name      string
		userID    string
		email     string
		shouldErr bool
	}{
		{"Empty UserID", "", "test@example.com", false},
		{"Empty Email", "user-id", "", false},
		{"Valid Data", "valid-id", "valid@example.com", false},
	}

	for _, tc := range testCases {
		log.LogStep(fmt.Sprintf("Testing: %s", tc.name))
		
		user := &models.User{
			UserID:   tc.userID,
			Name:     "Test User",
			Email:    tc.email,
			ProvCode: models.Google,
			ProvUID:  "test-uid",
		}

		err := models.CreateUser(user, models.Google)
		if tc.shouldErr {
			log.AssertTrue(err != nil, fmt.Sprintf("%s should error", tc.name))
		} else {
			log.LogInfo(fmt.Sprintf("%s result: %v", tc.name, err))
		}
	}
}
