package user

import (
	"context"
	"regexp"

	"github.com/jinzhu/copier"

	"github.com/tkane/tkblog/internal/tkblog/store"
	"github.com/tkane/tkblog/internal/pkg/errno"
	"github.com/tkane/tkblog/internal/pkg/model"
	"github.com/tkane/tkblog/pkg/token"
	"github.com/tkane/tkblog/pkg/auth"
	v1 "github.com/tkane/tkblog/pkg/api/tkblog/v1"
)

type UserBiz interface {
	Create(ctx context.Context, r *v1.CreateUserRequest) error
	Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error)
	ChangePwd(ctx context.Context, username string, r *v1.ChangePwdRequest) error 
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

func (b *userBiz) Login(ctx context.Context, r *v1.LoginRequest) (*v1.LoginResponse, error) {
	user, err := b.ds.Users().Get(ctx, r.Username)
	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	if err := auth.Compare(user.Password, r.Password); err != nil {
		return nil, errno.ErrPwdIncorrect
	}

	t, err := token.Sign(r.Username)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &v1.LoginResponse{Token: t}, nil
}

func (b *userBiz) ChangePwd(ctx context.Context, username string, r *v1.ChangePwdRequest) error {
	userM, err := b.ds.Users().Get(ctx, username)
	if err != nil {
		return err
	}

	if err := auth.Compare(userM.Password, r.OldPwd); err != nil {
		return errno.ErrPwdIncorrect
	}

	// create 的时候用了 hook 这里更新却单独写？
	userM.Password, _ = auth.Encrypt(r.NewPwd)
	if err := b.ds.Users().Update(ctx, userM); err != nil {
		return err
	}

	return nil
}