package redis

import (
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/gorm"

	"github.com/nooncall/owls/go/global"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	// todo, refactor to config~~
	return global.GVA_DB.Debug()
}
