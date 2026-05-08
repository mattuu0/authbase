package models

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbconn *gorm.DB = nil
)

func Init() {
	// データベースを開く
	
	// データベースの接続情報
	dbType := os.Getenv("DATABASE_TYPE")
	dsn := os.Getenv("DATABASE_DSN")

	var dialector gorm.Dialector

	switch dbType {
	case "postgres":
		dialector = postgres.Open(dsn)
	case "mysql":
		fallthrough
	default:
		dialector = mysql.Open(dsn)
	}

	// データベースを開く
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}

	// マイグレーション
	// db.AutoMigrate(&sample{})

	// グローバル変数に格納
	dbconn = db
}
