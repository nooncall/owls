package tidb_or_mysql

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/nooncall/owls/go/api/v1"
	"github.com/nooncall/owls/go/middleware"
)

type TidbOrMysqlRouter struct{}

func (s *TidbOrMysqlRouter) InitApiRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("db").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("db")

	{
		readApi := v1.ApiGroupApp.TiDBOrMysqlGroup.ReadApi

		apiRouter.POST("/read", readApi.ReadData)
		apiRouterWithoutRecord.POST("/table", readApi.GetTableInfo)
	}

	{
		taskApi := v1.ApiGroupApp.TiDBOrMysqlGroup.TaskApi

		apiRouter.POST("/task", taskApi.AddTask)
		apiRouter.PUT("/task", taskApi.UpdateTask)
		apiRouterWithoutRecord.GET("/task", taskApi.GetTask)
		apiRouterWithoutRecord.POST("/task/list", taskApi.ListTask)
		apiRouterWithoutRecord.POST("/task/review", taskApi.ListReviewTask)
		apiRouterWithoutRecord.POST("/task/history", taskApi.ListHistoryTask)
	}

	{
		backupApi := v1.ApiGroupApp.TiDBOrMysqlGroup.BackupApi

		apiRouter.POST("/backup/rollback", backupApi.Rollback)
		apiRouterWithoutRecord.POST("/backup/list", backupApi.ListRollbackData)
	}

	{
		clusterApi := v1.ApiGroupApp.TiDBOrMysqlGroup.ClusterApi

		apiRouter.POST("/cluster", clusterApi.AddCluster) // without record can paas
		apiRouter.PUT("/cluster", clusterApi.UpdateCluster)
		apiRouter.DELETE("/cluster", clusterApi.DelCluster)
		apiRouter.POST("/cluster/list", clusterApi.ListCluster)
		apiRouter.GET("/cluster/name/list", clusterApi.ListClusterName)
		apiRouter.GET("/cluster/db/list", clusterApi.ListDB)
		apiRouter.GET("/cluster/table/list", clusterApi.ListTable)
	}

	{
		adminApi := v1.ApiGroupApp.TiDBOrMysqlGroup.AdminApi

		apiRouter.POST("/admin", adminApi.AddAdmin)
		apiRouter.DELETE("/admin", adminApi.DelAdmin)
		apiRouterWithoutRecord.GET("/admin/list", adminApi.ListAdmin)
	}

	{
		ruleApi := v1.ApiGroupApp.TiDBOrMysqlGroup.RuleApi

		apiRouterWithoutRecord.POST("/rule/list", ruleApi.LisRule)
		apiRouter.PUT("/rule/status", ruleApi.UpdateRuleStatus)
	}
}
