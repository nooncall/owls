package system

import (
	"github.com/qingfeng777/owls/server/config"
)

// 配置文件结构体
type System struct {
	Config config.Server `json:"config"`
}
