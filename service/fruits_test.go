package service_test

import (
	"reflect"
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/go-gin-xorm-starter/repository"
	"github.com/gemcook/go-gin-xorm-starter/service"
	"github.com/gemcook/ptr"
)

type fruitsRepositoryMock struct {
	repository.FruitsInterface
	FakeGetAll  func() ([]*model.Fruit, error)
	FakeGetByID func(fruitID uint64) (*model.Fruit, error)
	FakeCreate  func(body *model.FruitBody) (*model.Fruit, error)
	FakeUpdate  func(fruitID uint64, notice *model.FruitBody) (*model.Fruit, error)
	FakeDelete  func(fruitID uint64) error
}

func (fr *fruitsRepositoryMock) GetAll() ([]*model.Fruit, error) {
	return fr.FakeGetAll()
}

func (fr *fruitsRepositoryMock) GetByID(fruitID uint64) (*model.Fruit, error) {
	return fr.FakeGetByID(fruitID)
}

func (fr *fruitsRepositoryMock) Create(body *model.FruitBody) (*model.Fruit, error) {
	return fr.FakeCreate(body)
}

func (fr *fruitsRepositoryMock) Update(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
	return fr.FakeUpdate(fruitID, body)
}

func (fr *fruitsRepositoryMock) Delete(fruitID uint64) error {
	return fr.FakeDelete(fruitID)
}

func TestFruits_GetAll(t *testing.T) {
	type fakes struct {
		getAll func() ([]*model.Fruit, error)
	}
	tests := []struct {
		name    string
		fakes   fakes
		want    []*model.Fruit
		wantErr bool
	}{
		{"success",
			fakes{
				getAll: func() ([]*model.Fruit, error) {
					return []*model.Fruit{&model.Fruit{FruitBody: model.FruitBody{Name: ptr.String("apple")}}}, nil
				}},
			[]*model.Fruit{&model.Fruit{FruitBody: model.FruitBody{Name: ptr.String("apple")}}},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fruitsRepositoryMock{
				FakeGetAll: tt.fakes.getAll,
			}
			f := service.NewFruits(repo)

			got, err := f.GetAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fruits.GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fruits.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFruits_GetByID(t *testing.T) {
	type fakes struct {
		getByID func(fruitID uint64) (*model.Fruit, error)
	}
	type args struct {
		fruitID uint64
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		want    *model.Fruit
		wantErr bool
	}{
		{"success",
			fakes{
				getByID: func(fruitID uint64) (*model.Fruit, error) {
					return &model.Fruit{
						Common:    model.Common{ID: 1},
						FruitBody: model.FruitBody{Name: ptr.String("apple")},
					}, nil
				},
			},
			args{fruitID: 1},
			&model.Fruit{
				Common:    model.Common{ID: 1},
				FruitBody: model.FruitBody{Name: ptr.String("apple")},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fruitsRepositoryMock{
				FakeGetByID: tt.fakes.getByID,
			}
			f := service.NewFruits(repo)

			got, err := f.GetByID(tt.args.fruitID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fruits.GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fruits.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFruits_Create(t *testing.T) {
	type fakes struct {
		create func(body *model.FruitBody) (*model.Fruit, error)
	}
	type args struct {
		body *model.FruitBody
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		want    *model.Fruit
		wantErr bool
	}{
		{"success",
			fakes{
				create: func(body *model.FruitBody) (*model.Fruit, error) {
					return &model.Fruit{
						Common:    model.Common{ID: 2},
						FruitBody: model.FruitBody{Name: ptr.String("apple")},
					}, nil
				}},
			args{
				body: &model.FruitBody{
					Name: ptr.String("apple"),
				},
			},
			&model.Fruit{
				Common:    model.Common{ID: 2},
				FruitBody: model.FruitBody{Name: ptr.String("apple")},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fruitsRepositoryMock{
				FakeCreate: tt.fakes.create,
			}
			f := service.NewFruits(repo)

			got, err := f.Create(tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fruits.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fruits.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFruits_Update(t *testing.T) {
	type fakes struct {
		update func(fruitID uint64, body *model.FruitBody) (*model.Fruit, error)
	}
	type args struct {
		fruitID uint64
		body    *model.FruitBody
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		want    *model.Fruit
		wantErr bool
	}{
		{"success",
			fakes{
				update: func(fruitID uint64, body *model.FruitBody) (*model.Fruit, error) {
					return &model.Fruit{
						Common:    model.Common{ID: fruitID},
						FruitBody: *body,
					}, nil
				},
			},
			args{
				fruitID: 1,
				body:    &model.FruitBody{Name: ptr.String("apple")},
			},
			&model.Fruit{
				Common:    model.Common{ID: 1},
				FruitBody: model.FruitBody{Name: ptr.String("apple")},
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fruitsRepositoryMock{
				FakeUpdate: tt.fakes.update,
			}
			f := service.NewFruits(repo)

			got, err := f.Update(tt.args.fruitID, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Fruits.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fruits.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFruits_Delete(t *testing.T) {
	type fakes struct {
		delete func(fruitID uint64) error
	}
	type args struct {
		fruitID uint64
	}
	tests := []struct {
		name    string
		fakes   fakes
		args    args
		wantErr bool
	}{
		{"success",
			fakes{
				delete: func(fruitID uint64) error {
					return nil
				},
			},
			args{fruitID: 1},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fruitsRepositoryMock{
				FakeDelete: tt.fakes.delete,
			}
			f := service.NewFruits(repo)

			if err := f.Delete(tt.args.fruitID); (err != nil) != tt.wantErr {
				t.Errorf("Fruits.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
