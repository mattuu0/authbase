// Package testing provides constants and fixtures for testing.
// This file contains test keys and common test data.
package testing

// TestJWTPrivateKey はテスト用のEd25519秘密鍵です
const TestJWTPrivateKey = `-----BEGIN PRIVATE KEY-----
MC4CAQAwBQYDK2VwBCIEIJ+DYvh6SEqVTm50DFtMDoQikTmiCqirVv9mWG9qfSnF
-----END PRIVATE KEY-----`

// TestJWTPublicKey はテスト用のEd25519公開鍵です
const TestJWTPublicKey = `-----BEGIN PUBLIC KEY-----
MCowBQYDK2VwAyEAGb9ECWmEzf6FQbrBZ9w7lshQhqowtrbLDFw4rXAxZuE=
-----END PUBLIC KEY-----`

// TestUserEmails はテストで使用するメールアドレスのリストです
var TestUserEmails = []string{
	"test1@example.com",
	"test2@example.com",
	"test3@example.com",
	"admin@example.com",
	"user@example.com",
}

// TestPasswords はテストで使用するパスワードのリストです
var TestPasswords = []string{
	"SecurePassword123!",
	"TestPass456!",
	"MyPassword789!",
}

// TestLabelColors はテストで使用するラベル色のリストです
var TestLabelColors = []string{
	"#FF0000", // Red
	"#00FF00", // Green
	"#0000FF", // Blue
	"#FFD700", // Gold
	"#800080", // Purple
}

// TestLabelNames はテストで使用するラベル名のリストです
var TestLabelNames = []string{
	"admin",
	"moderator",
	"premium",
	"verified",
	"beta-tester",
}