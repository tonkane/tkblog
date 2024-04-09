package user

import (
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/errno"
	"github.com/tkane/tkblog/internal/pkg/log"
	v1 "github.com/tkane/tkblog/pkg/api/tkblog/v1"
)

// 登录并返回一个 jwt token
func (ctrl *UserCtrl) Login(c *gin.Context) {
	log.C(c).Infow("Login function is called!")

	var r v1.LoginRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	resp, err := ctrl.b.Users().Login(c, &r)
	if err != nil {
		core.WriteResponse(c, err, nil)
		return 
	}
	core.WriteResponse(c, nil, resp)
}