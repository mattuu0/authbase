package services

import (
	"errors"
	"app/grpckit"
	"app/models"
	rediscache "app/redis-cache"
	"app/utils"
)

type SearchResult struct {
	UserID string `json:"userid"`
	Name   string `json:"name"`
}

func SearchByName(name string) ([]SearchResult,error) {
	utils.Println(name)

	// 検索する
	searchResult,err := grpckit.SearchUser("",name)

	// エラー処理
	if err != nil {
		return []SearchResult{},err
	}

	// 結果を変換
	result := []SearchResult{}
	for _, val := range searchResult.Users {
		// 結果に追加
		result = append(result, SearchResult{
			UserID: val.UserID,
			Name:   val.Name,
		})
	}

	return result,nil
}

func SendFriendRequest(myId,targetId string) error {
	// 自分が関わるリクエスト取得
	_, err := models.GetFriend(models.FriendArgs{
		UserID:   myId,
		TargetID: targetId,
	})

	// エラー処理
	if err == nil {
		return errors.New("Already friends")
	}

	// リクエストを作成
	return models.CreateRequest(models.FriendRequestArgs{
		UserID:   myId,
		TargetID: targetId,
	})
}

type FriendRequest struct {
	RequestID string `json:"id"`
	SenderID string	`json:"sender"`
	TargetID string `json:"target"`
}

func GetSentRequest(userid string) ([]FriendRequest,error) {
	// 送信済みを取得
	GetRequests,err := models.GetSentRequest(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{},err
	}

	// 返すリクエスト
	requests := []FriendRequest{}

	for _, request := range GetRequests {
		// リクエストを変換
		requests = append(requests, FriendRequest{
			RequestID: request.RequestID,
			SenderID:  request.SenderID,
			TargetID:  request.TargetID,
		})
	}

	return requests,nil
}

// 受信済みを取得
func GetRecvedRequest(userid string) ([]FriendRequest,error) {
	// 送信済みを取得
	GetRequests,err := models.GetRecvedRequest(userid)

	// エラー処理
	if err != nil {
		return []FriendRequest{},err
	}

	// 返すリクエスト
	requests := []FriendRequest{}

	for _, request := range GetRequests {
		// リクエストを変換
		requests = append(requests, FriendRequest{
			RequestID: request.RequestID,
			SenderID:  request.SenderID,
			TargetID:  request.TargetID,
		})
	}

	return requests,nil
}

// 承認する
func AcceptRequest(requestId string,myId string) error {
	// リクエストを取得
	request, err := models.GetRequestByID(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// ターゲットを検証する
	if request.TargetID != myId {
		return errors.New("invalid request")
	}

	// リクエストを承認
	err = models.AcceptRequest(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュ更新 (自分)
	err = CacheFriend(myId)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// キャッシュ更新 (相手)
	err = CacheFriend(request.SenderID)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// リクエストを削除
	return models.RemoveRequest(requestId)
}

func RejectRequest(myid,requestId string) error {
	// リクエストを取得
	request, err := models.GetRequestByID(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// ターゲットを検証する
	if request.TargetID != myid {
		return errors.New("invalid request")
	}

	// リクエストを削除
	return models.RemoveRequest(requestId)
}

func GetFriendList(userid string) ([]string,error) {
	// キャッシュから取得
	cached,err := rediscache.GetCacheFriend(userid)

	// 成功したとき
	if err == nil {
		return cached.FriendIds,nil
	}

	// フレンドリストを取得
	friends,err := models.GetFriendList(userid)

	// エラー処理
	if err != nil {
		return []string{},err
	}

	// キャッシュに保存
	err = CacheFriend(userid)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	return friends,nil
}

// フレンド削除する関数
func RemoveFriend(userid string,targetid string) error {
	// フレンド取得
	friend, err := models.GetFriend(models.FriendArgs{
		UserID:   userid,
		TargetID: targetid,
	})

	// エラー処理
	if err != nil {
		return err
	}

	// フレンドを削除
	err = models.RemoveFriend(friend.FriendID)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュ更新 (自分)
	err = CacheFriend(userid)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// キャッシュ更新 (相手)
	err = CacheFriend(targetid)

	// エラー処理
	if err != nil {
		// キャッシュの保存に失敗した場合
		utils.Println(err)
	}

	// フレンドを削除
	return nil
}

// ユーザーのフレンドをキャッシュする
func CacheFriend(userid string) error {
	// フレンドリストを取得
	friends,err := models.GetFriendList(userid)

	// エラー処理
	if err != nil {
		return err
	}

	// キャッシュ更新
	utils.Println("キャッシュ更新")
	utils.Println(userid)
	utils.Println(friends)

	// キャッシュに保存
	err = rediscache.AddCacheFriend(rediscache.CacheFriendArgs{
		UserID:   userid,
		Data: rediscache.FriendCache{
			FriendIds: friends,
		},
	})

	// エラー処理
	if err != nil {
		return err
	}

	return nil
}

// 送信済みリクエストをキャンセルする
func CancelRequest(myid,requestId string) error {
	// リクエストを取得
	request, err := models.GetRequestByID(requestId)

	// エラー処理
	if err != nil {
		return err
	}

	// 自分が送ったリクエストか
	if request.SenderID != myid {
		return errors.New("invalid request")
	}

	// リクエストを削除
	return models.RemoveRequest(requestId)
}