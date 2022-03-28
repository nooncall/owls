package router

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router/autocode"
	"github.com/flipped-aurora/gin-vue-admin/server/router/example"
	"github.com/flipped-aurora/gin-vue-admin/server/router/system"
	"github.com/flipped-aurora/gin-vue-admin/server/router/tidb_or_mysql"
)

type RouterGroup struct {
	System      system.RouterGroup
	Example     example.RouterGroup
	Autocode    autocode.RouterGroup
	TidbOrMysql tidb_or_mysql.RouterGroup
}

var RouterGroupApp = new(RouterGroup)
