package request

import (
	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
