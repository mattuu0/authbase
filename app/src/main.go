package main

import (
	"net/http"
	"app/controllers"
	"app/grpckit"
	"app/middlewares"
	"app/models"
	rediscache "app/redis-cache"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// モデル初期化
	models.Init()

	// redis 初期化
	rediscache.Init()

	// grpc 初期化
	grpckit.Init()

	// ミドルウェア初期化
	middlewares.Init()

	// ルーター
	router := echo.New()

	// router.Use(middleware.Recover())
	router.Use(middleware.Logger())
	// router.Use(middlewares.PocketAuth())

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, echo.Map{
			"result": "hello world",
		})
	}, middlewares.RequireAuth)

	// フレンドグループ
	friendg := router.Group("/friend")
	{
		// ミドルウェア設定
		friendg.Use(middlewares.RequireAuth)

		// 検索するエンドポイント
		friendg.POST("/search",controllers.SearchUser)

		// フレンド一覧取得
		friendg.GET("/list",controllers.GetFriendList)

		// リクエスト(送信)
		friendg.POST("/request",controllers.FriendRequest)

		// 送信済みリクエスト 取得
		friendg.GET("/sentrequest",controllers.GetSentRequest)

		// 受信済みリクエスト 取得
		friendg.GET("/recvrequest",controllers.RecvedRequest)

		// リクエストを承認する
		friendg.POST("/accept",controllers.AcceptRequest)

		// リクエストを拒否する
		friendg.POST("/reject",controllers.RejectRequest)

		// フレンド削除
		friendg.POST("/remove",controllers.RemoveFriend)

		// リクエストをキャンセル
		friendg.POST("/cancel",controllers.CancelRequest)
	}
	
	// websocket 用
	// router.GET("/ws", websocket.HandleWs)

	router.Logger.Fatal(router.Start(":8090"))
}