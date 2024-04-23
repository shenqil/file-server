package router

import (
	"fileServer/app/api"
	"fileServer/util/auth"

	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var _ IRouter = (*Router)(nil)

// RouterSet 注入router
var RouterSet = wire.NewSet(wire.Struct(new(Router), "*"), wire.Bind(new(IRouter), new(*Router)))

// IRouter 注册路由
type IRouter interface {
	Register(app *gin.Engine) error
	Prefixes() []string
}

// Router 路由管理器
type Router struct {
	Auth      auth.Auther
	AvatarAPI *api.Avatar
	FileAPI   *api.FileApi
}

func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

func (a *Router) Prefixes() []string {
	return []string{
		"/file-api/",
	}
}
