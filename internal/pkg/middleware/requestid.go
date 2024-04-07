package middleware

import (
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/known"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查请求头中是否有 X-request-ID 有就用，没有就生成
		requestID := c.Request.Header.Get(known.XRequestIDKey)
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// 将 requestID 放入 gin.Context中
		c.Set(known.XRequestIDKey, requestID)

		// 设置 响应头 header
		c.Writer.Header().Set(known.XRequestIDKey, requestID)
		c.Next()
	}
}

