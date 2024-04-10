package user

import (
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/log"
)

// 一个 get 方法放一个go，不太行
func (ctrl *UserCtrl) Get(c *gin.Context) {
	log.C(c).Infow("get user function called")

	user, err := ctrl.b.Users().Get(c, c.Param("name"))
	if err != nil {
		core.WriteResponse(c, err, nil)
		return 
	}

	core.WriteResponse(c, nil, user)
}