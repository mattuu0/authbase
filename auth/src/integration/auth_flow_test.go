// Package integration provides integration tests for complete authentication flows.
// Tests cover end-to-end scenarios including user creation, login, token usage, and logout.
package integration

import (
	"auth/controllers"
	"auth/middlewares"
	"auth/models"
	"auth/services"
	testtool "auth/testing"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRouter はテスト用のルーターをセットアップします
func setupTestRouter() *echo.Echo {
	router := echo.New()

	// Basic認証エンドポイント
	router.POST("/basic/signup", controllers.CreateBasicUser)
	router.POST("/basic/login", controllers.LoginBasicUser)

	// 認証が必要なエンドポイント
	router.GET("/me", controllers.GetMe, middlewares.RequireAuth)
	router.POST("/logout", controllers.Logout, middlewares.RequireAuth)
	router.GET("/token", controllers.GetToken, middlewares.RequireAuth)

	return router
}

// TestBasicAuthFlow tests the complete basic authentication flow
func TestBasicAuthFlow(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	// データベース接続を置き換え
	models.ReplaceDB(db)

	testtool.SetupTestEnv(t)

	// Basicプロバイダーを有効化
	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("complete signup and login flow", func(t *testing.T) {
		// 1. ユーザー登録
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

		// 2. 登録したトークンで /me にアクセス
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", signupToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var meResponse map[string]interface{}
		err = json.Unmarshal(rec.Body.Bytes(), &meResponse)
		require.NoError(t, err)
		assert.Equal(t, "integration@example.com", meResponse["email"])

		// 3. ログアウト
		req = httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", signupToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		// 4. ログアウト後は /me にアクセスできない
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", signupToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		// 5. 再度ログイン
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

		// 6. 新しいトークンで /me にアクセス
		req = httptest.NewRequest(http.MethodGet, "/me", nil)
		req.Header.Set("Authorization", newToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// TestAccessTokenFlow tests the access token generation flow
func TestAccessTokenFlow(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("get access token", func(t *testing.T) {
		// 1. ユーザー登録
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

		// 2. アクセストークンを取得
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
	})
}

// TestMultipleSessionsFlow tests managing multiple concurrent sessions
func TestMultipleSessionsFlow(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("multiple sessions for same user", func(t *testing.T) {
		// 1. ユーザー登録
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

		// 2. 別のデバイスからログイン（2つ目のセッション）
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

		// トークンが異なることを確認
		assert.NotEqual(t, token1, token2)

		// 3. 両方のトークンで /me にアクセスできることを確認
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

		// 4. 1つ目のセッションをログアウト
		req = httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", token1)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 5. 1つ目のトークンは無効、2つ目は有効
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
	})
}

// TestErrorHandling tests various error scenarios
func TestErrorHandling(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupTestRouter()

	t.Run("invalid credentials", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    "nonexistent@example.com",
			"password": "WrongPassword123!",
		}
		loginJSON, _ := json.Marshal(loginBody)

		req := httptest.NewRequest(http.MethodPost, "/basic/login", bytes.NewBuffer(loginJSON))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.NotEqual(t, http.StatusOK, rec.Code)
	})

	t.Run("malformed request body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("missing authorization header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/me", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
