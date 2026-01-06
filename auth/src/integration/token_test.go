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

// TestAccessTokenFlow tests the access token generation flow
func TestAccessTokenFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("get access token", func(t *testing.T) {
		testtool.RecordResult(t)
		// 1. ユーザー登録
		testtool.LogStep(t, "1. Signing up user")
		signupBody := map[string]string{
			"name":     "Token Test User",
			"email":    "token@example.com",
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
		sessionToken := signupResponse["token"].(string)
		testtool.LogSuccess(t, "User signed up")

		// 2. アクセストークンを取得
		testtool.LogStep(t, "2. Requesting access token using session token")
		req = httptest.NewRequest(http.MethodGet, "/token", nil)
		req.Header.Set("Authorization", sessionToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var tokenResponse map[string]interface{}
		err := json.Unmarshal(rec.Body.Bytes(), &tokenResponse)
		require.NoError(t, err)

		accessToken := tokenResponse["token"].(string)
		assert.NotEmpty(t, accessToken)
		assert.NotEqual(t, sessionToken, accessToken)
		testtool.LogSuccess(t, "Access token retrieved successfully")
	})
}
