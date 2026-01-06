// Package services_test provides unit tests for basic user authentication services.
// Tests cover user creation, login, password hashing, and error handling.
package services

import (
	"auth/models"
	testtool "auth/testing"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateBasicUser tests the creation of a basic authentication user
func TestCreateBasicUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	// 既存のデータベース接続を差し替え
	models.ReplaceDB(db)

	testtool.SetupTestEnv(t)

	// Basicプロバイダーを有効化
	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	// サービスを初期化
	Init()

	t.Run("successful user creation", func(t *testing.T) {
		args := CreateBasicUserArgs{
			Name:     "Test User",
			Email:    "newuser@example.com",
			Password: "SecurePassword123!",
		}

		token, result := CreateBasicUser(args)

		assert.True(t, result.Success)
		assert.Equal(t, http.StatusOK, result.Code)
		assert.NotEmpty(t, token)
		assert.NoError(t, result.Error)

		// ユーザーが作成されたことを確認
		user, getResult := models.GetUserByEmail(args.Email)
		require.NoError(t, getResult.Error)
		assert.Equal(t, args.Name, user.Name)
		assert.Equal(t, args.Email, user.Email)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("duplicate email", func(t *testing.T) {
		email := "duplicate@example.com"

		// 1人目を作成
		args1 := CreateBasicUserArgs{
			Name:     "First User",
			Email:    email,
			Password: "Password123!",
		}
		_, result1 := CreateBasicUser(args1)
		require.True(t, result1.Success)

		// 同じメールで2人目を作成しようとする
		args2 := CreateBasicUserArgs{
			Name:     "Second User",
			Email:    email,
			Password: "DifferentPassword123!",
		}
		_, result2 := CreateBasicUser(args2)

		assert.False(t, result2.Success)
		assert.Equal(t, http.StatusConflict, result2.Code)
		assert.Error(t, result2.Error)
	})

	t.Run("provider disabled", func(t *testing.T) {
		// プロバイダーを無効化
		provider.IsEnabled = 0
		db.Save(provider)

		args := CreateBasicUserArgs{
			Name:     "Test User",
			Email:    "disabled@example.com",
			Password: "Password123!",
		}

		_, result := CreateBasicUser(args)

		assert.False(t, result.Success)
		assert.Equal(t, http.StatusUnauthorized, result.Code)
		assert.Error(t, result.Error)

		// プロバイダーを再度有効化（後続のテストのため）
		provider.IsEnabled = 1
		db.Save(provider)
	})
}

// TestLoginBasicUser tests basic user login functionality
func TestLoginBasicUser(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	// Basicプロバイダーを有効化
	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	Init()

	// テストユーザーを作成
	email := "logintest@example.com"
	password := "TestPassword123!"
	createArgs := CreateBasicUserArgs{
		Name:     "Login Test User",
		Email:    email,
		Password: password,
	}
	_, createResult := CreateBasicUser(createArgs)
	require.True(t, createResult.Success)

	t.Run("successful login", func(t *testing.T) {
		loginArgs := LoginBasicUserArgs{
			Email:    email,
			Password: password,
		}

		token, result := LoginBasicUser(loginArgs)

		assert.True(t, result.Success)
		assert.Equal(t, http.StatusOK, result.Code)
		assert.NotEmpty(t, token)
		assert.NoError(t, result.Error)
	})

	t.Run("wrong password", func(t *testing.T) {
		loginArgs := LoginBasicUserArgs{
			Email:    email,
			Password: "WrongPassword123!",
		}

		_, result := LoginBasicUser(loginArgs)

		assert.False(t, result.Success)
		assert.Equal(t, http.StatusBadRequest, result.Code)
		assert.Error(t, result.Error)
	})

	t.Run("non-existent user", func(t *testing.T) {
		loginArgs := LoginBasicUserArgs{
			Email:    "nonexistent@example.com",
			Password: "Password123!",
		}

		_, result := LoginBasicUser(loginArgs)

		assert.False(t, result.Success)
		assert.Error(t, result.Error)
	})

	t.Run("provider disabled", func(t *testing.T) {
		// プロバイダーを無効化
		provider.IsEnabled = 0
		db.Save(provider)

		loginArgs := LoginBasicUserArgs{
			Email:    email,
			Password: password,
		}

		_, result := LoginBasicUser(loginArgs)

		assert.False(t, result.Success)
		assert.Equal(t, http.StatusUnauthorized, result.Code)
		assert.Error(t, result.Error)
	})
}

// TestPasswordHashing tests password hashing and verification
func TestPasswordHashing(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	Init()

	t.Run("password is hashed", func(t *testing.T) {
		password := "PlainTextPassword123!"
		args := CreateBasicUserArgs{
			Name:     "Hash Test User",
			Email:    "hashtest@example.com",
			Password: password,
		}

		_, result := CreateBasicUser(args)
		require.True(t, result.Success)

		// パスワードが平文で保存されていないことを確認
		user, getResult := models.GetUserByEmail(args.Email)
		require.NoError(t, getResult.Error)
		assert.NotEqual(t, password, user.PasswordHash)
		assert.NotEmpty(t, user.PasswordHash)
	})

	t.Run("password verification works", func(t *testing.T) {
		password := "VerifyPassword123!"
		email := "verify@example.com"

		// ユーザーを作成
		createArgs := CreateBasicUserArgs{
			Name:     "Verify User",
			Email:    email,
			Password: password,
		}
		_, createResult := CreateBasicUser(createArgs)
		require.True(t, createResult.Success)

		// 正しいパスワードでログイン
		loginArgs := LoginBasicUserArgs{
			Email:    email,
			Password: password,
		}
		_, loginResult := LoginBasicUser(loginArgs)
		assert.True(t, loginResult.Success)

		// 間違ったパスワードでログイン
		wrongArgs := LoginBasicUserArgs{
			Email:    email,
			Password: "WrongPassword123!",
		}
		_, wrongResult := LoginBasicUser(wrongArgs)
		assert.False(t, wrongResult.Success)
	})
}