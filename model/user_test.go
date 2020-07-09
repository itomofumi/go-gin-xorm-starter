package model_test

import (
	"testing"

	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/stretchr/testify/assert"
)

func TestUser_UserPublicData(t *testing.T) {
	user := model.User{
		Common: model.Common{
			ID: 1,
		},
	}

	assert.EqualValues(t, 1, user.GetPublicData().UserID)
}
