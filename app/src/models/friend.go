package models

import (
	"errors"
	"app/utils"

	"gorm.io/gorm"
)

type FriendRequest struct {
	SenderID  string `gorm:"primaryKey"` //送信者ID
	TargetID  string `gorm:"primaryKey"` //相手のID
	RequestID string
	Created   int64 `gorm:"autoCreateTime"` // 作成時間としてUNIX秒を使用する
}

type Friend struct {
	FriendID string `gorm:"unique"`         //フレンドID
	SenderID string `gorm:"primaryKey"`     //送信者ID
	TargetID string `gorm:"primaryKey"`     //相手のID
	Created  int64  `gorm:"autoCreateTime"` // 作成時間としてUNIX秒を使用する
}

type FriendRequestArgs struct {
	UserID   string //自分のID
	TargetID string //相手のID
}

type FriendArgs struct {
	UserID   string //自分のID
	TargetID string //相手のID
}

func CreateRequest(args FriendRequestArgs) error {
	// リクエストが存在するとき
	_, err := GetRequest(args)

	// エラー処理
	if err == gorm.ErrRecordNotFound {
		// レコードが存在しない場合

		// IDを生成
		RequestID, err := utils.Genid()
		// エラー処理
		if err != nil {
			return err
		}

		// データベースに保存
		result := dbconn.Save(&FriendRequest{
			SenderID:  args.UserID,
			TargetID:  args.TargetID,
			RequestID: RequestID,
		})

		// エラー処理
		if result.Error != nil {
			return result.Error
		}

		return nil
	} else if err == nil {
		return errors.New("Already sent")
	}

	return err
}

// フレンドかどうかを取得
func GetFriend(args FriendArgs) (Friend, error) {
	friend := Friend{}

	// データベースから取得
	result := dbconn.Where(&Friend{
		SenderID: args.UserID,
		TargetID: args.TargetID,
	}).Or(&Friend{
		SenderID: args.TargetID,
		TargetID: args.UserID,
	}).First(&friend)

	// エラー処理
	if result.Error != nil {
		return Friend{}, result.Error
	}

	return friend, nil
}

// 自分が関わるフレンドを取得
func GetRequest(args FriendRequestArgs) (FriendRequest, error) {
	request := FriendRequest{}

	// データベースから取得
	result := dbconn.Where(&FriendRequest{
		SenderID: args.UserID,
		TargetID: args.TargetID,
	}).Or(&FriendRequest{
		SenderID: args.TargetID,
		TargetID: args.UserID,
	}).First(&request)

	// エラー処理
	if result.Error != nil {
		return FriendRequest{}, result.Error
	}

	return request, nil
}

// 送信済みのリクエストを取得
func GetSentRequest(senderID string) ([]FriendRequest, error) {
	// フレンドリクエスト
	requests := []FriendRequest{}

	// リクエスト取得
	result := dbconn.Where(&FriendRequest{
		SenderID: senderID,
	}).Find(&requests)

	// エラー処理
	if result.Error != nil {
		return []FriendRequest{}, result.Error
	}

	return requests, nil
}

// 受信済みのリクエストを取得
func GetRecvedRequest(recverID string) ([]FriendRequest, error) {
	// フレンドリクエスト
	requests := []FriendRequest{}

	// リクエスト取得
	result := dbconn.Where(&FriendRequest{
		TargetID: recverID,
	}).Find(&requests)

	// エラー処理
	if result.Error != nil {
		return []FriendRequest{}, result.Error
	}

	return requests, nil
}

// リクエストを取得
func GetRequestByID(requestID string) (FriendRequest, error) {
	request := FriendRequest{}

	// データベースから取得
	result := dbconn.Where(&FriendRequest{
		RequestID: requestID,
	}).First(&request)

	// エラー処理
	if result.Error != nil {
		return FriendRequest{}, result.Error
	}

	return request, nil
}

func AcceptRequest(requestID string) error {
	// リクエストを取得
	request, err := GetRequestByID(requestID)

	// エラー処理
	if err != nil {
		return err
	}

	// フレンドIDを生成
	FriendID, err := utils.Genid()
	// エラー処理
	if err != nil {
		return err
	}

	// リクエストを承認
	return dbconn.Save(&Friend{
		FriendID: FriendID,
		SenderID: request.SenderID,
		TargetID: request.TargetID,
	}).Error
}

// リクエストを削除
func RemoveRequest(requestID string) error {
	return dbconn.Where(&FriendRequest{
		RequestID: requestID,
	}).Unscoped().Delete(&FriendRequest{}).Error
}

// フレンド一覧
func GetFriendList(userid string) ([]string, error) {
	// フレンドリストを取得
	friends := []Friend{}

	// フレンドリストを取得
	result := dbconn.Where(&Friend{
		SenderID: userid,
	}).Or(&Friend{
		TargetID: userid,
	}).Find(&friends)

	// エラー処理
	if result.Error != nil {
		return []string{}, result.Error
	}

	// 返すリスト
	return_friends := []string{}
	for _, friend := range friends {
		if friend.SenderID == userid {
			return_friends = append(return_friends, friend.TargetID)
		} else {
			return_friends = append(return_friends, friend.SenderID)
		}
	}

	return return_friends, nil
}

func RemoveFriend(friendid string) error {
	return dbconn.Where(&Friend{
		FriendID: friendid,
	}).Unscoped().Delete(&Friend{}).Error
}