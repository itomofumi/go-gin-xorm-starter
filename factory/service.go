package factory

import (
	"github.com/gemcook/go-gin-xorm-starter/infra"
	"github.com/gemcook/go-gin-xorm-starter/repository"
	"github.com/gemcook/go-gin-xorm-starter/service"
)

const (
	// ServiceKey はサービスファクトリ取得キー名
	ServiceKey = "service_factory"
)

// Servicer はサービスファクトリ
type Servicer interface {
	NewUsers() service.UsersInterface
	NewFruits() service.FruitsInterface
}

// ServiceFactory はサービスファクトリの実装
// インフラ層の依存情報を初期化時に注入する
type ServiceFactory struct {
	engine infra.EngineInterface
}

// NewService initializes factory with injected infra.
func NewService(engine infra.EngineInterface) *ServiceFactory {
	r := &ServiceFactory{
		engine: engine,
	}
	return r
}

// NewFruits returns Fruits service.
func (r *ServiceFactory) NewFruits() service.FruitsInterface {
	repo := repository.NewFruits(r.engine)
	return service.NewFruits(repo)
}

// NewUsers returns Users service.
func (r *ServiceFactory) NewUsers() service.UsersInterface {
	repo := repository.NewUsers(r.engine)
	return service.NewUsers(repo)
}
