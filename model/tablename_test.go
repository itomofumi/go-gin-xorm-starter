package model_test

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/model"
	"github.com/stretchr/testify/assert"
)

func TestTableName(t *testing.T) {

	type tableNamer interface{ TableName() string }
	tests := []struct {
		m     tableNamer
		table string
	}{
		{model.Common{}, ""},
		{model.Fruit{}, "fruits"},
		{model.User{}, "users"},
		{model.UserPublicData{}, "users"},
	}

	for _, tt := range tests {
		t.Run(tt.table, func(t *testing.T) {
			assert.Equal(t, tt.table, tt.m.TableName())
		})
	}
}
