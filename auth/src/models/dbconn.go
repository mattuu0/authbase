package models

import (
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

func Init() error {
	// データベース接続
	db, err := OpenDB()
	if err != nil {
		return err
	}

	// データベース接続確認
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Provider{})
	db.AutoMigrate(&Session{})
	db.AutoMigrate(&Label{})
	db.AutoMigrate(&AdminUser{})

	// グローバル変数に格納
	dbconn = db

	// プロバイダを初期化する
	InitProviders()

	return nil
}

func GetDB() *gorm.DB {
	return dbconn
}