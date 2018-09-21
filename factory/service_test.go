package factory_test

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/factory"
	"github.com/gemcook/go-gin-xorm-starter/infra"
)

type EngineMock struct {
	infra.EngineInterface
}

func TestNew(t *testing.T) {
	factory := factory.New(&EngineMock{})
	factory.NewFruits()
	factory.NewUsers()
}
