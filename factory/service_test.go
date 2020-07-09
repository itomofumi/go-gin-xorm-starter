package factory_test

import (
	"testing"

	"github.com/itomofumi/go-gin-xorm-starter/factory"
	"github.com/itomofumi/go-gin-xorm-starter/infra"
)

type EngineMock struct {
	infra.EngineInterface
}

type KVSClientMock struct {
	infra.KVSClientInterface
}

func TestNew(t *testing.T) {
	factory := factory.NewService(&EngineMock{}, &KVSClientMock{})
	factory.NewFruits()
	factory.NewUsers()
}
