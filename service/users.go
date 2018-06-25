package service

import (
	"fmt"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
)

// UsersInterface はサポーター管理サービスです
type UsersInterface interface {
	Create(*model.User) (*model.UserPublicData, error)
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
func (u *Users) Create(user *model.User) (*model.UserPublicData, error) {

	// 現在の登録状態を取得
	currentUser, found := u.GetByEmail(user.Email)

	if found {
		if currentUser.EmailVerified {
			return nil, fmt.Errorf("認証済みのユーザーのため、登録できません")
		}
		// 仮ユーザーなので、削除する
		err := u.Delete(currentUser.ID)
		if err != nil {
			// 失敗しても気にしない
		}
	}

	return u.repo.Create(user.Email, &user.UserProfile)
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
