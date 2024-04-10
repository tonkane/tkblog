package user

import (
	"github.com/tkane/tkblog/internal/tkblog/biz"
	"github.com/tkane/tkblog/internal/tkblog/store"

	"github.com/tkane/tkblog/pkg/auth"
)

type UserCtrl struct {
	a *auth.Authz
	b biz.IBiz
}

func New(ds store.IStore, a *auth.Authz) *UserCtrl {
	return &UserCtrl{a: a, b: biz.NewBiz(ds)}
}