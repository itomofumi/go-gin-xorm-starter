package repository

import (
	"fmt"

	"github.com/go-xorm/xorm"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/util"
)

// FruitsInterface is a fruits repository.
type FruitsInterface interface {
	GetAll() ([]*model.Fruit, error)
	GetByID(fruitID int64) (*model.Fruit, error)
	Create(fruit *model.FruitBody) (*model.Fruit, error)
	Update(fruitID int64, fruit *model.FruitBody) (*model.Fruit, error)
	Delete(fruitID int64) error
}

// Fruits implements FruitsInterface.
type Fruits struct {
	engine xorm.EngineInterface
}

// NewFruits initializes a fruits repository.
func NewFruits(engine xorm.EngineInterface) FruitsInterface {
	s := Fruits{engine}
	return &s
}

// GetAll gets all fruits.
func (n *Fruits) GetAll() ([]*model.Fruit, error) {
	list := make([]*model.Fruit, 0)
	err := n.engine.Where("is_deleted = ?", false).Find(&list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

// Create adds a new fruit and returns the created item.
func (n *Fruits) Create(fruit *model.FruitBody) (*model.Fruit, error) {
	data := model.Fruit{}
	data.FruitBody = fruit
	data.IsDeleted = false
	data.IsEnabled = true

	_, err := n.engine.InsertOne(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

// GetByID gets a fruit by the given ID.
func (n *Fruits) GetByID(fruitID int64) (*model.Fruit, error) {
	data := model.Fruit{}

	found, err := n.engine.ID(fruitID).Where("is_deleted = ?", false).Get(&data)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("data not found for id = %v", fruitID)
	}
	return &data, nil
}

// Update edits a fruit data by the given ID.
func (n *Fruits) Update(fruitID int64, fruit *model.FruitBody) (*model.Fruit, error) {
	now := util.GetFormatedTimeNow()

	sql := `
	UPDATE
		fruits
	SET
		updated_at = ?,
		name = ?,
		price = ?
	WHERE
		id = ?
		AND is_deleted = ?
	`

	_, err := n.engine.Exec(sql, now, fruit.Name, fruit.Price, fruitID, false)
	if err != nil {
		return nil, err
	}

	// gets the updated item.
	updated, err := n.GetByID(fruitID)
	if err != nil {
		return nil, err
	}

	return updated, nil
}

// Delete performs logical deletion by the given ID.
func (n *Fruits) Delete(fruitID int64) error {
	now := util.GetFormatedTimeNow()

	sql := `
	UPDATE
		fruits
	SET
		updated_at = ?,
		is_deleted = ?
	WHERE
		id = ?
		AND is_deleted = ?
	`

	_, err := n.engine.Exec(sql, now, true, fruitID, false)
	if err != nil {
		return err
	}
	return nil
}
