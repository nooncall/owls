package injection

import (
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/service/redis"
	"github.com/nooncall/owls/go/service/redis/redis_dao"
	"github.com/nooncall/owls/go/service/system"
	"github.com/nooncall/owls/go/service/system/login"
	service "github.com/nooncall/owls/go/service/tidb_or_mysql"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/admin"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/auth"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/auth/login_check"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/checker"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/dao"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/db_info"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/task"
)

func Injection() {
	task.SetBackupDao(dao.BackupDAO)
	task.SetTaskDao(dao.Task)
	task.SetSubTaskDao(dao.SubTask)
	task.SetDbTools(db_info.DBTool)
	task.SetChecker(checker.Checker)
	checker.SetRuleStatusDao(dao.Rule)
	db_info.SetClusterDao(dao.Cluster)
	auth.SetLoginService(login_check.LoginService)
	service.SetClock(service.RealClock{})
	admin.SetAdminDao(dao.Admin)
	redis.SetRedisTaskDao(redis_dao.Task)

	switch global.GVA_CONFIG.Login.Model {
	case "ldap": // todo, const
		system.SetUserLogin(login.LdapUserImpl())
	case "registry":
		system.SetUserLogin(login.RegistryUserImpl())
	default:
		panic("unknown login mode, implement me")
	}

	switch "" { //todo, support auth tool;
	case "conf":
		task.SetAuthTools(auth.ConfAuthService)
	case "net":
		task.SetAuthTools(auth.NetAuthService)
	case "admin":
		task.SetAuthTools(auth.AdminAuthService)
	case "mock":
		// MockInjection()
	default:
		task.SetAuthTools(auth.AdminAuthService)
	}
}
