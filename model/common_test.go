package model_test

import (
	"testing"

	"github.com/itomofumi/go-gin-xorm-starter/model"
	"github.com/stretchr/testify/assert"
)

func TestCommon_SetDefault(t *testing.T) {
	common := &model.Common{}

	common.SetDefault()
	assert := assert.New(t)
	assert.EqualValues(true, *common.IsEnabled)
	assert.EqualValues(false, *common.IsDeleted)
}
