package request

import (
	"fmt"

	"github.com/nooncall/owls/go/config"
)

type InitDB struct {
	DBType   string `json:"dbType"`                      // 数据库类型
	Host     string `json:"host"`                        // 服务器地址
	Port     string `json:"port"`                        // 数据库连接端口
	UserName string `json:"userName" binding:"required"` // 数据库用户名
	Password string `json:"password"`                    // 数据库密码
	DBName   string `json:"dbName" binding:"required"`   // 数据库名
}

// MysqlEmptyDsn msyql 空数据库 建库链接
// Author SliverHorn
func (i *InitDB) MysqlEmptyDsn() string {
	if i.Host == "" {
		i.Host = "127.0.0.1"
	}
	if i.Port == "" {
		i.Port = "3306"
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/", i.UserName, i.Password, i.Host, i.Port)
}

// PgsqlEmptyDsn pgsql 空数据库 建库链接
// Author SliverHorn
func (i *InitDB) PgsqlEmptyDsn() string {
	if i.Host == "" {
		i.Host = "127.0.0.1"
	}
	if i.Port == "" {
		i.Port = "3306"
	}
	return "host=" + i.Host + " user=" + i.UserName + " password=" + i.Password + " port=" + i.Port + " " + "sslmode=disable TimeZone=Asia/Shanghai"
}

// ToMysqlConfig 转换 config.Mysql
// Author [SliverHorn](https://github.com/SliverHorn)
func (i *InitDB) ToMysqlConfig() config.Mysql {
	return config.Mysql{
		Path:         i.Host,
		Port:         i.Port,
		Dbname:       i.DBName,
		Username:     i.UserName,
		Password:     i.Password,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogMode:      "error",
		Config:       "charset=utf8mb4&parseTime=True&loc=Local",
	}
}

// ToPgsqlConfig 转换 config.Pgsql
// Author [SliverHorn](https://github.com/SliverHorn)
func (i *InitDB) ToPgsqlConfig() config.Pgsql {
	return config.Pgsql{
		Path:         i.Host,
		Port:         i.Port,
		Dbname:       i.DBName,
		Username:     i.UserName,
		Password:     i.Password,
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogMode:      "error",
		Config:       "sslmode=disable TimeZone=Asia/Shanghai",
	}
}
