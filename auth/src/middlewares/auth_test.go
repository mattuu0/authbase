// Package middlewares_test provides unit tests for authentication middlewares.
// Tests cover authentication validation, banned user handling, and error cases.
package middlewares

import (
	"auth/models"
	"auth/services"
	testtool "auth/testing"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRequireAuth tests the authentication middleware
func TestRequireAuth(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)

	// DB接続を置き換え
	models.ReplaceDB(db)

	testtool.SetupTestEnv(t)

	services.Init()

	// テストユーザーとセッションを作成
	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "middleware@example.com", models.Google)

	// セッションを作成
	token, err := services.NewSession(services.SessionArgs{
		UserID:    user.UserID,
		RemoteIP:  "192.168.1.1",
		UserAgent: "Test Browser",
	})
	require.NoError(t, err)

	// テスト用のハンドラー（認証が成功したら呼ばれる）
	handler := func(c echo.Context) error {
		session := c.Get("session").(*models.Session)
		return c.JSON(http.StatusOK, echo.Map{
			"message": "authenticated",
			"userID":  session.UserID,
		})
	}

	t.Run("valid token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// ミドルウェアを適用
		h := RequireAuth(handler)
		err := h(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "authenticated")
		assert.Contains(t, rec.Body.String(), user.UserID)
	})

	t.Run("missing token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		// Authorizationヘッダーなし
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := RequireAuth(handler)
		err := h(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
		assert.Contains(t, rec.Body.String(), "unauthorized")
	})

	t.Run("invalid token", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "invalid-token-string")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := RequireAuth(handler)
		err := h(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})

	t.Run("banned user", func(t *testing.T) {
		// BANされたユーザーを作成
		bannedUser := testtool.CreateTestUser(t, db, "banned@example.com", models.Google)
		bannedUser.IsBanned = 1
		models.UpdateUser(bannedUser)

		// BANされたユーザーのセッションを作成
		bannedToken, err := services.NewSession(services.SessionArgs{
			UserID:    bannedUser.UserID,
			RemoteIP:  "192.168.1.2",
			UserAgent: "Banned Browser",
		})
		// BANされたユーザーなのでエラーになるはず
		if err == nil {
			// もしセッションが作成できてしまった場合（古いセッション）
			// ミドルウェアがブロックすることをテスト
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header.Set("Authorization", bannedToken)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := RequireAuth(handler)
			err := h(c)

			require.NoError(t, err)
			assert.Equal(t, http.StatusForbidden, rec.Code)
			assert.Contains(t, rec.Body.String(), "banned")
		}
	})

	t.Run("expired or deleted session", func(t *testing.T) {
		// セッションを作成してすぐに削除
		tempUser := testtool.CreateTestUser(t, db, "temp@example.com", models.Google)
		tempToken, err := services.NewSession(services.SessionArgs{
			UserID:    tempUser.UserID,
			RemoteIP:  "192.168.1.3",
			UserAgent: "Temp Browser",
		})
		require.NoError(t, err)

		// セッションを削除
		sessionID, _ := services.ValidateSessionToken(tempToken)
		services.DeleteSession(sessionID)

		// 削除されたセッションでアクセス
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", tempToken)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := RequireAuth(handler)
		err = h(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

// TestRequireAuthContextValues tests that middleware sets correct context values
func TestRequireAuthContextValues(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	services.Init()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "context@example.com", models.Google)

	token, err := services.NewSession(services.SessionArgs{
		UserID:    user.UserID,
		RemoteIP:  "192.168.1.4",
		UserAgent: "Context Test",
	})
	require.NoError(t, err)

	t.Run("context values are set", func(t *testing.T) {
		handler := func(c echo.Context) error {
			// コンテキストから値を取得
			session := c.Get("session")
			assert.NotNil(t, session)

			sessionObj, ok := session.(*models.Session)
			require.True(t, ok)
			assert.Equal(t, user.UserID, sessionObj.UserID)

			return c.JSON(http.StatusOK, echo.Map{"ok": true})
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		h := RequireAuth(handler)
		err := h(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

// TestMultipleMiddlewareChaining tests middleware chaining
func TestMultipleMiddlewareChaining(t *testing.T) {
	db := testtool.SetupTestDB(t)
	defer testtool.CleanupTestDB(t, db)
	testtool.SetupTestEnv(t)

	services.Init()

	testtool.CreateTestProvider(t, db, models.Google)
	user := testtool.CreateTestUser(t, db, "chain@example.com", models.Google)

	token, err := services.NewSession(services.SessionArgs{
		UserID:    user.UserID,
		RemoteIP:  "192.168.1.5",
		UserAgent: "Chain Test",
	})
	require.NoError(t, err)

	t.Run("chained middleware", func(t *testing.T) {
		// カスタムミドルウェア
		customMiddleware := func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set("custom", "value")
				return next(c)
			}
		}

		handler := func(c echo.Context) error {
			session := c.Get("session")
			custom := c.Get("custom")

			return c.JSON(http.StatusOK, echo.Map{
				"session": session != nil,
				"custom":  custom,
			})
		}

		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", token)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		// ミドルウェアをチェーン
		h := RequireAuth(customMiddleware(handler))
		err := h(c)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "value")
	})
}