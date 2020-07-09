package service

import (
	"fmt"

	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/itomofumi/go-gin-xorm-starter/repository"
)

// UsersInterface はサポーター管理サービスです
type UsersInterface interface {
	Create(email string, profile *model.UserProfile) (*model.UserPublicData, error)
	GetByID(uint64) (user *model.User, ok bool)
	GetByEmail(email string) (user *model.User, ok bool)
	Verify(userID uint64) error
	Update(uint64, *model.UserProfile) (*model.UserPublicData, error)
	Delete(uint64) error
}

// Users はサポーターのサービス実装
type Users struct {
	repo repository.UsersInterface
}

// NewUsers はユーザのサービスを初期化
func NewUsers(repo repository.UsersInterface) UsersInterface {
	u := Users{repo}
	return &u
}

// Create はユーザを登録
func (u *Users) Create(email string, profile *model.UserProfile) (*model.UserPublicData, error) {

	// Get user's verification status.
	currentUser, found := u.GetByEmail(email)

	if found {
		if currentUser.EmailVerified != nil && *currentUser.EmailVerified {
			return nil, fmt.Errorf("user is already verified")
		}
		// Delete temporary user.
		// Don't care wheather success or not.
		_ = u.Delete(currentUser.ID)
	}

	return u.repo.Create(email, profile)
}

// GetByID は指定のユーザを取得します
func (u *Users) GetByID(userID uint64) (user *model.User, ok bool) {
	return u.repo.GetByID(userID)
}

// GetByEmail は指定のユーザを取得します
func (u *Users) GetByEmail(email string) (user *model.User, ok bool) {
	return u.repo.GetByEmail(email)
}

// Verify はユーザーを認証済みにします
func (u *Users) Verify(userID uint64) error {
	return u.repo.Verify(userID)
}

// Update はユーザを削除
func (u *Users) Update(id uint64, user *model.UserProfile) (*model.UserPublicData, error) {
	return u.repo.Update(id, user)
}

// Delete はユーザを削除
func (u *Users) Delete(id uint64) error {
	return u.repo.Delete(id)
}
