package tkblog

import (
	"github.com/gin-gonic/gin"

	"github.com/tkane/tkblog/internal/pkg/log"
	"github.com/tkane/tkblog/internal/pkg/core"
	"github.com/tkane/tkblog/internal/pkg/errno"

	"github.com/tkane/tkblog/internal/tkblog/store"
	"github.com/tkane/tkblog/internal/tkblog/controller/v1/user"

	mw "github.com/tkane/tkblog/internal/pkg/middleware"

	"github.com/tkane/tkblog/pkg/auth"
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

	authz, err := auth.NewAuthz(store.S.DB())
	if err != nil {
		return err
	}

	uc := user.New(store.S, authz)

	g.POST("/login", uc.Login)

	v1 := g.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			userv1.POST("", uc.Create)
			userv1.PUT(":name/change-password", uc.ChangePwd)
			userv1.Use(mw.Authn(), mw.Authz(authz))
			// Auth 中间件在这之后才有效
			userv1.GET("test", func(c *gin.Context){
				core.WriteResponse(c, nil, gin.H{"status": "ok"})
			})
			userv1.GET(":name", uc.Get)
		}
	}

	return nil
}