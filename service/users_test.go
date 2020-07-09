package service_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/itomofumi/go-gin-xorm-starter/repository"
	"github.com/itomofumi/go-gin-xorm-starter/service"
	"github.com/itomofumi/ptr"
)

// usersRepositoryMock is a mock for Users repository.
type usersRepositoryMock struct {
	repository.UsersInterface
	FakeGetByEmail func(email string) (user *model.User, ok bool)
	FakeCreate     func(email string, profile *model.UserProfile) (*model.UserPublicData, error)
	FakeDelete     func(id uint64) error
	FakeGetByID    func(id uint64) (user *model.User, ok bool)
	FakeVerify     func(userID uint64) error
	FakeUpdate     func(id uint64, profile *model.UserProfile) (*model.UserPublicData, error)
}

func (ur *usersRepositoryMock) GetByEmail(email string) (user *model.User, ok bool) {
	return ur.FakeGetByEmail(email)
}

func (ur *usersRepositoryMock) Create(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
	return ur.FakeCreate(email, profile)
}

func (ur *usersRepositoryMock) Delete(id uint64) error {
	return ur.FakeDelete(id)
}

func (ur *usersRepositoryMock) GetByID(id uint64) (user *model.User, ok bool) {
	return ur.FakeGetByID(id)
}

func (ur *usersRepositoryMock) Verify(userID uint64) error {
	return ur.FakeVerify(userID)
}

func (ur *usersRepositoryMock) Update(id uint64, profile *model.UserProfile) (*model.UserPublicData, error) {
	return ur.FakeUpdate(id, profile)
}

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
			repo := &usersRepositoryMock{
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

func TestUsers_Verify(t *testing.T) {
	type fakes struct {
		verify func(userID uint64) error
	}
	type args struct {
		userID uint64
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		wantErr bool
	}{
		{"success",
			fakes{
				verify: func(userID uint64) error {
					return nil
				},
			},
			args{userID: 1}, false,
		},
		{"failure",
			fakes{
				verify: func(userID uint64) error {
					return fmt.Errorf("cannot verify")
				},
			},
			args{userID: 1}, true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &usersRepositoryMock{
				FakeVerify: tt.fakes.verify,
			}
			u := service.NewUsers(repo)

			if err := u.Verify(tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Users.Verify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUsers_GetByID(t *testing.T) {
	type fakes struct {
		getByID func(userID uint64) (user *model.User, ok bool)
	}
	type args struct {
		userID uint64
	}
	tests := []struct {
		name     string
		fakes    fakes
		args     args
		wantUser *model.User
		wantOk   bool
	}{
		{"success",
			fakes{
				getByID: func(userID uint64) (user *model.User, ok bool) {
					return &model.User{
						Common: model.Common{ID: 1},
						Email:  "foo@example.com",
					}, true
				},
			},
			args{userID: 1},
			&model.User{
				Common: model.Common{ID: 1},
				Email:  "foo@example.com",
			},
			true,
		},
		{"failure",
			fakes{
				getByID: func(userID uint64) (user *model.User, ok bool) {
					return nil, false
				},
			},
			args{userID: 999},
			nil,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &usersRepositoryMock{
				FakeGetByID: tt.fakes.getByID,
			}
			u := service.NewUsers(repo)
			gotUser, gotOk := u.GetByID(tt.args.userID)
			if !reflect.DeepEqual(gotUser, tt.wantUser) {
				t.Errorf("Users.GetByID() gotUser = %v, want %v", gotUser, tt.wantUser)
			}
			if gotOk != tt.wantOk {
				t.Errorf("Users.GetByID() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestUsers_Update(t *testing.T) {
	type fakes struct {
		update func(id uint64, user *model.UserProfile) (*model.UserPublicData, error)
	}
	type args struct {
		id   uint64
		user *model.UserProfile
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		want    *model.UserPublicData
		wantErr bool
	}{
		{"success",
			fakes{
				update: func(id uint64, user *model.UserProfile) (*model.UserPublicData, error) {
					return &model.UserPublicData{
						UserID:      id,
						UserProfile: *user,
					}, nil
				},
			},
			args{id: 1, user: &model.UserProfile{DisplayName: ptr.String("foo")}},
			&model.UserPublicData{
				UserID:      1,
				UserProfile: model.UserProfile{DisplayName: ptr.String("foo")},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &usersRepositoryMock{
				FakeUpdate: tt.fakes.update,
			}
			u := service.NewUsers(repo)
			got, err := u.Update(tt.args.id, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("Users.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Users.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
