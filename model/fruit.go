package model

// Fruit is a model
type Fruit struct {
	Common    `xorm:"extends"`
	FruitBody `xorm:"extends"`
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

// IsValid checks fruit data.
func (f *FruitBody) IsValid() bool {
	if f.Name == nil || *f.Name == "" {
		return false
	}
	if f.Price == nil || *f.Price < 0 {
		return false
	}
	return true
}
