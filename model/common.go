package model

import (
	"time"

	"github.com/gemcook/ptr"
)

// Common テーブル共通項目を定義
type Common struct {
	ID        uint64     `xorm:"pk autoincr index(pk)" json:"id"`
	IsDeleted *bool      `xorm:"default false notnull" json:"-"`
	IsEnabled *bool      `xorm:"default true notnull" json:"-"`
	CreatedAt *time.Time `xorm:"created notnull" json:"-"`
	UpdatedAt *time.Time `xorm:"updated notnull" json:"-"`
}

// TableName should not be called
func (Common) TableName() string {
	return ""
}

// SetDefault sets init data
func (m *Common) SetDefault() {
	m.IsDeleted = ptr.Bool(false)
	m.IsEnabled = ptr.Bool(true)
	m.CreatedAt = nil
	m.UpdatedAt = nil
}
