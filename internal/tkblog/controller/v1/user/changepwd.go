package user

import (
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/errno"
	"github.com/tkane/tkblog/internal/pkg/log"
	v1 "github.com/tkane/tkblog/pkg/api/tkblog/v1"
)


func (ctrl *UserCtrl) ChangePwd(c *gin.Context) {
	log.C(c).Infow("change pwd function is called!")

	var r v1.ChangePwdRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		core.WriteResponse(c, errno.ErrBind, nil)
		return
	}

	if _, err := govalidator.ValidateStruct(r); err != nil {
		core.WriteResponse(c, errno.ErrInvalidParameter.SetMessage(err.Error()), nil)
		return
	}

	// 奇葩，为啥把名字放URL路径上
	if err := ctrl.b.Users().ChangePwd(c, c.Param("name"), &r); err != nil {
		core.WriteResponse(c, err, nil)
		return
	}
	
	// 不写返回值你是真的懂
	core.WriteResponse(c, nil, gin.H{"status": "ok"})
}