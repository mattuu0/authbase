package integration

import (
	"auth/models"
	"auth/services"
	testtool "auth/testing"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBasicAuthFlow tests the complete basic authentication flow
func TestBasicAuthFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	testtool.SetupTestEnv(t)

	// Basicプロバイダーを有効化
	testtool.LogStep(t, "Enabling Basic Provider")
	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("complete signup and login flow", func(t *testing.T) {
		testtool.RecordResult(t)
		// 1. ユーザー登録
		testtool.LogStep(t, "1. Signing up new user")
		signupBody := map[string]string{
			"name":     "Integration Test User",
			"email":    "integration@example.com",
			"password": "SecurePassword123!",
		}
		signupJSON, _ := json.Marshal(signupBody)

		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var signupResponse map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &signupResponse)
		require.NoError(t, err)

		signupToken := signupResponse["token"].(string)
		assert.NotEmpty(t, signupToken)
		testtool.LogSuccess(t, "User signed up")

		// 2. 登録したトークンで /me にアクセス
		testtool.LogStep(t, "2. Accessing /me with signup token")
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", signupToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var meResponse map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &meResponse)
		require.NoError(t, err)
		assert.Equal(t, "integration@example.com", meResponse["email"])
		testtool.LogSuccess(t, "Access verified")

		// 3. ログアウト
		testtool.LogStep(t, "3. Logging out")
		req = httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", signupToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Logged out")

		// 4. ログアウト後は /me にアクセスできない
		testtool.LogStep(t, "4. Verifying token invalidation after logout")
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", signupToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testtool.LogSuccess(t, "Access denied as expected")

		// 5. 再度ログイン
		testtool.LogStep(t, "5. Logging in again")
		loginBody := map[string]string{
			"email":    "integration@example.com",
			"password": "SecurePassword123!",
		}
		loginJSON, _ := json.Marshal(loginBody)

		req = httptest.NewRequest(http.MethodPost, "/basic/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var loginResponse map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &loginResponse)
		require.NoError(t, err)

		newToken := loginResponse["token"].(string)
		assert.NotEmpty(t, newToken)
		testtool.LogSuccess(t, "Logged in again")

		// 6. 新しいトークンで /me にアクセス
		testtool.LogStep(t, "6. Accessing /me with new token")
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", newToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Access verified with new token")
	})
}
