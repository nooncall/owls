package request

import (
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/model/system"
)

type SysOperationRecordSearch struct {
	system.SysOperationRecord
	request.PageInfo
}
