package user

import (
	"context"
	"regexp"

	"github.com/jinzhu/copier"

	"github.com/tkane/tkblog/internal/tkblog/store"
	"github.com/tkane/tkblog/internal/pkg/errno"
	"github.com/tkane/tkblog/internal/pkg/model"
	v1 "github.com/tkane/tkblog/pkg/api/tkblog/v1"
)

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
}

type userBiz struct {
	ds store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(ds store.IStore) *userBiz {
	return &userBiz{ds: ds}
}

func (b *userBiz) Create(ctx context.Context, r *v1.CreateUserRequest) error {
	var userM model.UserM
	_ = copier.Copy(&userM, r)

	if err := b.ds.Users().Create(ctx, &userM); err != nil {
		// 检查报错信息，如果username有重复就报错？
		if match, _ := regexp.MatchString("Duplicate entry '.*' for key 'username'", err.Error()); match {
			return errno.ErrUserAlreadyExist
		}
		return err
	}

	return nil
}