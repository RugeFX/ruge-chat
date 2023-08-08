package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Model reference :
// 'reg_num', 'profile_picture', 'password', 'level_id', 'created_at',

type User struct {
	ID             uuid.UUID      `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `json:"deletedAt" gorm:"index"`
	Username       string         `json:"username" gorm:"not null;unique;type:varchar(25)"`
	Email          string         `json:"email" gorm:"not null;unique"`
	ProfilePicture string         `json:"profilePicture" gorm:"default:'default.png'"`
	Password       string         `json:"-" gorm:"not null;type:varchar(80)"`
	Level          int            `json:"level" gorm:"default:0"`
}

type ReqUser struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	ProfilePicture string `json:"profilePicture"`
}
