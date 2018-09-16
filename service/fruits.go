package service

import (
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
)

// FruitsInterface はフルーツサービスのインタフェース
type FruitsInterface interface {
	GetAll() ([]*model.Fruit, error)
	GetByID(fruitID uint64) (*model.Fruit, error)
	Create(body *model.FruitBody) (*model.Fruit, error)
	Update(fruitID uint64, notice *model.FruitBody) (*model.Fruit, error)
	Delete(fruitID uint64) error
}

// Fruits はフルーツサービス
type Fruits struct {
	repo repository.FruitsInterface
}

// NewFruits はフルーツサービスの初期化
func NewFruits(repo repository.FruitsInterface) FruitsInterface {
	c := Fruits{repo}
	return &c
}

// GetAll はフルーツを全て取得します
func (n *Fruits) GetAll() ([]*model.Fruit, error) {
	return n.repo.GetAll()
}

// GetByID は指定のフルーツを取得します
func (n *Fruits) GetByID(fruitID uint64) (*model.Fruit, error) {
	return n.repo.GetByID(fruitID)
}

// Create はフルーツを新規追加します
func (n *Fruits) Create(body *model.FruitBody) (*model.Fruit, error) {
	return n.repo.Create(body)
}

// Update はフルーツを更新します
func (n *Fruits) Update(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
	return n.repo.Update(fruitID, body)
}

// Delete はフルーツを削除します
func (n *Fruits) Delete(fruitID uint64) error {
	return n.repo.Delete(fruitID)
}
