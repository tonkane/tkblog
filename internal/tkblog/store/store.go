package store

import (
	"sync"

	"gorm.io/gorm"
)

var (
	once sync.Once
	// 全局变量
	S *datastore
)

type IStore interface {
	Users() UserStore
	DB() *gorm.DB
}

type datastore struct {
	db *gorm.DB
}

var _ IStore = (*datastore)(nil)

func NewStore(db *gorm.DB) *datastore {
	once.Do(func ()  {
		S = &datastore{db}
	})

	return S
}

// 公用层应该控制反转？
func (ds *datastore) Users() UserStore {
	return newUsers(ds.db)
}

// 返回DB实例
func (ds *datastore) DB() *gorm.DB {
	return ds.db
}