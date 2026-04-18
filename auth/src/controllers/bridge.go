package controllers

import (
	"auth/models"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// IssueBridgeToken はログイン済みセッションからユーザーを取得し、アクセストークンを関連付けたブリッジトークンを発行します
func IssueBridgeToken(ctx echo.Context) error {
	session, ok := ctx.Get("session").(*models.Session)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}

	// 現在のアクセストークンをリクエストヘッダー等から取得
	accessToken := ctx.Request().Header.Get("Authorization")
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	token, err := services.IssueBridgeToken(session.UserID, accessToken)
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
