package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/errno"
	"github.com/tkane/tkblog/internal/pkg/known"

	"github.com/tkane/tkblog/pkg/token"
)

// 鉴权中间件
func Authn() gin.HandlerFunc {
	return func (c *gin.Context) {
		username, err := token.ParseRequest(c)
		if err != nil {
			core.WriteResponse(c, errno.ErrTokenInvalid, nil)
			c.Abort()
			return
		}
		c.Set(known.XusernameKey, username)
		c.Next()
	}
}