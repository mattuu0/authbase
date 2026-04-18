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

	// 現在のリフレッシュトークン（セッション）をリクエストクッキー等から取得
	// ※既存のセッション管理方式に合わせて取得
	refreshToken, err := ctx.Cookie("sessionid") // 仮定：sessionid クッキーから取得
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "refresh token not found"})
	}

	token, err := services.IssueBridgeToken(session.UserID, refreshToken.Value)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": "failed to issue bridge token"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"bridge_token": token})
}

// ExchangeBridgeToken は Authorization ヘッダーから受け取ったJWT形式のブリッジトークンを検証し、アクセストークンを返却します
func ExchangeBridgeToken(ctx echo.Context) error {
	bridgeTokenString := ctx.Request().Header.Get("Authorization")
	if bridgeTokenString == "" {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": "Authorization header is required"})
	}

	// "Bearer " プレフィックスがある場合は削除
	if len(bridgeTokenString) > 7 && bridgeTokenString[:7] == "Bearer " {
		bridgeTokenString = bridgeTokenString[7:]
	}

	accessToken, err := services.ExchangeBridgeToken(bridgeTokenString)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"access_token": accessToken})
}
