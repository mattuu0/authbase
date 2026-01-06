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

// TestMultipleSessionsFlow tests managing multiple concurrent sessions
func TestMultipleSessionsFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("multiple sessions for same user", func(t *testing.T) {
		testtool.RecordResult(t)
		// 1. ユーザー登録
		testtool.LogStep(t, "1. Signing up user (Session 1)")
		signupBody := map[string]string{
			"name":     "Multi Session User",
			"email":    "multisession@example.com",
			"password": "Password123!",
		}
		signupJSON, _ := json.Marshal(signupBody)

		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var signupResponse map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &signupResponse)
		token1 := signupResponse["token"].(string)
		testtool.LogSuccess(t, "Session 1 created")

		// 2. 別のデバイスからログイン（2つ目のセッション）
		testtool.LogStep(t, "2. Logging in from another device (Session 2)")
		loginBody := map[string]string{
			"email":    "multisession@example.com",
			"password": "Password123!",
		}
		loginJSON, _ := json.Marshal(loginBody)

		req = httptest.NewRequest(http.MethodPost, "/basic/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var loginResponse map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &loginResponse)
		token2 := loginResponse["token"].(string)
		testtool.LogSuccess(t, "Session 2 created")

		// トークンが異なることを確認
		assert.NotEqual(t, token1, token2)

		// 3. 両方のトークンで /me にアクセスできることを確認
		testtool.LogStep(t, "3. Verifying both sessions are active")
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", token1)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", token2)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Both sessions active")

		// 4. 1つ目のセッションをログアウト
		testtool.LogStep(t, "4. Logging out Session 1")
		req = httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", token1)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Session 1 logged out")

		// 5. 1つ目のトークンは無効、2つ目は有効
		testtool.LogStep(t, "5. Verifying Session 1 is invalid and Session 2 is valid")
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", token1)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", token2)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Session validation correct")
	})
}
