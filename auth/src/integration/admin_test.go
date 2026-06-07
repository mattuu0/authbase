// Package integration provides integration tests for admin authentication flows.
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

func setupAdminRouter() *echo.Echo {
	router := echo.New()

	router.POST("/admin/signup", controllers.CreateAdminUser)
	router.POST("/admin/login", controllers.LoginAdminUser)
	router.GET("/admin/status", controllers.GetAdminStatus)
	router.GET("/admin/info", controllers.GetAdminInfo, middlewares.RequireAdminAuth)
	router.POST("/admin/logout", controllers.AdminLogout, middlewares.RequireAdminAuth)

	// admin が必要な API
	apig := router.Group("/api")
	apig.Use(middlewares.RequireAdminAuth)
	userg := apig.Group("/user")
	userg.GET("/all", controllers.GetAllUsers)
	userg.POST("", controllers.CreateUser)
	userg.PUT("", controllers.UpdateUser)
	userg.DELETE("", controllers.DeleteOauth)
	userg.PUT("/ban", controllers.ToggleBan)
	labelg := apig.Group("/labels")
	labelg.GET("", controllers.GetLabels)
	labelg.POST("", controllers.CreateLabel)
	labelg.PUT("", controllers.UpdateLabel)
	labelg.DELETE("", controllers.DeleteLabel)
	sessiong := apig.Group("/session")
	sessiong.GET("", controllers.GetSessions)
	sessiong.DELETE("", controllers.DeleteSession)

	return router
}

// loginAsAdmin はテスト用の管理者ログインを行い、クッキー付きのレスポンスを返します
func loginAsAdmin(t *testing.T, router *echo.Echo, username, password string) []*http.Cookie {
	t.Helper()
	body, _ := json.Marshal(map[string]string{"username": username, "password": password})
	req := httptest.NewRequest(http.MethodPost, "/admin/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code, "admin login should succeed")
	return rec.Result().Cookies()
}

// TestAdminSignupAndLoginFlow はAdmin登録・ログイン・情報取得の完全フローをテストします
func TestAdminSignupAndLoginFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	services.Init()

	router := setupAdminRouter()

	t.Run("admin signup, login, and info flow", func(t *testing.T) {
		testtool.RecordResult(t)

		// 1. 初期ステータスは管理者なし
		testtool.LogStep(t, "1. Check initial admin status")
		req := httptest.NewRequest(http.MethodGet, "/admin/status", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var status map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &status))
		assert.Equal(t, false, status["HasSystemUser"])
		testtool.LogSuccess(t, "No admin user initially")

		// 2. 管理者ユーザーを作成
		testtool.LogStep(t, "2. Create admin user")
		body, _ := json.Marshal(map[string]string{"username": "admin-user", "password": "AdminPass123!"})
		req = httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Admin user created")

		// 3. ステータスが管理者ありになる
		testtool.LogStep(t, "3. Check status after signup")
		req = httptest.NewRequest(http.MethodGet, "/admin/status", nil)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &status))
		assert.Equal(t, true, status["HasSystemUser"])
		testtool.LogSuccess(t, "Admin status updated")

		// 4. ログイン
		testtool.LogStep(t, "4. Login as admin")
		cookies := loginAsAdmin(t, router, "admin-user", "AdminPass123!")
		assert.NotEmpty(t, cookies)
		testtool.LogSuccess(t, "Admin logged in")

		// 5. 管理者情報を取得
		testtool.LogStep(t, "5. Get admin info")
		req = httptest.NewRequest(http.MethodGet, "/admin/info", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var info map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &info))
		assert.Equal(t, "admin-user", info["Username"])
		assert.Empty(t, info["PasswordHash"], "password_hash should not be exposed")
		testtool.LogSuccess(t, "Admin info retrieved correctly")

		// 6. ログアウト
		testtool.LogStep(t, "6. Admin logout")
		req = httptest.NewRequest(http.MethodPost, "/admin/logout", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		testtool.LogSuccess(t, "Admin logged out")

		// 7. ログアウト後は admin/info にアクセス不可
		testtool.LogStep(t, "7. Verify access denied after logout")
		req = httptest.NewRequest(http.MethodGet, "/admin/info", nil)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testtool.LogSuccess(t, "Access denied after logout")
	})
}

// TestAdminDuplicateSignup は管理者の二重登録がエラーになることをテストします
func TestAdminDuplicateSignup(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)
	services.Init()
	router := setupAdminRouter()

	body, _ := json.Marshal(map[string]string{"username": "first-admin", "password": "Pass1234!"})
	req := httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)

	body, _ = json.Marshal(map[string]string{"username": "second-admin", "password": "Pass5678!"})
	req = httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec = httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusInternalServerError, rec.Code, "second admin signup should fail")
}

// TestAdminAPIAccess は管理者APIのアクセス制御をテストします
func TestAdminAPIAccess(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()
	router := setupAdminRouter()

	// 管理者なしで API アクセス
	t.Run("api access without admin session returns 401", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	// 管理者を作成してログイン
	body, _ := json.Marshal(map[string]string{"username": "api-admin", "password": "ApiAdmin123!"})
	req := httptest.NewRequest(http.MethodPost, "/admin/signup", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	cookies := loginAsAdmin(t, router, "api-admin", "ApiAdmin123!")

	// 管理者セッションで API アクセス
	t.Run("api access with admin session succeeds", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodGet, "/api/user/all", nil)
		for _, c := range cookies {
			req.AddCookie(c)
		}
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
