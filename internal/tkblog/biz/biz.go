package biz

import (
	"github.com/tkane/tkblog/internal/tkblog/biz/user"
	"github.com/tkane/tkblog/internal/tkblog/store"
)

type IBiz interface {
	Users() user.UserBiz
}

type biz struct {
	ds store.IStore
}

var _ IBiz = (*biz)(nil)

func NewBiz(ds store.IStore) *biz {
	return &biz{ds: ds}
}

func (b *biz) Users() user.UserBiz {
	return user.New(b.ds)
}