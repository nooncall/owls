package initialize

import (
	//"github.com/nooncall/owls/go/plugin/email" // 本地插件仓库地址模式
	"github.com/gin-gonic/gin"
	"github.com/nooncall/owls/go/plugin/example_plugin"
	"github.com/nooncall/owls/go/utils/plugin"
)

func PluginInit(group *gin.RouterGroup, Plugin ...plugin.Plugin) {
	for i := range Plugin {
		PluginGroup := group.Group(Plugin[i].RouterPath())
		Plugin[i].Register(PluginGroup)
	}
}

func InstallPlugin(PublicGroup *gin.RouterGroup, PrivateGroup *gin.RouterGroup) {
	//  添加开放权限的插件 示例
	PluginInit(PublicGroup, example_plugin.ExamplePlugin)
}
