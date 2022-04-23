package autocode

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/qingfeng777/owls/server/api/v1"
	"github.com/qingfeng777/owls/server/middleware"
)

type AutoCodeExampleRouter struct{}

func (s *AutoCodeExampleRouter) InitSysAutoCodeExampleRouter(Router *gin.RouterGroup) {
	autoCodeExampleRouter := Router.Group("autoCodeExample").Use(middleware.OperationRecord())
	autoCodeExampleRouterWithoutRecord := Router.Group("autoCodeExample")
	autoCodeExampleApi := v1.ApiGroupApp.AutoCodeApiGroup.AutoCodeExampleApi
	{
		autoCodeExampleRouter.POST("createSysAutoCodeExample", autoCodeExampleApi.CreateAutoCodeExample)   // 新建AutoCodeExample
		autoCodeExampleRouter.DELETE("deleteSysAutoCodeExample", autoCodeExampleApi.DeleteAutoCodeExample) // 删除AutoCodeExample
		autoCodeExampleRouter.PUT("updateSysAutoCodeExample", autoCodeExampleApi.UpdateAutoCodeExample)    // 更新AutoCodeExample
	}
	{
		autoCodeExampleRouterWithoutRecord.GET("findSysAutoCodeExample", autoCodeExampleApi.FindAutoCodeExample)       // 根据ID获取AutoCodeExample
		autoCodeExampleRouterWithoutRecord.GET("getSysAutoCodeExampleList", autoCodeExampleApi.GetAutoCodeExampleList) // 获取AutoCodeExample列表
	}
}
