package model

// Fruit is a model
type Fruit struct {
	Common     `xorm:"extends"`
	*FruitBody `xorm:"extends"`
}

// FruitBody the main data
type FruitBody struct {
	Name  *string `json:"name"`
	Price *int    `json:"price"`
}

// TableName はテーブル名を返す
func (Fruit) TableName() string {
	return "fruits"
}
