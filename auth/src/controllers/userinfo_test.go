// Package controllers_test provides unit tests for the GetUserInfo controller.
package controllers_test

import (
	"auth/controllers"
	"auth/models"
	"auth/services"
	testtool "auth/testing"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserInfoEcho() *echo.Echo {
	e := echo.New()
	e.GET("/userinfo", controllers.GetUserInfo)
	return e
}

func TestGetUserInfo(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	models.ReplaceDB(db)
	testtool.SetupTestEnv(t)
	services.Init()

	testtool.CreateTestProvider(t, db, models.Basic)
	testtool.CreateTestLabel(t, db, "admin", "#ff0000")

	// アクセストークンを発行するためのユーザーを作成
	user := testtool.CreateTestUser(t, db, "userinfo@example.com", models.Basic)
	require.NoError(t, user.AddLabel("admin"))

	e := setupUserInfoEcho()

	t.Run("valid access token returns user info", func(t *testing.T) {
		accessToken, err := services.GetAccessToken(user.UserID)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()

		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))

		assert.Equal(t, user.UserID, resp["user_id"])
		assert.Equal(t, user.Email, resp["email"])
		assert.Equal(t, user.Name, resp["name"])
		assert.Equal(t, string(models.Basic), resp["prov_code"])

		labels, ok := resp["labels"].([]interface{})
		require.True(t, ok)
		assert.Len(t, labels, 1)
		assert.Equal(t, "admin", labels[0])
	})

	t.Run("Bearer prefix is stripped correctly", func(t *testing.T) {
		accessToken, err := services.GetAccessToken(user.UserID)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("token without Bearer prefix returns 401", func(t *testing.T) {
		accessToken, err := services.GetAccessToken(user.UserID)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", accessToken)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("missing Authorization header returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Contains(t, resp["error"], "Authorization")
	})

	t.Run("invalid token returns 401", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.string")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("session token (HS512) is rejected as access token", func(t *testing.T) {
		// セッショントークン（HS512）はアクセストークン（EdDSA）として無効
		sessionToken, err := services.GenSessionToken("some-session-id")
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+sessionToken)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("response contains exp field", func(t *testing.T) {
		accessToken, err := services.GetAccessToken(user.UserID)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodGet, "/userinfo", nil)
		req.Header.Set("Authorization", "Bearer "+accessToken)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]interface{}
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotNil(t, resp["exp"])

		exp, ok := resp["exp"].(float64)
		require.True(t, ok)
		assert.Greater(t, exp, float64(0))
	})
}
