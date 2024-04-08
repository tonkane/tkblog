package model

import (
	"time"

	"gorm.io/gorm"
	"github.com/tkane/tkblog/pkg/auth"
)

type UserM struct {
	ID        int64     `gorm:"column:id;primary_key"` //
	Username  string    `gorm:"column:username"`       //
	Password  string    `gorm:"column:password"`       //
	Nickname  string    `gorm:"column:nickname"`       //
	Email     string    `gorm:"column:email"`          //
	Phone     string    `gorm:"column:phone"`          //
	CreatedAt time.Time `gorm:"column:createdAt"`      //
	UpdatedAt time.Time `gorm:"column:updatedAt"`      //
}

// TableName sets the insert table name for this struct type
func (u *UserM) TableName() string {
	return "user"
}

// GORM 的 BeforeCreate Hook 功能
func (u *UserM) BeforeCreate(tx *gorm.DB) (err error) {
	// 加密
	u.Password, err = auth.Encrypt(u.Password)
	if err != nil {
		return err
	}
	return nil
}
