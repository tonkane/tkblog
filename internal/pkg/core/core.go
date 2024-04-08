package core

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/errno"
)

// 定义返回的结构体
type ErrResponse struct {
	Code string `json:"code"`
	Message string `json:"message"`
}

// 包装返回方法，将错误信息写入
func WriteResponse(c *gin.Context, err error, data interface{}) {
	if err != nil {
		httpCode, code, message := errno.Decode(err)
		c.JSON(httpCode, ErrResponse{
			code, message,
		})
		return
	}

	c.JSON(http.StatusOK, data)
}