package repository_test

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/repository"
)

func TestFruits_GetAll(t *testing.T) {
	engine, cleanup := Setup(t)
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
