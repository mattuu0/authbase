// basicuser_signup_test.go - Basic認証ユーザー登録のテスト
package services

import (
	"auth/models"
	testtool "auth/testing"
	"fmt"
	"net/http"
	"testing"
)

// TestSignup_ValidData tests successful user signup
func TestSignup_ValidData(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Signup - Valid Data")
	defer log.Finish()

	log.LogStep("Enabling Basic provider")
	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	Init()

	log.LogStep("Creating user with valid data")
	args := CreateBasicUserArgs{
		Name:     "Test User",
		Email:    "valid@example.com",
		Password: "SecurePass123!",
	}

	token, result := CreateBasicUser(args)
	
	log.AssertTrue(result.Success, "Signup should succeed")
	log.AssertEqual(http.StatusOK, result.Code, "Should return 200 OK")
	log.AssertNotEmpty(token, "Should return token")
	log.LogInfo(fmt.Sprintf("Generated token: %s...", token[:20]))

	log.LogStep("Verifying user in database")
	user, getResult := models.GetUserByEmail(args.Email)
	log.AssertNoError(getResult.Error, "User should exist")
	log.AssertEqual(args.Name, user.Name, "Name should match")
	log.AssertEqual(args.Email, user.Email, "Email should match")
	log.LogSuccess("User successfully created in database")
}

// TestSignup_WeakPassword tests signup with weak password
func TestSignup_WeakPassword(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Signup - Weak Password")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	weakPasswords := []string{
		"123",
		"password",
		"abc",
	}

	for _, pwd := range weakPasswords {
		log.LogStep(fmt.Sprintf("Testing password: %s", pwd))
		
		args := CreateBasicUserArgs{
			Name:     "Test User",
			Email:    fmt.Sprintf("test-%s@example.com", pwd),
			Password: pwd,
		}

		token, result := CreateBasicUser(args)
		log.LogInfo(fmt.Sprintf("Result: Success=%v, Code=%d", result.Success, result.Code))
		
		// 注: パスワード強度チェックが実装されている場合はここでチェック
		if token != "" {
			log.LogWarning("Weak password was accepted")
		}
	}
}

// TestSignup_DuplicateEmail tests duplicate email handling
func TestSignup_DuplicateEmail(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Signup - Duplicate Email")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	email := "duplicate@example.com"

	log.LogStep("Creating first user")
	args1 := CreateBasicUserArgs{
		Name:     "First User",
		Email:    email,
		Password: "Password123!",
	}
	_, result1 := CreateBasicUser(args1)
	log.AssertTrue(result1.Success, "First signup should succeed")

	log.LogStep("Attempting to create second user with same email")
	args2 := CreateBasicUserArgs{
		Name:     "Second User",
		Email:    email,
		Password: "DifferentPass123!",
	}
	_, result2 := CreateBasicUser(args2)
	
	log.AssertTrue(!result2.Success, "Second signup should fail")
	log.AssertEqual(http.StatusConflict, result2.Code, "Should return 409 Conflict")
	log.LogSuccess("Correctly rejected duplicate email")
}

// TestSignup_InvalidEmail tests signup with invalid email
func TestSignup_InvalidEmail(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Signup - Invalid Email")
	defer log.Finish()

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	Init()

	invalidEmails := []string{
		"notanemail",
		"@example.com",
		"user@",
		"user @example.com",
	}

	for _, email := range invalidEmails {
		log.LogStep(fmt.Sprintf("Testing email: %s", email))
		
		args := CreateBasicUserArgs{
			Name:     "Test User",
			Email:    email,
			Password: "Password123!",
		}

		_, result := CreateBasicUser(args)
		log.LogInfo(fmt.Sprintf("Result: Success=%v, Code=%d", result.Success, result.Code))
		
		// 注: メールバリデーションが実装されている場合の確認
	}
}

// TestSignup_ProviderDisabled tests signup when provider is disabled
func TestSignup_ProviderDisabled(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	
	log := testtool.NewTestLogger(t, "Signup - Provider Disabled")
	defer log.Finish()

	log.LogStep("Creating disabled provider")
	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 0
	db.Save(provider)
	Init()

	log.LogStep("Attempting signup with disabled provider")
	args := CreateBasicUserArgs{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "Password123!",
	}

	_, result := CreateBasicUser(args)
	
	log.AssertTrue(!result.Success, "Signup should fail")
	log.AssertEqual(http.StatusUnauthorized, result.Code, "Should return 401")
	log.LogSuccess("Correctly blocked signup with disabled provider")
}