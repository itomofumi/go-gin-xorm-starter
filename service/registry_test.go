package service

import (
	"testing"

	"github.com/gemcook/go-gin-xorm-starter/infra"
)

type EngineMock struct {
	infra.EngineInterface
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry(&EngineMock{})
	registry.NewFruits()
	registry.NewUsers()
}
