package model

import (
	"testing"

	"github.com/gemcook/ptr"
)

func TestFruitBody_IsValid(t *testing.T) {
	type fields struct {
		Name  string
		Price int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"valid", fields{"Apple", 100}, true},
		{"invalid: no name", fields{"", 100}, false},
		{"invalid: minus price", fields{"Mango", -100}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FruitBody{
				Name:  ptr.String(tt.fields.Name),
				Price: ptr.Int(tt.fields.Price),
			}
			if got := f.IsValid(); got != tt.want {
				t.Errorf("FruitBody.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}
