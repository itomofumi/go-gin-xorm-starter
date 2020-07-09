package service

import (
	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/itomofumi/go-gin-xorm-starter/repository"
)

// FruitsInterface defines fruits service interface.
type FruitsInterface interface {
	GetAll() ([]*model.Fruit, error)
	GetByID(fruitID uint64) (*model.Fruit, error)
	Create(body *model.FruitBody) (*model.Fruit, error)
	Update(fruitID uint64, notice *model.FruitBody) (*model.Fruit, error)
	Delete(fruitID uint64) error
}

// Fruits implements fruits service.
type Fruits struct {
	repo repository.FruitsInterface
}

// NewFruits initializes fruits service.
func NewFruits(repo repository.FruitsInterface) FruitsInterface {
	f := Fruits{repo}
	return &f
}

// GetAll returns all fruits.
func (f *Fruits) GetAll() ([]*model.Fruit, error) {
	return f.repo.GetAll()
}

// GetByID returns a fruit specified by the given id.
func (f *Fruits) GetByID(fruitID uint64) (*model.Fruit, error) {
	return f.repo.GetByID(fruitID)
}

// Create creates a new fruit.
func (f *Fruits) Create(body *model.FruitBody) (*model.Fruit, error) {
	return f.repo.Create(body)
}

// Update updates a fruit specified by the given id.
func (f *Fruits) Update(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
	return f.repo.Update(fruitID, body)
}

// Delete deletes a fruit specified by the given id.
func (f *Fruits) Delete(fruitID uint64) error {
	return f.repo.Delete(fruitID)
}
