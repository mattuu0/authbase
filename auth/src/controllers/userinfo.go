package controllers

import (
	"auth/services"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUserInfo(ctx echo.Context) error {
	raw := ctx.Request().Header.Get("Authorization")
	if raw == "" {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header is required"})
	}

	tokenString := raw
	if strings.HasPrefix(raw, "Bearer ") {
		tokenString = raw[7:]
	}

	info, err := services.ParseAccessToken(tokenString)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired token"})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"user_id":   info.UserID,
		"name":      info.Name,
		"email":     info.Email,
		"labels":    info.Labels,
		"prov_code": info.ProvCode,
		"prov_uid":  info.ProvUid,
		"exp":       info.Exp,
	})
}
