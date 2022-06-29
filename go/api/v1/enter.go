package v1

import (
	"github.com/qingfeng777/owls/server/api/v1/auth"
	"github.com/qingfeng777/owls/server/api/v1/autocode"
	"github.com/qingfeng777/owls/server/api/v1/example"
	"github.com/qingfeng777/owls/server/api/v1/system"
	"github.com/qingfeng777/owls/server/api/v1/task"
	"github.com/qingfeng777/owls/server/api/v1/tidb_or_mysql"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	AutoCodeApiGroup autocode.ApiGroup
	TiDBOrMysqlGroup tidb_or_mysql.ApiGroup
	Task             task.TaskApi
	Auth             auth.AuthApi
}

var ApiGroupApp = new(ApiGroup)

type ListData struct {
	Total    int64       `json:"total"`
	List     interface{} `json:"list"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}
