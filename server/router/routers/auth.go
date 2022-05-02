
package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/qingfeng777/owls/server/api/v1"
	"github.com/qingfeng777/owls/server/middleware"
)

func (s *Routers) InitAuthRouter(Router *gin.RouterGroup) {
	recordRouter := Router.Use(middleware.OperationRecord())

	{
		auth := v1.ApiGroupApp.Auth

		recordRouter.POST("/auth", auth.AddAuth)
		recordRouter.DELETE("/auth", auth.DelAuth)
		Router.POST("/auth/list", auth.ListAuth)
		Router.POST("/data-type/list", auth.ListDataType)
	}
}
