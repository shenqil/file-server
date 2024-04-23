package router

import (
	"fileServer/app/middleware"
	_ "fileServer/docs"

	"github.com/gin-gonic/gin"
)

// RegisterAPI 注册 api 组路由器
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/file-api")

	g.Use(middleware.UserAuthMiddleware(a.Auth))
	g.Use(middleware.RateLimiterMiddleware())

	v1 := g.Group("/v1")
	{
		gFiles := v1.Group("files")
		{
			// gFiles.GET(":name", a.Get)
			gFiles.POST("", a.FileAPI.Upload)
			// gFiles.DELETE(":name", a.DemoAPI.Delete)
		}

		gAvatar := v1.Group("avatars")
		{
			// gFiles.GET(":name", a.Get)
			gAvatar.POST("", a.AvatarAPI.Upload)
			// gFiles.DELETE(":name", a.DemoAPI.Delete)
		}
	}
}
