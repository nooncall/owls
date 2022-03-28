package task

import (
	"database/sql"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

type DBInfo struct {
	DB        *sql.DB
	DefaultDB *sql.DB
	DBName    string
}

func (db *DBInfo) CloseConn() {
	if err := db.DB.Close(); err != nil {
		logger.Errorf("close db conn err: %s", err.Error())
	}
	if err := db.DefaultDB.Close(); err != nil {
		logger.Errorf("close default db conn err: %s", err.Error())
	}
}

type dbTools interface {
	GetDBConn(dbName, cluster string) (*DBInfo, error)
}

var dbTool dbTools

func SetDbTools(impl dbTools) {
	dbTool = impl
}

type OwlAccount struct {
	ID       uint64 `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Passwd   string `json:"passwd" gorm:"column:passwd"`
	Dbtype   int64  `json:"dbtype" gorm:"column:dbtype"`
}
