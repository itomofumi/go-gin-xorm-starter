package repository_test

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
	"github.com/stretchr/testify/assert"
)

func TestUsers_GetByEmail(t *testing.T) {
	engine, cleanup := Setup(t)
	defer cleanup()

	users := repository.NewUsers(engine)

	email := "test@example.com"
	result, ok := users.GetByEmail(email)
	if !ok {
		t.Fatalf("Users.GetByEmail() could not get user by email = %s", email)
	}

	assert := assert.New(t)
	assert.EqualValues(email, result.Email)
}

func TestUsers_Create(t *testing.T) {
	engine, cleanup := Setup(t)
	defer cleanup()

	users := repository.NewUsers(engine)

	name := "foobar"
	email := "foobar@example.com"
	body := model.UserProfile{
		DisplayName: &name,
	}
	result, err := users.Create(email, &body)
	if err != nil {
		t.Fatalf("Users.Create() returned an unexpected error=%v", err)
	}

	assert := assert.New(t)
	assert.EqualValues(name, *result.DisplayName)
}

func TestUsers_Update(t *testing.T) {
	engine, cleanup := Setup(t)
	defer cleanup()

	users := repository.NewUsers(engine)

	var id uint64 = 1
	name := "foobar"
	body := model.UserProfile{
		DisplayName: &name,
	}
	result, err := users.Update(id, &body)
	if err != nil {
		t.Fatalf("Users.Update() returned an unexpected error=%v", err)
	}

	assert := assert.New(t)
	assert.EqualValues(name, *result.DisplayName)
}

func TestUsers_Delete(t *testing.T) {
	engine, cleanup := Setup(t)
	defer cleanup()

	users := repository.NewUsers(engine)

	var id uint64 = 1
	email := "test@example.com"
	err := users.Delete(id)
	if err != nil {
		t.Fatalf("Users.Delete() returned an unexpected error=%v", err)
	}

	_, ok := users.GetByEmail(email)
	if ok {
		t.Fatalf("Users.Delete() could not delete user by email = %s", email)
	}
}
