package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/itomofumi/go-gin-xorm-starter/handler"
	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/itomofumi/go-gin-xorm-starter/service"
	"github.com/itomofumi/ptr"
	"github.com/stretchr/testify/assert"
)

// UsersMock is a mock of users.
type UsersMock struct {
	service.UsersInterface
	FakeGetByEmail func(email string) (user *model.User, ok bool)
	FakeCreate     func(email string, profile *model.UserProfile) (*model.UserPublicData, error)
}

func (fm *UsersMock) GetByEmail(email string) (user *model.User, ok bool) {
	return fm.FakeGetByEmail(email)
}

func (fm *UsersMock) Create(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
	return fm.FakeCreate(email, profile)
}

var testUsers = []*model.User{
	{
		Common: model.Common{ID: 1},
		Email:  "foo@example.com",
		UserPublicData: model.UserPublicData{
			UserProfile: model.UserProfile{DisplayName: ptr.String("foo")},
		},
	},
	{
		Common: model.Common{ID: 2},
		Email:  "bar@example.com",
		UserPublicData: model.UserPublicData{
			UserProfile: model.UserProfile{DisplayName: ptr.String("bar")},
		},
	},
}

func TestGetMe(t *testing.T) {
	defer Setup()()

	type fakes struct {
		getByEmail func(email string) (*model.User, bool)
	}
	type args struct {
		email string
	}
	tests := []struct {
		name       string
		fakes      fakes
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				getByEmail: func(email string) (*model.User, bool) {
					for _, v := range testUsers {
						if v.Email == email {
							return v, true
						}
					}
					return nil, false
				},
			},
			args{email: "foo@example.com"},
			http.StatusOK,
			testUsers[0],
		},
		{"bad request",
			fakes{
				getByEmail: func(email string) (*model.User, bool) {
					return nil, false
				},
			},
			args{email: "xxx"},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("user not found")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Users := &UsersMock{
				FakeGetByEmail: tt.fakes.getByEmail,
			}
			factory := &ServiceFactoryMock{
				UsersMock: Users,
			}

			c, w := createGinTestContext(factory)
			c.Set("email", tt.args.email)
			handler.GetMe(c)
			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case *model.User:
				var res *model.User
				json.Unmarshal(w.Body.Bytes(), &res)
				assert.Equal(t, want.Email, res.Email)
				assert.Equal(t, want.DisplayName, res.DisplayName)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}

func TestPostUser(t *testing.T) {
	defer Setup()()
	type fakes struct {
		create func(email string, profile *model.UserProfile) (*model.UserPublicData, error)
	}
	type args struct {
		body interface{}
	}
	tests := []struct {
		name       string
		fakes      fakes
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				create: func(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
					return testUsers[0].GetPublicData(), nil
				},
			},
			args{
				body: model.UserCreateBody{
					Email:       "foo@example.com",
					UserProfile: testUsers[0].UserProfile,
				},
			},
			http.StatusCreated,
			testUsers[0].GetPublicData(),
		},
		{"bad request",
			fakes{
				create: func(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			args{
				body: model.UserCreateBody{
					Email:       "foo@example.com",
					UserProfile: testUsers[0].UserProfile,
				},
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
		{"wrong email format",
			fakes{
				create: func(email string, profile *model.UserProfile) (*model.UserPublicData, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			args{
				body: struct {
					Email string `json:"email"`
				}{Email: "aaa"},
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "request body mismatch", "Key: 'UserCreateBody.Email' Error:Field validation for 'Email' failed on the 'email' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			users := &UsersMock{
				FakeCreate: tt.fakes.create,
			}
			factory := &ServiceFactoryMock{
				UsersMock: users,
			}

			c, w := createGinTestContext(factory)
			b, _ := json.Marshal(tt.args.body)
			c.Request, _ = http.NewRequest("POST", "/users", bytes.NewBuffer(b))

			handler.PostUser(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case *model.UserPublicData:
				var res *model.UserPublicData
				json.Unmarshal(w.Body.Bytes(), &res)
				assert.Equal(t, want, res)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}
