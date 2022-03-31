package injection

import (
	service "github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/admin"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/auth"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/auth/login_check"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/checker"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/dao"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/db_info"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/task"
)

func Injection()  {
	task.SetBackupDao(dao.BackupDAO)
	//task.SetTaskDao(dao.Task)
	task.SetSubTaskDao(dao.SubTask)
	task.SetDbTools(db_info.DBTool)
	task.SetChecker(checker.Checker)
	checker.SetRuleStatusDao(dao.Rule)
	db_info.SetClusterDao(dao.Cluster)
	auth.SetLoginService(login_check.LoginService)
	service.SetClock(service.RealClock{})
	admin.SetAdminDao(dao.Admin)
}