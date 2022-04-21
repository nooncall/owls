package request

import (
	"github.com/qingfeng777/owls/server/model/autocode"
	"github.com/qingfeng777/owls/server/model/common/request"
)

type {{.StructName}}Search struct{
    autocode.{{.StructName}}
    request.PageInfo
}