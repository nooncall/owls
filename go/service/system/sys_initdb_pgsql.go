package system

import (
	"path/filepath"

	"github.com/nooncall/owls/go/config"
	"github.com/nooncall/owls/go/global"
	model "github.com/nooncall/owls/go/model/system"
	"github.com/nooncall/owls/go/source/example"
	"github.com/nooncall/owls/go/source/system"
	"github.com/nooncall/owls/go/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// writePgsqlConfig pgsql 回写配置
// Author [SliverHorn](https://github.com/SliverHorn)
func (initDBService *InitDBService) writePgsqlConfig(pgsql config.Pgsql) error {
	global.GVA_CONFIG.System.DbType = "pgsql"
	global.GVA_CONFIG.Pgsql = pgsql
	cs := utils.StructToMap(global.GVA_CONFIG)
	for k, v := range cs {
		global.GVA_VP.Set(k, v)
	}
	global.GVA_VP.Set("jwt.signing-key", uuid.NewV4().String())
	return global.GVA_VP.WriteConfig()
}

func (initDBService *InitDBService) initPgsqlDB(conf model.InitDBData) error {
	dsn := conf.PgsqlEmptyDsn()
	createSql := "CREATE DATABASE " + conf.DBName
	if err := initDBService.createDatabase(dsn, "pgx", createSql); err != nil {
		return err
	} // 创建数据库

	pgsqlConfig := conf.ToPgsqlConfig()
	if pgsqlConfig.Dbname == "" {
		return nil
	} // 如果没有数据库名, 则跳出初始化数据

	if db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  pgsqlConfig.Dsn(), // DSN data source name
		PreferSimpleProtocol: false,
	}), &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}); err != nil {
		return nil
	} else {
		global.GVA_DB = db
	}

	if err := initDBService.initTables(); err != nil {
		global.GVA_DB = nil
		return err
	}

	if err := initDBService.initPgsqlData(&conf); err != nil {
		global.GVA_DB = nil
		return err
	}

	if err := initDBService.writePgsqlConfig(pgsqlConfig); err != nil {
		return err
	}

	global.GVA_CONFIG.AutoCode.Root, _ = filepath.Abs("..")
	return nil
}

// initPgsqlData pgsql 初始化数据
// Author [SliverHorn](https://github.com/SliverHorn)
func (initDBService *InitDBService) initPgsqlData(initData *model.InitDBData) error {
	return model.PgsqlDataInitialize(initData,
		system.Api,
		system.User,
		system.Casbin,
		system.BaseMenu,
		system.Authority,
		system.Dictionary,
		system.UserAuthority,
		system.DataAuthorities,
		system.AuthoritiesMenus,
		system.DictionaryDetail,
		system.ViewAuthorityMenuPostgres,

		example.FilePgsql,
	)
}
