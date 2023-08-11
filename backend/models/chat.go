package models

import (
	"time"

	"gorm.io/gorm"
)

// Model reference :
// 'reg_num', 'profile_picture', 'password', 'level_id', 'created_at',

type Chat struct {
	ID             int            `json:"id" gorm:"primaryKey"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Username       string         `json:"username" gorm:"not null;unique;type:varchar(25)"`
	Email          string         `json:"email" gorm:"not null;unique"`
	ProfilePicture string         `json:"profilePicture" gorm:"default:'default.png'"`
	Password       string         `json:"-" gorm:"not null;type:varchar(80)"`
	Level          int            `json:"level" gorm:"default:0"`
}
