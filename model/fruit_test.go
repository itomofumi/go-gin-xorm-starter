package model_test

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/gemcook/ptr"
)

func TestFruitBodyStructLevelValidation(t *testing.T) {
	type fields struct {
		Name  string
		Price int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"valid", fields{"Apple", 100}, false},
		{"invalid: no name", fields{"", 100}, true},
		{"invalid: minus price", fields{"Mango", -100}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &model.FruitBody{
				Name:  ptr.String(tt.fields.Name),
				Price: ptr.Int(tt.fields.Price),
			}
			v := &model.StructValidator{}
			if err := v.ValidateStruct(f); (err != nil) != tt.wantErr {
				t.Errorf("StructValidation for FruitBody{} error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
