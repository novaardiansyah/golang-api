package models

import (
	"strings"
	"time"

	"golang-api/internal/config"

	"gorm.io/gorm"
)

type User struct {
	ID                   uint           `gorm:"primaryKey" json:"id"`
	Code                 string         `gorm:"size:255" json:"code"`
	Name                 string         `gorm:"size:255;not null" json:"name"`
	Email                string         `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Password             string         `gorm:"size:255;not null" json:"-"`
	HasAllowNotification *bool          `gorm:"default:false" json:"has_allow_notification"`
	NotificationToken    *string        `gorm:"size:255" json:"-"`
	AvatarUrl            *string        `gorm:"size:255" json:"-"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty" swaggertype:"string"`

	AvatarUrlFormatted *string `gorm:"-" json:"avatar_url"`
}

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	u.AvatarUrlFormatted = u.GetAvatarUrl()
	return
}

func (u User) GetAvatarUrl() *string {
	if u.AvatarUrl == nil || *u.AvatarUrl == "" {
		return u.AvatarUrl
	}
	if strings.HasPrefix(*u.AvatarUrl, "https") {
		return u.AvatarUrl
	}
	fullUrl := config.MainUrl + "/storage/" + *u.AvatarUrl
	return &fullUrl
}

func (User) TableName() string {
	return "users"
}
