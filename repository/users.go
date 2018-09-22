package repository

import (
	"fmt"

	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/util"
	"github.com/gemcook/ptr"
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
func NewUsers(engine infra.EngineInterface) *Users {
	u := Users{
		engine: engine,
	}
	return &u
}

// GetByEmail returns an user who has the given email.
func (u *Users) GetByEmail(email string) (user *model.User, ok bool) {
	var result model.User
	ok, err := u.engine.Where(
		`
		is_deleted = ? 
		AND is_enabled = ? 
		AND email = ?
		`,
		false, true, email).Get(&result)

	if err != nil {
		return nil, false
	}
	if !ok {
		return nil, false
	}
	return &result, true
}

// GetByID はIDでユーザー情報を取得します
func (u *Users) GetByID(id uint64) (user *model.User, ok bool) {
	user = &model.User{}
	ok, err := u.engine.ID(id).Get(user)
	if err != nil || !ok {
		return nil, false
	}
	user.UserID = id
	return user, true
}

// Create adds a new user.
func (u *Users) Create(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
	user := model.User{}
	user.Common.SetDefault()
	user.Email = email
	user.EmailVerified = ptr.Bool(false)
	user.DisplayName = profile.DisplayName
	user.About = profile.About
	user.AvatarURL = profile.AvatarURL
	user.LastLoginAt = util.GetTimeNow()

	session := u.engine.NewSession()
	defer session.Close()

	err := session.Begin()
	if err != nil {
		return nil, err
	}
	var existsUser model.User
	found, err := session.Where(
		`
		is_enabled = ? 
		AND email = ? 
		AND email_verified = ?
		`, true, email, true).Get(&existsUser)

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
func (u *Users) Verify(userID uint64) error {
	user := model.User{}
	user.EmailVerified = ptr.Bool(true)
	_, err := u.engine.ID(userID).Where(
		`
		is_deleted = ? 
		AND is_enabled = ?
		`, false, true).Update(&user)
	if err != nil {
		return err
	}

	return nil
}

// Update updates user's profile data.
func (u *Users) Update(id uint64, profile *model.UserProfile) (*model.UserPublicData, error) {
	if profile == nil {
		return nil, fmt.Errorf("profile must not be nil")
	}
	user := model.User{}
	user.UserProfile = *profile

	_, err := u.engine.ID(id).Update(&user)
	if err != nil {
		return nil, err
	}

	var updated model.User
	_, err = u.engine.ID(id).Get(&updated)
	if err != nil {
		return nil, err
	}

	return updated.GetPublicData(), nil
}

// Delete sets is_deleted = true
func (u *Users) Delete(id uint64) error {
	user := model.User{}
	user.IsDeleted = ptr.Bool(true)

	_, err := u.engine.ID(id).Where("is_deleted = ?", false).Update(&user)
	if err != nil {
		return err
	}

	return nil
}
