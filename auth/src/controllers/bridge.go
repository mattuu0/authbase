package controllers

import (
	"auth/models"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// IssueBridgeToken はログイン済みセッションからユーザーを取得し、リフレッシュトークンを関連付けたブリッジトークンを発行します
func IssueBridgeToken(ctx echo.Context) error {
	session, ok := ctx.Get("session").(*models.Session)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	// ミドルウェアで認証されたリフレッシュトークン（ヘッダーのAuthorization）を取得
	refreshToken := ctx.Request().Header.Get("Authorization")
	if len(refreshToken) > 7 && refreshToken[:7] == "Bearer " {
		refreshToken = refreshToken[7:]
	}

	token, err := services.IssueBridgeToken(session.UserID, refreshToken)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to issue bridge token"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"bridge_token": token})
}

// ExchangeBridgeToken は Authorization ヘッダーから受け取ったJWT形式のブリッジトークンを検証し、アクセストークンとリフレッシュトークンを返却します
func ExchangeBridgeToken(ctx echo.Context) error {
	bridgeTokenString := ctx.Request().Header.Get("Authorization")
	if bridgeTokenString == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Authorization header is required"})
	}

	// "Bearer " プレフィックスがある場合は削除
	if len(bridgeTokenString) > 7 && bridgeTokenString[:7] == "Bearer " {
		bridgeTokenString = bridgeTokenString[7:]
	}

	tokens, err := services.ExchangeBridgeToken(bridgeTokenString)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"refresh_token": tokens["refresh_token"],
	})
}
