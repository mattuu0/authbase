// Package testing provides test utilities and setup functions for the auth service.
// This package contains helper functions to initialize test databases,
// mock configurations, and common test fixtures.
package testing

import (
	"auth/models"
	"fmt"
	"os"
	"testing"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// LogStep prints a formatted step message to the test log
func LogStep(t *testing.T, msg string) {
	t.Helper()
	t.Logf("\n🔹 %s ...\n", msg)
}

// LogSuccess prints a formatted success message
func LogSuccess(t *testing.T, msg string) {
	t.Helper()
	t.Logf("✅ %s\n", msg)
}

// TestDB はテスト用のデータベース接続を保持します
var TestDB *gorm.DB

// SetupTestDB はテスト用のSQLiteデータベースをセットアップします
// 各テストの開始時に呼び出されることを想定しています
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	LogStep(t, "Initializing Test Database (In-Memory)")

	// SQLiteデータベースを作成 (インメモリ)
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// テーブルのマイグレーション
	err = db.AutoMigrate(
		&models.User{},
		&models.Provider{},
		&models.Session{},
		&models.Label{},
		&models.AdminUser{},
	)
	if err != nil {
		t.Fatalf("Failed to migrate test database: %v", err)
	}

	models.SetDB(db)
	TestDB = db
	
	LogSuccess(t, "Test Database initialized")
	return db
}

// CleanupTestDB はテスト用データベースをクリーンアップします
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	t.Helper()
	
	LogStep(t, "Cleaning up Test Database")

	sqlDB, err := db.DB()
	if err != nil {
		t.Logf("Failed to get sql.DB: %v", err)
		return
	}

	if err := sqlDB.Close(); err != nil {
		t.Logf("Failed to close database: %v", err)
	}
	
	LogSuccess(t, "Test Database cleaned up")
}

// SetupTestEnv はテスト用の環境変数をセットアップします
func SetupTestEnv(t *testing.T) {
	t.Helper()

	// JWT用の秘密鍵と公開鍵（テスト用のダミー値）
	os.Setenv("JWT_PRIVATE_KEY", `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIJ+DYvh6SEqVTm50DFtMDoQikTmiCqirVv9mWG9qfSnF
-----END PRIVATE KEY-----`)

	os.Setenv("JWT_PUBLIC_KEY", `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAGb9ECWmEzf6FQbrBZ9w7lshQhqowtrbLDFw4rXAxZuE=
-----END PUBLIC KEY-----`)

	os.Setenv("TOKEN_SECRET", "test-secret-key-for-testing-purposes-only")
	os.Setenv("ADMIN_SESSION_KEY", "test-admin-session-key")
	os.Setenv("DB_DSN", ":memory:")
}

// CreateTestProvider はテスト用のプロバイダーを作成します
func CreateTestProvider(t *testing.T, db *gorm.DB, providerCode models.ProviderCode) *models.Provider {
	t.Helper()

	provider := &models.Provider{
		ProviderName: string(providerCode),
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		CallbackURL:  fmt.Sprintf("http://localhost/callback/%s", providerCode),
		ProviderCode: providerCode,
		IsEnabled:    1,
	}

	if err := db.Create(provider).Error; err != nil {
		t.Fatalf("Failed to create test provider: %v", err)
	}

	return provider
}

// CreateTestUser はテスト用のユーザーを作成します
func CreateTestUser(t *testing.T, db *gorm.DB, email string, providerCode models.ProviderCode) *models.User {
	t.Helper()

	user := &models.User{
		UserID:   fmt.Sprintf("test-user-%s", email),
		Name:     "Test User",
		Email:    email,
		ProvCode: providerCode,
		ProvUID:  fmt.Sprintf("prov-uid-%s", email),
	}

	if err := models.CreateUser(user, providerCode); err != nil {
		t.Fatalf("Failed to create test user: %v", err)
	}

	return user
}

// CreateTestAdminUser はテスト用の管理者ユーザーを作成します
func CreateTestAdminUser(t *testing.T, db *gorm.DB, username, password string) *models.AdminUser {
	t.Helper()

	adminUser := &models.AdminUser{
		UserID:       fmt.Sprintf("admin-%s", username),
		Username:     username,
		PasswordHash: password, // テスト用なので平文
		IsSystem:     0,
	}

	if err := db.Create(adminUser).Error; err != nil {
		t.Fatalf("Failed to create test admin user: %v", err)
	}

	return adminUser
}

// CreateTestLabel はテスト用のラベルを作成します
func CreateTestLabel(t *testing.T, db *gorm.DB, name, color string) *models.Label {
	t.Helper()

	label := &models.Label{
		Name:  name,
		Color: color,
	}

	if err := db.Create(label).Error; err != nil {
		t.Fatalf("Failed to create test label: %v", err)
	}

	return label
}