package user

import (
	"github.com/tkane/tkblog/internal/tkblog/biz"
	"github.com/tkane/tkblog/internal/tkblog/store"
)

type UserCtrl struct {
	b biz.IBiz
}

func New(ds store.IStore) *UserCtrl {
	return &UserCtrl{b: biz.NewBiz(ds)}
}