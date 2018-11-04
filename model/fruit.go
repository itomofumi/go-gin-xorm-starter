package model

import "github.com/go-playground/validator"

// Fruit is a model
type Fruit struct {
	Common    `xorm:"extends"`
	FruitBody `xorm:"extends"`
}

// FruitBody the main data
type FruitBody struct {
	Name  *string `json:"name" binding:"required,min=1"`
	Price *int    `json:"price"`
}

// TableName はテーブル名を返す
func (Fruit) TableName() string {
	return "fruits"
}

// FruitBodyStructLevelValidation contains FruitBody custom struct level validations.
func FruitBodyStructLevelValidation(sl validator.StructLevel) {

	fruitBody := sl.Current().Interface().(FruitBody)

	if fruitBody.Price == nil || *fruitBody.Price < 0 {
		sl.ReportError(fruitBody.Price, "Price", "price", "notminus", "")
	}
}
