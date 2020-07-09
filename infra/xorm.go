package infra

import "xorm.io/xorm"

// EngineSelector provides Select function.
type EngineSelector interface {
	Select(str string) *xorm.Session
}

// EngineInterface  xorm.EngineInterfaceの不足を追加
type EngineInterface interface {
	xorm.EngineInterface
	EngineSelector
}
