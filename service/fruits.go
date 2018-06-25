package service

import (
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
)

// FruitsInterface はフルーツサービスのインタフェース
type FruitsInterface interface {
	GetAll() ([]*model.Fruit, error)
	GetByID(noticeID int64) (*model.Fruit, error)
	Create(notice *model.FruitBody) (*model.Fruit, error)
	Update(noticeID int64, notice *model.FruitBody) (*model.Fruit, error)
	Delete(noticeID int64) error
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
func (n *Fruits) GetByID(noticeID int64) (*model.Fruit, error) {
	return n.repo.GetByID(noticeID)
}

// Create はフルーツを新規追加します
func (n *Fruits) Create(notice *model.FruitBody) (*model.Fruit, error) {
	return n.repo.Create(notice)
}

// Update はフルーツを更新します
func (n *Fruits) Update(noticeID int64, notice *model.FruitBody) (*model.Fruit, error) {
	return n.repo.Update(noticeID, notice)
}

// Delete はフルーツを削除します
func (n *Fruits) Delete(noticeID int64) error {
	return n.repo.Delete(noticeID)
}
