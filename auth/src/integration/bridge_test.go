// Package integration provides integration tests for bridge token flows.
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

	"gorm.io/gorm"
)

func setupBridgeRouter(db *gorm.DB) *echo.Echo {
	router := echo.New()
	router.POST("/basic/signup", controllers.CreateBasicUser)
	router.POST("/basic/login", controllers.LoginBasicUser)
	router.GET("/me", controllers.GetMe, middlewares.RequireAuth)
	router.POST("/bridge/issue", controllers.IssueBridgeToken, middlewares.RequireAuth)
	router.GET("/bridge/exchange", controllers.ExchangeBridgeToken)
	return router
}

// TestBridgeTokenFlow はブリッジトークンの発行・交換フローをテストします
func TestBridgeTokenFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	require.NoError(t, db.AutoMigrate(&models.BridgeToken{}))
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()

	router := setupBridgeRouter(db)

	t.Run("issue and exchange bridge token", func(t *testing.T) {
		testtool.RecordResult(t)

		// 1. ユーザーを登録してセッショントークンを取得
		testtool.LogStep(t, "1. Signup user")
		signupBody, _ := json.Marshal(map[string]string{
			"name": "Bridge User", "email": "bridge@example.com", "password": "BridgePass123!",
		})
		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		var signupResp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &signupResp))
		sessionToken := signupResp["token"].(string)
		testtool.LogSuccess(t, "User signed up")

		// 2. ブリッジトークンを発行
		testtool.LogStep(t, "2. Issue bridge token")
		req = httptest.NewRequest(http.MethodPost, "/bridge/issue", nil)
		req.Header.Set("Authorization", sessionToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var bridgeResp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &bridgeResp))
		bridgeToken, ok := bridgeResp["bridge_token"].(string)
		require.True(t, ok)
		assert.NotEmpty(t, bridgeToken)
		testtool.LogSuccess(t, "Bridge token issued")

		// 3. ブリッジトークンを交換してリフレッシュトークンを取得
		testtool.LogStep(t, "3. Exchange bridge token")
		req = httptest.NewRequest(http.MethodGet, "/bridge/exchange", nil)
		req.Header.Set("Authorization", "Bearer "+bridgeToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		var exchangeResp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &exchangeResp))
		assert.NotEmpty(t, exchangeResp["refresh_token"])
		testtool.LogSuccess(t, "Bridge token exchanged")
	})

	t.Run("bridge token is single use", func(t *testing.T) {
		testtool.RecordResult(t)

		signupBody, _ := json.Marshal(map[string]string{
			"name": "OneTime User", "email": "onetime@example.com", "password": "OneTimePass123!",
		})
		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		sessionToken := resp["token"].(string)

		req = httptest.NewRequest(http.MethodPost, "/bridge/issue", nil)
		req.Header.Set("Authorization", sessionToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		bridgeToken := resp["bridge_token"].(string)

		// 1回目の交換は成功
		req = httptest.NewRequest(http.MethodGet, "/bridge/exchange", nil)
		req.Header.Set("Authorization", "Bearer "+bridgeToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)

		// 2回目は失敗（ワンタイムトークン）
		req = httptest.NewRequest(http.MethodGet, "/bridge/exchange", nil)
		req.Header.Set("Authorization", "Bearer "+bridgeToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testtool.LogSuccess(t, "Bridge token correctly rejected on second use")
	})

	t.Run("bridge exchange without token returns 400", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodGet, "/bridge/exchange", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("bridge issue without auth returns 401", func(t *testing.T) {
		testtool.RecordResult(t)
		req := httptest.NewRequest(http.MethodPost, "/bridge/issue", nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}
