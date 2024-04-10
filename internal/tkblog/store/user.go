package store

import (
	"context"

	"gorm.io/gorm"

	"github.com/tkane/tkblog/internal/pkg/model"
)

type UserStore interface {
	Create(ctx context.Context, user *model.UserM) error
	Get(ctx context.Context, username string) (*model.UserM, error)
	Update(ctx context.Context, user *model.UserM) error
	List(ctx context.Context, offset, limit int) (int64, []*model.UserM, error)
}

type users struct {
	db *gorm.DB
}

var _ UserStore = (*users)(nil)

func newUsers(db *gorm.DB) *users {
	return &users{db}
}


func (u *users) Create(ctx context.Context, user *model.UserM) error {
	return u.db.Create(&user).Error
}

// 获取
func (u *users) Get(ctx context.Context, username string) (*model.UserM, error) {
	var user model.UserM
	if err := u.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// 更新记录
func (u *users) Update(ctx context.Context, user *model.UserM) error {
	return u.db.Save(user).Error
}

// List
func (u *users) List(ctx context.Context, offset, limit int) (count int64, ret []*model.UserM, err error) {
	// offset -1 limit -1 ? 难道还要重置 ?
	err = u.db.Offset(offset).Limit(defaultLimit(limit)).Order("id asc").Find(&ret).
		Offset(-1).
		Limit(-1).
		Count(&count).
		Error
	return
}