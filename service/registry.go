package service

import (
	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/gemcook/go-gin-xorm-starter/repository"
)

const (
	// RegistryKey はサービスレジストリ取得キー名
	RegistryKey = "service_registry"
)

// RegistryInterface はサービスレジストリ
type RegistryInterface interface {
	NewUsers() UsersInterface
	NewFruits() FruitsInterface
}

// Registry はサービスレジストリの実装
// インフラ層の依存情報を初期化時に注入する
type Registry struct {
	engine infra.EngineInterface
}

// NewRegistry initializes registry with injected orm  は依存するORMを注入してサービスレジストリを初期化
func NewRegistry(engine infra.EngineInterface) RegistryInterface {
	r := &Registry{
		engine: engine,
	}
	return r
}

// NewFruits returns Fruits service.
func (r *Registry) NewFruits() FruitsInterface {
	repo := repository.NewFruits(r.engine)
	return NewFruits(repo)
}

// NewUsers returns Users service.
func (r *Registry) NewUsers() UsersInterface {
	repo := repository.NewUsers(r.engine)
	return NewUsers(repo)
}
