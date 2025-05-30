package models

import (
	"log"
	"os"
	// "os"

	// "gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	// "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func Init() {
	// データベースを開く
	// db, err := gorm.Open(sqlite.Open(os.Getenv("DBPATH")), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	// データベースの接続情報
	dsn := os.Getenv("DATABASE_DSN")

	// データベースを開く
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// マイグレーション
	db.AutoMigrate(&FriendRequest{})
	db.AutoMigrate(&Friend{})

	// グローバル変数に格納
	dbconn = db
}
