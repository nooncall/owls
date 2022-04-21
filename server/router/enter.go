package router

import (
	"github.com/qingfeng777/owls/server/router/autocode"
	"github.com/qingfeng777/owls/server/router/example"
	"github.com/qingfeng777/owls/server/router/system"
	"github.com/qingfeng777/owls/server/router/tidb_or_mysql"
)

type RouterGroup struct {
	System      system.RouterGroup
	Example     example.RouterGroup
	Autocode    autocode.RouterGroup
	TidbOrMysql tidb_or_mysql.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
