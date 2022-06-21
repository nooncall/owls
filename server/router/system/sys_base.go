package system

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/qingfeng777/owls/server/api/v1"
)

type BaseRouter struct{}

func (s *BaseRouter) InitBaseRouter(Router *gin.RouterGroup) (R gin.IRoutes) {
	baseRouter := Router.Group("base")
	baseApi := v1.ApiGroupApp.SystemApiGroup.BaseApi
	{
		baseRouter.GET("login/model", baseApi.LoginModel)
		baseRouter.POST("register", baseApi.Register)
		baseRouter.POST("login", baseApi.Login)
		baseRouter.POST("captcha", baseApi.Captcha)
	}
	return baseRouter
}
