// Package integration provides integration tests for complete authentication flows.
// Tests cover end-to-end scenarios including user creation, login, token usage, and logout.
package integration

import (
	"auth/controllers"
	"auth/middlewares"

	"github.com/labstack/echo/v4"
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
