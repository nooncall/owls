package router

import (
	"github.com/nooncall/owls/go/router/autocode"
	"github.com/nooncall/owls/go/router/example"
	"github.com/nooncall/owls/go/router/routers"
	"github.com/nooncall/owls/go/router/system"
	"github.com/nooncall/owls/go/router/tidb_or_mysql"
)

type RouterGroup struct {
	System      system.RouterGroup
	Example     example.RouterGroup
	Autocode    autocode.RouterGroup
	TidbOrMysql tidb_or_mysql.RouterGroup
	Routers     routers.Routers
}

var RouterGroupApp = new(RouterGroup)
