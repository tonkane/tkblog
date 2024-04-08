package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/errno"
	"github.com/tkane/tkblog/internal/pkg/log"
	v1 "github.com/tkane/tkblog/pkg/api/tkblog/v1"
)

func (ctrl *UserCtrl) Create(c *gin.Context) {
	log.C(c).Infow("create user function called!")

	var r v1.CreateUserRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	if err := ctrl.b.Users().Create(c, &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}

	core.WriteResponse(c, nil, nil)
}