package dao

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/qingfeng777/owls/server/global"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	// todo, refactor to config~~
	return global.GVA_DB.Debug()
}
