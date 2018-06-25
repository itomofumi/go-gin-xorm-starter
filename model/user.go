package model

import (
	"time"
)

// User ユーザー情報を格納
type User struct {
	Common         `xorm:"extends"`
	Email          string     `xorm:"VARCHAR(120) notnull index(email)" json:"email"`
	EmailVerified  bool       `xorm:"notnull" json:"email_verified"`
	LastLoginAt    *time.Time `json:"last_login_at"`
	UserPublicData `xorm:"extends"`
}

// TableName represents db table name
func (User) TableName() string {
	return "users"
}

// GetPublicData は公開用のユーザー情報を取得
func (u *User) GetPublicData() *UserPublicData {
	pub := u.UserPublicData
	pub.PublicID = u.ID
	return &pub
}

// UserPublicData has public user data
type UserPublicData struct {
	PublicID    uint64 `xorm:"-" json:"id"`
	UserProfile `xorm:"extends"`
}

// UserProfile has user's editable profile data
type UserProfile struct {
	DisplayName *string `json:"display_name"`
	About       *string `json:"about"`
	AvatarURL   *string `json:"avatar_url"`
}

// UserCreateBody contains new user data.
type UserCreateBody struct {
	Email string `json:"email"`
	UserProfile
}

// TableName represents db table name
func (UserPublicData) TableName() string {
	return "users"
}
