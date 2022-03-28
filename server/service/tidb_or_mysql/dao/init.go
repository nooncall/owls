package dao

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	//return DB
	return global.GVA_DB
}

func InitDB() {
	conn := fmt.Sprintf("%s:%s@(%s:%d)/%s?charset=utf8mb4", config.Conf.DB.User,
		config.Conf.DB.Password, config.Conf.DB.Address, config.Conf.DB.Port, config.Conf.DB.DBName)

	logger.Info(conn)
/*	var err error
	gorm.Open()
	DB, err = gorm.Open("mysql", conn)
	if err != nil {
		logger.Errorf("init DB error : %s", err.Error())
		log.Fatal("init DB error : ", err.Error())
	}*/

	/*if config.Conf.DB.MaxIdleConn > 0 {
		DB.DB().SetMaxIdleConns(config.Conf.DB.MaxIdleConn)
	}
	if config.Conf.DB.MaxOpenConn > 0 {
		DB.DB().SetMaxOpenConns(config.Conf.DB.MaxOpenConn)
	}
	DB.SingularTable(true)
	DB.Callback().Update().Replace("gorm:update_time_stamp", func(scope *gorm.Scope) {})

	if config.Conf.Server.ShowSql {
		DB.LogMode(true)
		DB.Logger.LogMode(gorm.)*/
	//}
}
