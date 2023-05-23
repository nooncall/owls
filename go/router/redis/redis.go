package redis

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nooncall/owls/go/api/v1"
	"github.com/nooncall/owls/go/middleware"
)

type RedisRouter struct{}

func (s *RedisRouter) InitApiRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("cache").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("cache")

	{
		readApi := v1.ApiGroupApp.Redis.ReadApi

		apiRouter.POST("/read", readApi.ReadData)
	}

	// task 接口复用，task包含写接口

	// admin, cluster 接口复用
	{
		ruleApi := v1.ApiGroupApp.Redis.ReadApi

		apiRouterWithoutRecord.GET("/rule/list", ruleApi.ListRule)
		apiRouter.PUT("/rule/status", ruleApi.UpdateRuleStatus)
	}
}
