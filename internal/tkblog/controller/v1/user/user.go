package user

import (
	"github.com/tkane/tkblog/internal/tkblog/biz"
	"github.com/tkane/tkblog/internal/tkblog/store"

	"github.com/tkane/tkblog/pkg/auth"
	pb "github.com/tkane/tkblog/pkg/proto/tkblog/v1"
)

type UserCtrl struct {
	a *auth.Authz
	b biz.IBiz
	pb.UnimplementedTkBlogServer
}

func New(ds store.IStore, a *auth.Authz) *UserCtrl {
	return &UserCtrl{a: a, b: biz.NewBiz(ds)}
}