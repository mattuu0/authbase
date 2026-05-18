// Package integration provides integration tests for the /userinfo endpoint.
package integration

import (
	"auth/controllers"
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

func setupUserInfoRouter() *echo.Echo {
	router := echo.New()

	router.POST("/basic/signup", controllers.CreateBasicUser)
	router.POST("/basic/login", controllers.LoginBasicUser)
	router.GET("/token", controllers.GetToken, func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			if token == "" {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
			session, err := services.GetSession(token)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
			}
			c.Set("session", session)
			return next(c)
		}
	})
	router.GET("/userinfo", controllers.GetUserInfo)

	return router
}

// TestUserInfoFlow はサインアップ → セッショントークン取得 → アクセストークン取得 → /userinfo の完全フローをテストします
func TestUserInfoFlow(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)

	services.Init()

	router := setupUserInfoRouter()

	t.Run("signup to userinfo full flow", func(t *testing.T) {
		testtool.RecordResult(t)

		// 1. サインアップしてセッショントークンを取得
		testtool.LogStep(t, "1. Signing up user")
		signupBody, _ := json.Marshal(map[string]string{
			"name":     "UserInfo Tester",
			"email":    "userinfoflow@example.com",
			"password": "FlowPass123!",
		})

		req := httptest.NewRequest(http.MethodPost, "/basic/signup", bytes.NewBuffer(signupBody))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var signupResp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &signupResp))
		sessionToken := signupResp["token"].(string)
		testtool.LogSuccess(t, "Session token obtained")

		// 2. セッショントークンでアクセストークンを取得
		testtool.LogStep(t, "2. Getting access token")
		req = httptest.NewRequest(http.MethodGet, "/token", nil)
		req.Header.Set("Authorization", sessionToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var tokenResp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &tokenResp))
		accessToken := tokenResp["token"].(string)
		assert.NotEmpty(t, accessToken)
		testtool.LogSuccess(t, "Access token obtained")

		// 3. アクセストークンで /userinfo にアクセス
		testtool.LogStep(t, "3. Fetching /userinfo with access token")
		req = httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)

		var userInfo map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &userInfo))

		assert.Equal(t, "userinfoflow@example.com", userInfo["email"])
		assert.Equal(t, "UserInfo Tester", userInfo["name"])
		assert.NotEmpty(t, userInfo["user_id"])
		assert.Equal(t, "basic", userInfo["prov_code"])
		assert.NotNil(t, userInfo["exp"])
		testtool.LogSuccess(t, "/userinfo returned correct user data")

		// 4. セッショントークンで /userinfo はアクセス不可（署名方式が違う）
		testtool.LogStep(t, "4. Verifying session token is rejected at /userinfo")
		req = httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+sessionToken)
		rec = httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		testtool.LogSuccess(t, "Session token correctly rejected at /userinfo")
	})
}

// TestUserInfoWithLabels はラベル付きユーザーの /userinfo レスポンスを検証します
func TestUserInfoWithLabels(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()

	t.Run("labels appear in userinfo response", func(t *testing.T) {
		testtool.RecordResult(t)

		// ラベルとユーザーを準備
		testtool.CreateTestLabel(t, db, "moderator", "#aabbcc")
		testtool.CreateTestLabel(t, db, "verified", "#112233")
		user := testtool.CreateTestUser(t, db, "labeled-userinfo@example.com", models.Basic)
		require.NoError(t, user.AddLabel("moderator"))
		require.NoError(t, user.AddLabel("verified"))

		// アクセストークン取得
		accessToken, err := services.GetAccessToken(user.UserID)
		require.NoError(t, err)

		router := setupUserInfoRouter()
		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var userInfo map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &userInfo))

		labels, ok := userInfo["labels"].([]interface{})
		require.True(t, ok)
		labelStrs := make([]string, len(labels))
		for i, l := range labels {
			labelStrs[i] = l.(string)
		}
		assert.ElementsMatch(t, []string{"moderator", "verified"}, labelStrs)
		testtool.LogSuccess(t, "Labels returned correctly in /userinfo")
	})
}

// TestUserInfoTokenExpiry はアクセストークンの有効期限フィールドを検証します
func TestUserInfoTokenExpiry(t *testing.T) {
	testtool.RecordResult(t)
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	provider := testtool.CreateTestProvider(t, db, models.Basic)
	provider.IsEnabled = 1
	db.Save(provider)
	services.Init()

	t.Run("exp field is in the future", func(t *testing.T) {
		testtool.RecordResult(t)

		user := testtool.CreateTestUser(t, db, "expiry@example.com", models.Basic)
		accessToken, err := services.GetAccessToken(user.UserID)
		require.NoError(t, err)

		router := setupUserInfoRouter()
		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var userInfo map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &userInfo))

		exp, ok := userInfo["exp"].(float64)
		require.True(t, ok)
		assert.Greater(t, int64(exp), int64(0))
		testtool.LogSuccess(t, "exp field is valid")
	})
}
