package models

import (
	"auth/logger"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func OpenDB() (*gorm.DB,error) {
	dsn := os.Getenv("DB_DSN")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	// エラー処理
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigreteTable(db *gorm.DB) error {
	logger.Println("マイグレーションを実行しています...")

	// データベース接続確認
	err := db.AutoMigrate(&User{}, &Provider{}, &Session{}, &Label{}, &AdminUser{})	

	// エラー処理
	if err != nil {
		return err
	}

	logger.Println("マイグレーションを実行しました")

	return nil
}

func Init() error {
	logger.Println("データベース接続を確立しています...")

	// データベース接続
	db, err := OpenDB()
	if err != nil {
		logger.PrintErr("データベース接続に失敗しました", err)
		return err
	}

	// データベース接続確認
	err = MigreteTable(db)

	// エラー処理
	if err != nil {
		return err
	}

	// グローバル変数に格納
	dbconn = db

	// プロバイダを初期化する
	InitProviders()

	return nil
}

func GetDB() *gorm.DB {
	return dbconn
}