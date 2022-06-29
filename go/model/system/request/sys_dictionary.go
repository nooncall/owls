package request

import (
	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
