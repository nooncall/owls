package routers

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/qingfeng777/owls/server/api/v1"
	"github.com/qingfeng777/owls/server/middleware"
)

func (s *Routers) InitTaskRouter(Router *gin.RouterGroup) {
	recordRouter := Router.Use(middleware.OperationRecord())

	{
		taskApi := v1.ApiGroupApp.Task

		recordRouter.POST("/task", taskApi.AddTask)
		recordRouter.PUT("/task", taskApi.UpdateTask)
		Router.GET("/task", taskApi.GetTask)
		Router.POST("/task/list", taskApi.ListTask)
		Router.POST("/task/review", taskApi.ListReviewTask)
		Router.POST("/task/history", taskApi.ListHistoryTask)
	}
}
