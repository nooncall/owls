package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nooncall/owls/go/api/v1"
	"github.com/nooncall/owls/go/middleware"
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
