package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/controller"
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/service"
	"github.com/stretchr/testify/assert"
)

// FruitsMock is a mock of fruits.
type FruitsMock struct {
	service.FruitsInterface
	FakeGetAll  func() ([]*model.Fruit, error)
	FakeGetByID func(fruitID uint64) (*model.Fruit, error)
	FakeCreate  func(fruit *model.FruitBody) (*model.Fruit, error)
	FakeUpdate  func(fruitID uint64, fruit *model.FruitBody) (*model.Fruit, error)
	FakeDelete  func(fruitID uint64) error
}

func (fm *FruitsMock) GetAll() ([]*model.Fruit, error) {
	return fm.FakeGetAll()
}

func (fm *FruitsMock) GetByID(fruitID uint64) (*model.Fruit, error) {
	return fm.FakeGetByID(fruitID)
}

func (fm *FruitsMock) Create(fruit *model.FruitBody) (*model.Fruit, error) {
	return fm.FakeCreate(fruit)
}

func (fm *FruitsMock) Update(fruitID uint64, fruit *model.FruitBody) (*model.Fruit, error) {
	return fm.FakeUpdate(fruitID, fruit)
}

func (fm *FruitsMock) Delete(fruitID uint64) error {
	return fm.FakeDelete(fruitID)
}

var testData = []*model.Fruit{
	{
		Common: model.Common{ID: 1},
		FruitBody: &model.FruitBody{
			Price: 100,
			Name:  "Apple",
		},
	},
	{
		Common: model.Common{ID: 2},
		FruitBody: &model.FruitBody{
			Price: 200,
			Name:  "Mango",
		},
	},
}

func TestGetFruits(t *testing.T) {
	defer Setup()()

	type args struct {
		getAll func() ([]*model.Fruit, error)
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success", args{
			getAll: func() ([]*model.Fruit, error) {
				return testData, nil
			}},
			http.StatusOK,
			testData,
		},
		{"bad request", args{
			getAll: func() ([]*model.Fruit, error) {
				return nil, fmt.Errorf("some error")
			}},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeGetAll: tt.args.getAll,
			}
			registry := &RegistryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(registry)
			controller.GetFruits(c)
			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case []*model.Fruit:
				var res []*model.Fruit
				json.Unmarshal(w.Body.Bytes(), &res)
				assert.Equal(t, want, res)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}

func TestGetFruitByID(t *testing.T) {
	defer Setup()()

	type args struct {
		id      uint64
		getByID func(id uint64) (*model.Fruit, error)
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success", args{
			id: 1,
			getByID: func(id uint64) (*model.Fruit, error) {
				for _, v := range testData {
					if v.ID == id {
						return v, nil
					}
				}
				return nil, fmt.Errorf("data not found for id = %v", id)
			}},
			http.StatusOK,
			testData[0],
		},
		{"bad request", args{
			id: 9999,
			getByID: func(id uint64) (*model.Fruit, error) {
				return nil, fmt.Errorf("data not found for id = %v", id)
			}},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("data not found for id = %v", 9999)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeGetByID: tt.args.getByID,
			}
			registry := &RegistryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(registry)
			c.Set("fruit-id", tt.args.id)
			controller.GetFruitByID(c)
			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case *model.Fruit:
				var res *model.Fruit
				json.Unmarshal(w.Body.Bytes(), &res)
				assert.Equal(t, want, res)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}

func TestPostFruit(t *testing.T) {
	defer Setup()()

	type args struct {
		body   interface{}
		create func(fruit *model.FruitBody) (*model.Fruit, error)
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success", args{
			body: testData[0].FruitBody,
			create: func(fruit *model.FruitBody) (*model.Fruit, error) {
				return testData[0], nil
			}},
			http.StatusCreated,
			testData[0],
		},
		{"bad request", args{
			body: testData[0].FruitBody,
			create: func(fruit *model.FruitBody) (*model.Fruit, error) {
				return nil, fmt.Errorf("some error")
			}},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
		{"wrong body", args{
			body: struct {
				Price string `json:"price"`
			}{Price: "aaa"},
			create: func(fruit *model.FruitBody) (*model.Fruit, error) {
				return nil, fmt.Errorf("some error")
			}},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeCreate: tt.args.create,
			}
			registry := &RegistryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(registry)
			b, _ := json.Marshal(tt.args.body)
			c.Request, _ = http.NewRequest("POST", "/fruits", bytes.NewBuffer(b))

			controller.PostFruit(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case *model.Fruit:
				var res *model.Fruit
				json.Unmarshal(w.Body.Bytes(), &res)
				assert.Equal(t, want, res)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}

func TestPutFruit(t *testing.T) {
	defer Setup()()

	type args struct {
		id     uint64
		body   *model.FruitBody
		update func(fruitID uint64, fruit *model.FruitBody) (*model.Fruit, error)
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success", args{
			id:   1,
			body: testData[0].FruitBody,
			update: func(fruitID uint64, fruit *model.FruitBody) (*model.Fruit, error) {
				d := testData[0]
				d.FruitBody = fruit
				return d, nil
			}},
			http.StatusOK,
			testData[0],
		},
		{"bad request", args{
			id:   1,
			body: testData[0].FruitBody,
			update: func(fruitID uint64, fruit *model.FruitBody) (*model.Fruit, error) {
				return nil, fmt.Errorf("some error")
			}},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeUpdate: tt.args.update,
			}
			registry := &RegistryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(registry)
			b, _ := json.Marshal(tt.args.body)
			c.Request, _ = http.NewRequest("PUT", "/fruits/:fruit-id", bytes.NewBuffer(b))
			c.Set("fruit-id", tt.args.id)

			controller.PutFruit(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case *model.Fruit:
				var res *model.Fruit
				json.Unmarshal(w.Body.Bytes(), &res)
				assert.Equal(t, want, res)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}

func TestDeleteFruit(t *testing.T) {
	defer Setup()()

	type args struct {
		id     uint64
		delete func(fruitID uint64) error
	}
	tests := []struct {
		name       string
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success", args{
			id:     1,
			delete: func(fruitID uint64) error { return nil },
		},
			http.StatusNoContent,
			nil,
		},
		{"bad request", args{
			id: 9999,
			delete: func(fruitID uint64) error {
				return fmt.Errorf("data not found for id = %v", fruitID)
			}},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("data not found for id = %v", 9999)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeDelete: tt.args.delete,
			}
			registry := &RegistryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(registry)
			c.Set("fruit-id", tt.args.id)

			controller.DeleteFruit(c)

			assert.Equal(t, tt.wantStatus, w.Code)

			switch want := tt.want.(type) {
			case nil:
				assert.Equal(t, want, nil)
			case *model.ErrorResponse:
				testErrorResponse(t, want, w)
			}
		})
	}
}
