package repository

import (
	"fmt"

	"github.com/go-xorm/xorm"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/ptr"
)

// FruitsInterface is a fruits repository.
type FruitsInterface interface {
	GetAll() ([]*model.Fruit, error)
	GetByID(fruitID uint64) (*model.Fruit, error)
	Create(body *model.FruitBody) (*model.Fruit, error)
	Update(fruitID uint64, body *model.FruitBody) (*model.Fruit, error)
	Delete(fruitID uint64) error
}

// Fruits implements FruitsInterface.
type Fruits struct {
	engine xorm.EngineInterface
}

// NewFruits initializes a fruits repository.
func NewFruits(engine xorm.EngineInterface) *Fruits {
	f := Fruits{engine}
	return &f
}

// GetAll gets all fruits.
func (f *Fruits) GetAll() ([]*model.Fruit, error) {
	list := make([]*model.Fruit, 0)
	err := f.engine.Where("is_deleted = ?", false).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Create adds a new fruit and returns the created item.
func (f *Fruits) Create(body *model.FruitBody) (*model.Fruit, error) {
	fruit := model.Fruit{}
	if body != nil {
		fruit.FruitBody = *body
	}
	fruit.IsDeleted = ptr.Bool(false)
	fruit.IsEnabled = ptr.Bool(true)

	_, err := f.engine.InsertOne(&fruit)
	if err != nil {
		return nil, err
	}

	return &fruit, nil
}

// GetByID gets a fruit by the given ID.
func (f *Fruits) GetByID(fruitID uint64) (*model.Fruit, error) {
	fruit := model.Fruit{}

	found, err := f.engine.ID(fruitID).Where("is_deleted = ?", false).Get(&fruit)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("data not found for id = %v", fruitID)
	}
	return &fruit, nil
}

// Update edits a fruit data by the given ID.
func (f *Fruits) Update(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
	if body == nil {
		return nil, fmt.Errorf("body must not be nil")
	}
	fruit := model.Fruit{
		FruitBody: *body,
	}

	_, err := f.engine.ID(fruitID).Where("is_deleted = ?", false).Update(&fruit)

	if err != nil {
		return nil, err
	}

	// gets the updated item.
	updated, err := f.GetByID(fruitID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// Delete performs logical deletion by the given ID.
func (f *Fruits) Delete(fruitID uint64) error {
	fruit := model.Fruit{}
	fruit.IsDeleted = ptr.Bool(true)

	_, err := f.engine.ID(fruitID).Where("is_deleted = ?", false).Update(&fruit)
	if err != nil {
		return err
	}
	return nil
}
