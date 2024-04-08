package tkblog

import (
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/log"
	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/errno"

	"github.com/tkane/tkblog/internal/tkblog/store"
	"github.com/tkane/tkblog/internal/tkblog/controller/v1/user"
)

func installRouters(g *gin.Engine) error {
	// 404 handler
	g.NoRoute(func (c *gin.Context)  {
		core.WriteResponse(c, errno.ErrPageNotFound, nil)
	})
	// healthz handler
	g.GET("/healthz", func (c *gin.Context)  {
		// 打印 X-request-id
		log.C(c).Infow("healthz is called!")
		core.WriteResponse(c, nil, gin.H{"status": "ok"})
	})

	uc := user.New(store.S)

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
		}
	}

	return nil
}