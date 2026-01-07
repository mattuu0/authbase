// basicuser_login_test.go - Basic認証ログインのテスト
package services

import (
	"auth/models"
	testtool "auth/testing"
	"fmt"
	"net/http"
	"testing"
)

// TestLogin_ValidCredentials tests successful login
func TestLogin_ValidCredentials(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Login - Valid Credentials")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	email := "logintest@example.com"
	password := "TestPass123!"

	log.LogStep("Creating user account")
	createArgs := CreateBasicUserArgs{
		Name:     "Login Test User",
		Email:    email,
		Password: password,
	}
	_, createResult := CreateBasicUser(createArgs)
	log.AssertTrue(createResult.Success, "User creation should succeed")

	log.LogStep("Logging in with correct credentials")
	loginArgs := LoginBasicUserArgs{
		Email:    email,
		Password: password,
	}

	token, result := LoginBasicUser(loginArgs)
	
	log.AssertTrue(result.Success, "Login should succeed")
	log.AssertEqual(http.StatusOK, result.Code, "Should return 200 OK")
	log.AssertNotEmpty(token, "Should return token")
	log.LogInfo(fmt.Sprintf("Login token: %s...", token[:20]))
}

// TestLogin_WrongPassword tests login with incorrect password
func TestLogin_WrongPassword(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Login - Wrong Password")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	email := "wrongpass@example.com"
	correctPassword := "CorrectPass123!"
	wrongPassword := "WrongPass123!"

	log.LogStep("Creating user")
	createArgs := CreateBasicUserArgs{
		Name:     "Test User",
		Email:    email,
		Password: correctPassword,
	}
	_, createResult := CreateBasicUser(createArgs)
	log.AssertTrue(createResult.Success, "User creation")

	log.LogStep("Attempting login with wrong password")
	loginArgs := LoginBasicUserArgs{
		Email:    email,
		Password: wrongPassword,
	}

	_, result := LoginBasicUser(loginArgs)
	
	log.AssertTrue(!result.Success, "Login should fail")
	log.AssertEqual(http.StatusBadRequest, result.Code, "Should return 400")
	log.LogSuccess("Correctly rejected wrong password")
}

// TestLogin_NonExistentUser tests login with non-existent email
func TestLogin_NonExistentUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Login - Non-Existent User")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	log.LogStep("Attempting login with non-existent email")
	loginArgs := LoginBasicUserArgs{
		Email:    "nonexistent@example.com",
		Password: "Password123!",
	}

	_, result := LoginBasicUser(loginArgs)
	
	log.AssertTrue(!result.Success, "Login should fail")
	log.LogSuccess("Correctly rejected non-existent user")
}

// TestLogin_EmptyCredentials tests login with empty credentials
func TestLogin_EmptyCredentials(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Login - Empty Credentials")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	testCases := []struct {
		name     string
		email    string
		password string
	}{
		{"Empty Email", "", "Password123!"},
		{"Empty Password", "test@example.com", ""},
		{"Both Empty", "", ""},
	}

	for _, tc := range testCases {
		log.LogStep(fmt.Sprintf("Testing: %s", tc.name))
		
		loginArgs := LoginBasicUserArgs{
			Email:    tc.email,
			Password: tc.password,
		}

		_, result := LoginBasicUser(loginArgs)
		log.LogInfo(fmt.Sprintf("%s result: Success=%v", tc.name, result.Success))
	}
}

// TestLogin_ProviderDisabled tests login when provider is disabled
func TestLogin_ProviderDisabled(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Login - Provider Disabled")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	email := "provdisabled@example.com"
	password := "Password123!"

	log.LogStep("Creating user while provider is enabled")
	createArgs := CreateBasicUserArgs{
		Name:     "Test User",
		Email:    email,
		Password: password,
	}
	_, createResult := CreateBasicUser(createArgs)
	log.AssertTrue(createResult.Success, "User creation")

	log.LogStep("Disabling provider")
	provider.IsEnabled = 0
	db.Save(provider)

	log.LogStep("Attempting login with disabled provider")
	loginArgs := LoginBasicUserArgs{
		Email:    email,
		Password: password,
	}

	_, result := LoginBasicUser(loginArgs)
	
	log.AssertTrue(!result.Success, "Login should fail")
	log.AssertEqual(http.StatusUnauthorized, result.Code, "Should return 401")
	log.LogSuccess("Correctly blocked login with disabled provider")
}

// TestLogin_MultipleAttempts tests multiple login attempts
func TestLogin_MultipleAttempts(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Login - Multiple Attempts")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	email := "multiattempt@example.com"
	password := "Password123!"

	log.LogStep("Creating user")
	createArgs := CreateBasicUserArgs{
		Name:     "Test User",
		Email:    email,
		Password: password,
	}
	CreateBasicUser(createArgs)

	log.LogStep("Performing multiple login attempts")
	for i := 1; i <= 5; i++ {
		loginArgs := LoginBasicUserArgs{
			Email:    email,
			Password: password,
		}

		token, result := LoginBasicUser(loginArgs)
		log.AssertTrue(result.Success, fmt.Sprintf("Attempt %d should succeed", i))
		log.LogInfo(fmt.Sprintf("Attempt %d: Got token %s...", i, token[:20]))
	}
}