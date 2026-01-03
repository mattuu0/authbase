package controllers

import (
	"auth/models"
	"auth/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Logout(ctx echo.Context) error {

	// セッションを取得

	session, ok := ctx.Get("session").(*models.Session)



	// エラー処理

	if !ok {

		return ctx.JSON(http.StatusUnauthorized, echo.Map{"error": "unauthorized"})

	}



	// セッションを削除

	if err := services.Logout(session); err != nil {

		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})

	}



	return ctx.JSON(http.StatusOK, echo.Map{"message": "success"})

}



func LoginPage(ctx echo.Context) error {



	isPopup := ctx.QueryParam("ispopup")



	isMobile := ctx.QueryParam("ismobile")







	return ctx.Render(http.StatusOK, "login.html", map[string]interface{}{



		"IsPopup":  isPopup,



		"IsMobile": isMobile,



	})



}















func OpenAppPage(ctx echo.Context) error {

	return ctx.Render(http.StatusOK, "open-app.html", map[string]interface{}{

		"CustomScheme": CUSTOM_SCHEME,

	})

}
