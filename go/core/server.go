package core

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/initialize"
	"github.com/nooncall/owls/go/service/system"
	"github.com/nooncall/owls/go/utils/logger"
)

type server interface {
	ListenAndServe() error
}

func RunWindowsServer() {
	// 从db加载jwt数据
	if global.GVA_DB != nil {
		system.LoadAll()
	}

	Router := gin.Default()

	initialize.Routers(Router)
	Router.Static("/form-generator", "./resource/page")
	frontRouter(Router)
	docRouter(Router)

	address := fmt.Sprintf(":%d", global.GVA_CONFIG.System.Addr)
	s := initServer(address, Router)
	// 保证文本顺序输出
	// In order to ensure that the text order output can be deleted
	time.Sleep(10 * time.Microsecond)
	global.GVA_LOG.Info("server run success on ", zap.String("address", address))

	fmt.Printf(`
	欢迎使用 github.com/nooncall/owls/go
	当前版本:V0.0.1
	加群方式:微信号：xxx QQ群：xxxx
	默认自动化文档地址:http://127.0.0.1%s/swagger/index.html
	默认前端文件运行地址:http://127.0.0.1:8080
	如果项目让您获得了收益，那就帮忙宣传一下吧！
`, address)
	global.GVA_LOG.Error(s.ListenAndServe().Error())
}

func docRouter(r *gin.Engine) {
	currentDir := getCurrentDirectory()
	// currentDir = "/Users/mingbai/openS/vue/database-manager/bin"
	logger.Info("doc router current dir is: ", currentDir)
	r.Static("/docs", filepath.Join(currentDir, "./docs-static"))

	r.GET("/doc", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "docs/")
	})

	r.NoRoute(func(c *gin.Context) {
		if !strings.Contains(c.Request.RequestURI, "/api") {
			path := strings.Split(c.Request.URL.Path, "/")
			if len(path) > 1 {
				c.File(filepath.Join(currentDir, "./docs-static") + "/index.html")
				return
			}
		}
	})
}

func frontRouter(r *gin.Engine) {
	currentDir := getCurrentDirectory()
	// currentDir = "/Users/mingbai/openS/vue/database-manager/bin"
	logger.Info("current dir is: ", currentDir)
	r.Static("/owls", filepath.Join(currentDir, "./static"))

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "owls/")
	})

	r.NoRoute(func(c *gin.Context) {
		if !strings.Contains(c.Request.RequestURI, "/api") {
			path := strings.Split(c.Request.URL.Path, "/")
			if len(path) > 1 {
				c.File(filepath.Join(currentDir, "./static") + "/index.html")
				return
			}
		}
	})
}

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		logger.Errorf("get current dir err: %s", err.Error())
		return ""
	}
	return strings.Replace(dir, "\\", "/", -1)
}
