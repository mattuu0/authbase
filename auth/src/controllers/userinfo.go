package controllers

import (
	"auth/models"
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

	tokenString, ok := strings.CutPrefix(raw, "Bearer ")
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "Authorization header must use Bearer scheme"})
	}

	info, err := services.ParseAccessToken(tokenString)
	if err != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "invalid or expired token"})
	}

	// BANされたユーザーはアクセス不可
	user, result := models.GetUser(info.UserID)
	if result.Error != nil {
		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})
	}
	if user.IsBanned == 1 {
		return ctx.JSON(http.StatusForbidden, echo.Map{"error": "Your account has been banned"})
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
