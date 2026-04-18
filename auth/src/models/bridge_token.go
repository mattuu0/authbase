package models

import (
	"time"

	"gorm.io/gorm"
)

type BridgeToken struct {
	gorm.Model
	Token       string    `gorm:"primaryKey;index"`
	UserID      string    `gorm:"index"`
	ExpiresAt   time.Time
	IsUsed      bool      `gorm:"default:false"`
	AccessToken string    `gorm:"type:text"`
}
