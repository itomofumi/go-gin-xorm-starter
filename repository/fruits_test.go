package repository_test

import (
	"testing"

	"github.com/gemcook/ptr"
	"github.com/stretchr/testify/assert"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
)

func TestFruits_GetAll(t *testing.T) {
	engine, cleanup := setupDB(t)
	defer cleanup()

	fruits := repository.NewFruits(engine)
	list, err := fruits.GetAll()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 11 {
		t.Errorf("Fruits.GetAll() returned wrong number of result. got=%d, want=%d", len(list), 11)
	}
}

func TestFruits_GetByID(t *testing.T) {
	engine, cleanup := setupDB(t)
	defer cleanup()

	fruits := repository.NewFruits(engine)

	var id uint64 = 1
	result, err := fruits.GetByID(id)
	if err != nil {
		t.Errorf("Fruits.GetByID() returned an unexpected error=%v", err)
	}

	assert := assert.New(t)
	assert.Equal(id, result.ID)
	assert.Equal("Apple", *result.Name)
}

func TestFruits_Create(t *testing.T) {
	engine, cleanup := setupDB(t)
	defer cleanup()

	fruits := repository.NewFruits(engine)

	body := model.FruitBody{
		Name:  ptr.String("Lemon"),
		Price: ptr.Int(123),
	}
	result, err := fruits.Create(&body)
	if err != nil {
		t.Errorf("Fruits.Create() returned an unexpected error=%v", err)
	}
	assert := assert.New(t)
	assert.Equal("Lemon", *result.Name)
	assert.Equal(123, *result.Price)
}

func TestFruits_Update(t *testing.T) {
	engine, cleanup := setupDB(t)
	defer cleanup()

	fruits := repository.NewFruits(engine)

	var id uint64 = 1
	body := model.FruitBody{
		Price: ptr.Int(999),
	}
	result, err := fruits.Update(id, &body)
	if err != nil {
		t.Errorf("Fruits.Update() returned an unexpected error=%v", err)
	}

	assert := assert.New(t)
	assert.Equal(999, *result.Price)
}

func TestFruits_Delete(t *testing.T) {
	engine, cleanup := setupDB(t)
	defer cleanup()

	fruits := repository.NewFruits(engine)

	var id uint64 = 1
	err := fruits.Delete(1)

	result, err := fruits.GetByID(id)
	if err == nil {
		t.Errorf("Fruits.Delete() could not delete fruit. id=%d, got=%+v", id, result)
	}
}
