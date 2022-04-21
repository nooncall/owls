package v1

import (
	"github.com/qingfeng777/owls/server/api/v1/autocode"
	"github.com/qingfeng777/owls/server/api/v1/example"
	"github.com/qingfeng777/owls/server/api/v1/system"
	"github.com/qingfeng777/owls/server/api/v1/tidb_or_mysql"
)

type ApiGroup struct {
	SystemApiGroup   system.ApiGroup
	ExampleApiGroup  example.ApiGroup
	AutoCodeApiGroup autocode.ApiGroup
	TiDBOrMysqlGroup tidb_or_mysql.ApiGroup
}

var ApiGroupApp = new(ApiGroup)
