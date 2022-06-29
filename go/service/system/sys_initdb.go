package system

import (
	"database/sql"
	"fmt"

	adapter "github.com/casbin/gorm-adapter/v3"
	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/example"
	"github.com/qingfeng777/owls/server/model/system"
	"github.com/qingfeng777/owls/server/model/system/request"
	"github.com/qingfeng777/owls/server/service/auth/auth"
	tasks "github.com/qingfeng777/owls/server/service/task"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/checker"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/db_info"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/task"
)

type InitDBService struct{}

// InitDB 创建数据库并初始化 总入口
// Author [piexlmax](https://github.com/piexlmax)
// Author [SliverHorn](https://github.com/SliverHorn)
// Author [songzhibin97](https://github.com/songzhibin97)
func (initDBService *InitDBService) InitDB(conf request.InitDB) error {
	switch conf.DBType {
	case "mysql":
		return initDBService.initMysqlDB(conf)
	case "pgsql":
		return initDBService.initPgsqlDB(conf)
	default:
		return initDBService.initMysqlDB(conf)
	}
}

// initTables 初始化表
// Author [SliverHorn](https://github.com/SliverHorn)
func (initDBService *InitDBService) initTables() error {
	return global.GVA_DB.AutoMigrate(
		system.SysApi{},
		system.SysUser{},
		system.SysBaseMenu{},
		system.SysAuthority{},
		system.JwtBlacklist{},
		system.SysDictionary{},
		system.SysAutoCodeHistory{},
		system.SysOperationRecord{},
		system.SysDictionaryDetail{},
		system.SysBaseMenuParameter{},
		system.SysBaseMenuBtn{},
		system.SysAuthorityBtn{},

		adapter.CasbinRule{},

		example.ExaFile{},
		example.ExaCustomer{},
		example.ExaFileChunk{},
		example.ExaFileUploadAndDownload{},

		task.OwlTask{},
		task.OwlBackup{},
		task.OwlSubtask{},
		task.OwlExecItem{},
		checker.OwlRuleStatus{},
		db_info.OwlCluster{},

		tasks.Task{},
		auth.Auth{},
	)
}

// createDatabase 创建数据库(mysql)
// Author [SliverHorn](https://github.com/SliverHorn)
// Author: [songzhibin97](https://github.com/songzhibin97)

func (initDBService *InitDBService) createDatabase(dsn string, driver string, createSql string) error {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		return err
	}
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(db)
	if err = db.Ping(); err != nil {
		return err
	}
	_, err = db.Exec(createSql)
	return err
}
