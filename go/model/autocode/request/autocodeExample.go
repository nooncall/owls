// 自动生成模板SysDictionaryDetail
package request

import (
	"github.com/qingfeng777/owls/server/model/autocode"
	"github.com/qingfeng777/owls/server/model/common/request"
)

// 如果含有time.Time 请自行import time包
type AutoCodeExampleSearch struct {
	autocode.AutoCodeExample
	request.PageInfo
}
