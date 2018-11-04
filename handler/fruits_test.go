package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/handler"
	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/service"
	"github.com/gemcook/ptr"
	"github.com/stretchr/testify/assert"
)

// FruitsMock is a mock of fruits.
type FruitsMock struct {
	service.FruitsInterface
	FakeGetAll  func() ([]*model.Fruit, error)
	FakeGetByID func(fruitID uint64) (*model.Fruit, error)
	FakeCreate  func(body *model.FruitBody) (*model.Fruit, error)
	FakeUpdate  func(fruitID uint64, body *model.FruitBody) (*model.Fruit, error)
	FakeDelete  func(fruitID uint64) error
}

func (fm *FruitsMock) GetAll() ([]*model.Fruit, error) {
	return fm.FakeGetAll()
}

func (fm *FruitsMock) GetByID(fruitID uint64) (*model.Fruit, error) {
	return fm.FakeGetByID(fruitID)
}

func (fm *FruitsMock) Create(body *model.FruitBody) (*model.Fruit, error) {
	return fm.FakeCreate(body)
}

func (fm *FruitsMock) Update(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
	return fm.FakeUpdate(fruitID, body)
}

func (fm *FruitsMock) Delete(fruitID uint64) error {
	return fm.FakeDelete(fruitID)
}

var testFruits = []*model.Fruit{
	{
		Common: model.Common{ID: 1},
		FruitBody: model.FruitBody{
			Price: ptr.Int(100),
			Name:  ptr.String("Apple"),
		},
	},
	{
		Common: model.Common{ID: 2},
		FruitBody: model.FruitBody{
			Price: ptr.Int(200),
			Name:  ptr.String("Mango"),
		},
	},
}

func TestGetFruits(t *testing.T) {
	defer Setup()()

	type fakes struct {
		getAll func() ([]*model.Fruit, error)
	}
	tests := []struct {
		name       string
		fakes      fakes
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				getAll: func() ([]*model.Fruit, error) {
					return testFruits, nil
				},
			},
			http.StatusOK,
			testFruits,
		},
		{"bad request",
			fakes{
				getAll: func() ([]*model.Fruit, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeGetAll: tt.fakes.getAll,
			}
			factory := &ServiceFactoryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(factory)
			handler.GetFruits(c)
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
	type fakes struct {
		getByID func(id uint64) (*model.Fruit, error)
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name       string
		fakes      fakes
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				getByID: func(id uint64) (*model.Fruit, error) {
					for _, v := range testFruits {
						if v.ID == id {
							return v, nil
						}
					}
					return nil, fmt.Errorf("data not found for id = %v", id)
				},
			},
			args{id: 1},
			http.StatusOK,
			testFruits[0],
		},
		{"bad request",
			fakes{
				getByID: func(id uint64) (*model.Fruit, error) {
					return nil, fmt.Errorf("data not found for id = %v", id)
				},
			},
			args{id: 9999},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("data not found for id = %v", 9999)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeGetByID: tt.fakes.getByID,
			}
			factory := &ServiceFactoryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(factory)
			c.Set("fruit-id", tt.args.id)
			handler.GetFruitByID(c)
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

	type fakes struct {
		create func(body *model.FruitBody) (*model.Fruit, error)
	}
	type args struct {
		body interface{}
	}
	tests := []struct {
		name       string
		fakes      fakes
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				create: func(body *model.FruitBody) (*model.Fruit, error) {
					return testFruits[0], nil
				},
			},
			args{
				body: testFruits[0].FruitBody,
			},
			http.StatusCreated,
			testFruits[0],
		},
		{"bad request",
			fakes{
				create: func(body *model.FruitBody) (*model.Fruit, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			args{
				body: testFruits[0].FruitBody,
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
		{"wrong body type",
			fakes{
				create: nil,
			},
			args{
				body: struct {
					Price string `json:"price"`
				}{Price: "aaa"},
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "json: cannot unmarshal string into Go struct field FruitBody.price of type int"),
		},
		{"validation error: price is minus",
			fakes{
				create: nil,
			},
			args{
				body: model.FruitBody{Name: ptr.String("Apple"), Price: ptr.Int(-10)},
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "Key: 'FruitBody.Price' Error:Field validation for 'Price' failed on the 'notminus' tag"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeCreate: tt.fakes.create,
			}
			factory := &ServiceFactoryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(factory)
			b, _ := json.Marshal(tt.args.body)
			c.Request, _ = http.NewRequest("POST", "/fruits", bytes.NewBuffer(b))

			handler.PostFruit(c)

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

	type fakes struct {
		update func(fruitID uint64, body *model.FruitBody) (*model.Fruit, error)
	}
	type args struct {
		id   uint64
		body *model.FruitBody
	}
	tests := []struct {
		name       string
		fakes      fakes
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				update: func(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
					d := testFruits[0]
					d.FruitBody = *body
					return d, nil
				},
			},
			args{
				id:   1,
				body: &testFruits[0].FruitBody,
			},
			http.StatusOK,
			testFruits[0],
		},
		{"bad request",
			fakes{
				update: func(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
					return nil, fmt.Errorf("some error")
				},
			},
			args{
				id:   1,
				body: &testFruits[0].FruitBody,
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, "some error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeUpdate: tt.fakes.update,
			}
			factory := &ServiceFactoryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(factory)
			b, _ := json.Marshal(tt.args.body)
			c.Request, _ = http.NewRequest("PUT", "/fruits/:fruit-id", bytes.NewBuffer(b))
			c.Set("fruit-id", tt.args.id)

			handler.PutFruit(c)

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

	type fakes struct {
		delete func(fruitID uint64) error
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name       string
		fakes      fakes
		args       args
		wantStatus int
		want       interface{}
	}{
		{"success",
			fakes{
				delete: func(fruitID uint64) error { return nil },
			},
			args{
				id: 1,
			},
			http.StatusNoContent,
			nil,
		},
		{"bad request",
			fakes{
				delete: func(fruitID uint64) error {
					return fmt.Errorf("data not found for id = %v", fruitID)
				},
			},
			args{
				id: 9999,
			},
			http.StatusBadRequest,
			model.NewErrorResponse("400", model.ErrorParam, fmt.Errorf("data not found for id = %v", 9999)),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fruits := &FruitsMock{
				FakeDelete: tt.fakes.delete,
			}
			factory := &ServiceFactoryMock{
				FruitsMock: fruits,
			}

			c, w := createGinTestContext(factory)
			c.Set("fruit-id", tt.args.id)

			handler.DeleteFruit(c)

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
