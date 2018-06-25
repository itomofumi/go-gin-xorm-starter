package repository

import (
	"fmt"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/util"
)

// UsersInterface has users data.
type UsersInterface interface {
	GetByEmail(email string) (user *model.User, ok bool)
	GetByID(id uint64) (user *model.User, ok bool)
	Create(email string, profile *model.UserProfile) (*model.UserPublicData, error)
	Verify(userID uint64) error
	Update(id uint64, profile *model.UserProfile) (*model.UserPublicData, error)
	Delete(id uint64) error
}

// Users has users data.
type Users struct {
	engine infra.EngineInterface
}

// NewUsers initializes Users
func NewUsers(engine infra.EngineInterface) UsersInterface {
	r := Users{
		engine: engine,
	}
	return &r
}

// GetByEmail returns an user who has the given email.
func (r *Users) GetByEmail(email string) (user *model.User, ok bool) {
	var u model.User
	ok, err := r.engine.Where("is_deleted = ? AND is_enabled = ? AND email = ?", false, true, email).Get(&u)
	if err != nil {
		return nil, false
	}
	if !ok {
		return nil, false
	}
	return &u, true
}

// GetByID はIDでユーザー情報を取得します
func (r *Users) GetByID(id uint64) (user *model.User, ok bool) {
	user = &model.User{}
	ok, err := r.engine.ID(id).Get(user)
	if err != nil || !ok {
		return nil, false
	}
	return user, true
}

// Create adds a new user.
func (r *Users) Create(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
	user := model.User{}
	user.Common.UnsetDefaltCols()
	user.Email = email
	user.DisplayName = profile.DisplayName
	user.About = profile.About
	user.AvatarURL = profile.AvatarURL
	user.LastLoginAt = util.GetTimeNow()

	session := r.engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return nil, err
	}
	var existsUser model.User
	found, err := session.Where("is_enabled = ? AND email = ? AND email_verified = ?", true, email, true).Get(&existsUser)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	if found {
		session.Rollback()
		return nil, fmt.Errorf("verified user exists")
	}

	_, err = session.InsertOne(&user)
	if err != nil {
		session.Rollback()
		return nil, err
	}
	err = session.Commit()
	if err != nil {
		return nil, err
	}

	return user.GetPublicData(), nil
}

// Verify updates user as verified
func (r *Users) Verify(userID uint64) error {
	sql := "UPDATE users SET identity_verified = ? WHERE id = ? AND is_deleted = ? AND is_enabled = ?"
	_, err := r.engine.Exec(sql, true, userID, false, true)
	if err != nil {
		return err
	}

	return nil
}

// Update updates user's profile data.
func (r *Users) Update(id uint64, profile *model.UserProfile) (*model.UserPublicData, error) {
	now := util.GetFormatedTimeNow()

	sql := "UPDATE users SET updated_at = ?, display_name = ?, about = ?, avatar_url = ? WHERE id = ?"

	_, err := r.engine.Exec(sql, now, profile.DisplayName, profile.About, profile.AvatarURL, id)
	if err != nil {
		return nil, err
	}

	var user model.User
	ok, err := r.engine.ID(id).Get(&user)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("target user not found")
	}
	return user.GetPublicData(), nil
}

// Delete sets is_deleted = true
func (r *Users) Delete(id uint64) error {
	now := util.GetFormatedTimeNow()

	sql := "UPDATE users SET is_deleted = ?, is_enabled = ?, updated_at = ? WHERE id = ?"

	_, err := r.engine.Exec(sql, true, false, now, id)
	if err != nil {
		return err
	}

	return nil
}
