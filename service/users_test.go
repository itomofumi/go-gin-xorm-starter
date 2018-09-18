package service_test

import (
	"reflect"
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
	"github.com/gemcook/go-gin-xorm-starter/service"
	"github.com/gemcook/ptr"
)

// UsersRepositoryMock is a mock for Users repository.
type UsersRepositoryMock struct {
	repository.UsersInterface
	FakeGetByEmail func(email string) (user *model.User, ok bool)
	FakeCreate     func(email string, profile *model.UserProfile) (*model.UserPublicData, error)
	FakeDelete     func(id uint64) error
	// FakeGetByID    func(id uint64) (user *model.User, ok bool)
	// FakeVerify     func(userID uint64) error
	// FakeUpdate     func(id uint64, profile *model.UserProfile) (*model.UserPublicData, error)
}

func (ur *UsersRepositoryMock) GetByEmail(email string) (user *model.User, ok bool) {
	return ur.FakeGetByEmail(email)
}

func (ur *UsersRepositoryMock) Create(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
	return ur.FakeCreate(email, profile)
}

func (ur *UsersRepositoryMock) Delete(id uint64) error {
	return ur.FakeDelete(id)
}

// func (ur *UsersRepositoryMock) GetByID(id uint64) (user *model.User, ok bool) {
// 	return ur.FakeGetByID(id)
// }

// func (ur *UsersRepositoryMock) Verify(userID uint64) error {
// 	return ur.FakeVerify(userID)
// }

// func (ur *UsersRepositoryMock) Update(id uint64, profile *model.UserProfile) (*model.UserPublicData, error) {
// 	return ur.FakeUpdate(id, profile)
// }

var testUsers = []*model.User{
	{
		Common: model.Common{ID: 1},
		Email:  "foo@example.com",
		UserPublicData: model.UserPublicData{
			UserProfile: model.UserProfile{DisplayName: ptr.String("foo")},
		},
	},
}

func TestUsers_Create(t *testing.T) {

	type fakes struct {
		getByEmail func(email string) (user *model.User, ok bool)
		create     func(email string, profile *model.UserProfile) (*model.UserPublicData, error)
		delete     func(id uint64) error
	}
	type args struct {
		email   string
		profile *model.UserProfile
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		want    *model.UserPublicData
		wantErr bool
	}{
		{
			"success for first creation",
			fakes{
				getByEmail: func(email string) (user *model.User, ok bool) {
					return nil, false
				},
				create: func(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
					return testUsers[0].GetPublicData(), nil
				},
				delete: func(id uint64) error { return nil },
			},
			args{
				email:   "foo@example.com",
				profile: &testUsers[0].UserProfile,
			},
			testUsers[0].GetPublicData(),
			false,
		},
		{
			"success for secondary creation",
			fakes{
				getByEmail: func(email string) (user *model.User, ok bool) {
					return &model.User{EmailVerified: ptr.Bool(false)}, true
				},
				create: func(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
					return testUsers[0].GetPublicData(), nil
				},
				delete: func(id uint64) error { return nil },
			},
			args{
				email:   "foo@example.com",
				profile: &testUsers[0].UserProfile,
			},
			testUsers[0].GetPublicData(),
			false,
		},
		{
			"failure for verified email",
			fakes{
				getByEmail: func(email string) (user *model.User, ok bool) {
					return &model.User{EmailVerified: ptr.Bool(true)}, true
				},
				create: func(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
					return nil, nil
				},
				delete: func(id uint64) error { return nil },
			},
			args{
				email:   "foo@example.com",
				profile: &testUsers[0].UserProfile,
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &UsersRepositoryMock{
				FakeCreate:     tt.fakes.create,
				FakeGetByEmail: tt.fakes.getByEmail,
				FakeDelete:     tt.fakes.delete,
			}
			u := service.NewUsers(repo)

			got, err := u.Create(tt.args.email, tt.args.profile)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}
