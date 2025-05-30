package controllers

import (
	"net/http"
	"app/services"
	"app/utils"

	"github.com/labstack/echo/v4"
)

type SearchArgs struct {
	UserName string `json:"name"`
}

func SearchUser(ctx echo.Context) error {
	var args SearchArgs

	// リクエストをバインド
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// バリデーション
	if args.UserName == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// 検索する
	result,err := services.SearchByName(args.UserName)

	// エラー処理
	if err != nil {
		utils.Println("failed to find user : " + err.Error())
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.JSON(http.StatusOK,result)
}

type RequestArgs struct {
	UserID string `json:"userid"`
	RequestID string `json:"requestid"`
}

// フレンドリクエストを送信
func FriendRequest(ctx echo.Context) error {
	var args RequestArgs

	// リクエストをバインド
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// バリデーション
	if args.UserID == "" {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// ユーザー情報を取得 (送信者)
	myid := ctx.Get("UserID").(string)

	if args.UserID == myid {
		return ctx.NoContent(http.StatusBadRequest)
	}

	// リクエストを作成する
	err := services.SendFriendRequest(myid,args.UserID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}

func GetSentRequest(ctx echo.Context) error {
	// ユーザー情報を取得 (送信者)
	myid := ctx.Get("UserID").(string)

	// リクエストを取得
	requests,err := services.GetSentRequest(myid)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK,requests)
}

// 受信済みを取得する
func RecvedRequest(ctx echo.Context) error {
	// ユーザー情報を取得 (受信者)
	myid := ctx.Get("UserID").(string)

	// リクエストを取得
	requests,err := services.GetRecvedRequest(myid)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK,requests)
}

func AcceptRequest(ctx echo.Context) error {
	// ユーザー情報を取得 (承認者)
	myid := ctx.Get("UserID").(string)

	// bind
	var args RequestArgs
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// リクエストを取得
	err := services.AcceptRequest(args.RequestID,myid)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}

func RejectRequest(ctx echo.Context) error {
	// bind
	var args RequestArgs
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// ユーザー情報を取得 (拒否者)
	myid := ctx.Get("UserID").(string)

	// リクエストを取得
	err := services.RejectRequest(myid,args.RequestID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusInternalServerError)
	}

	return ctx.NoContent(http.StatusOK)
}

func GetFriendList(ctx echo.Context) error {
	// ユーザー情報を取得
	myid := ctx.Get("UserID").(string)

	// フレンドリストを取得
	friends,err := services.GetFriendList(myid)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.JSON(http.StatusOK,friends)
}

type RemoveArgs struct {
	UserID string `json:"userid"`
}
func RemoveFriend(ctx echo.Context) error {
	// bind
	var args RemoveArgs
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// ユーザー情報を取得
	myid := ctx.Get("UserID").(string)

	// リクエストを取得
	err := services.RemoveFriend(myid,args.UserID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}

func CancelRequest(ctx echo.Context) error {
	// bind
	var args RequestArgs
	if err := ctx.Bind(&args); err != nil {
		utils.Println("failed to bind json : " + err.Error())
		return ctx.NoContent(http.StatusBadRequest)
	}

	// ユーザー情報を取得
	myid := ctx.Get("UserID").(string)

	// リクエストを取得
	err := services.CancelRequest(myid,args.RequestID)

	// エラー処理
	if err != nil {
		utils.Println(err)
		return ctx.NoContent(http.StatusConflict)
	}

	return ctx.NoContent(http.StatusOK)
}