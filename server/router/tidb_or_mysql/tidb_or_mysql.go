package tidb_or_mysql

import (
	v1 "github.com/flipped-aurora/gin-vue-admin/server/api/v1"
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type TidbOrMysqlRouter struct{}

func (s *TidbOrMysqlRouter) InitApiRouter(Router *gin.RouterGroup) {
	apiRouter := Router.Group("db").Use(middleware.OperationRecord())
	apiRouterWithoutRecord := Router.Group("db")

	{
		taskApi := v1.ApiGroupApp.TiDBOrMysqlGroup.TaskApi

		apiRouter.POST("/task", taskApi.AddTask)
		apiRouter.PUT("/task", taskApi.UpdateTask)
		apiRouterWithoutRecord.GET("/task", taskApi.GetTask)
		apiRouterWithoutRecord.GET("/task/review", taskApi.ListReviewerTask)
		apiRouterWithoutRecord.GET("/task/exec", taskApi.ListExecTask)
		apiRouterWithoutRecord.GET("/task/history", taskApi.ListHistoryTask)
	}

	{
		backupApi := v1.ApiGroupApp.TiDBOrMysqlGroup.BackupApi

		apiRouter.POST("/backup/rollback", backupApi.Rollback)
		apiRouterWithoutRecord.GET("/backup/list", backupApi.ListRollbackData)
	}

	{
		clusterApi := v1.ApiGroupApp.TiDBOrMysqlGroup.ClusterApi

		apiRouter.POST("/cluster", clusterApi.AddCluster) // without record can paas
		apiRouter.PUT("/cluster", clusterApi.UpdateCluster)
		apiRouter.DELETE("/cluster", clusterApi.DelCluster)
		apiRouterWithoutRecord.POST("/cluster/list", clusterApi.ListCluster)
		apiRouterWithoutRecord.GET("/cluster/db/list", clusterApi.ListDB)
	}

	{
		adminApi := v1.ApiGroupApp.TiDBOrMysqlGroup.AdminApi

		apiRouter.POST("/admin", adminApi.AddAdmin)
		apiRouter.DELETE("/admin", adminApi.DelAdmin)
		apiRouterWithoutRecord.GET("/admin/list", adminApi.ListAdmin)
	}

	{
		ruleApi := v1.ApiGroupApp.TiDBOrMysqlGroup.RuleApi

		apiRouterWithoutRecord.GET("/rule/list", ruleApi.LisRule)
		apiRouter.PUT("/rule/status", ruleApi.UpdateRuleStatus)
	}
}
