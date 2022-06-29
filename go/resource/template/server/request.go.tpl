package request

import (
	"github.com/nooncall/owls/go/model/autocode"
	"github.com/nooncall/owls/go/model/common/request"
)

type {{.StructName}}Search struct{
    autocode.{{.StructName}}
    request.PageInfo
}