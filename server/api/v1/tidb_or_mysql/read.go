package tidb_or_mysql

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/qingfeng777/owls/server/model/common/response"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/task"
)

type ReadApi struct{}

func (readApi *ReadApi) ReadData(ctx *gin.Context) {
	f := "ReadData() -->"

	var req task.SqlParam
	if err := ctx.BindJSON(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	rollBackData, err := task.ReadData(&req)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: read failed, err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(rollBackData, ctx)
}

func (readApi *ReadApi) GetTableInfo(ctx *gin.Context) {
	f := "GetTableInfo()-->"
	var req task.SqlParam
	if err := ctx.BindJSON(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}
	
	if err := task.GetTableInfo(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: get table info failed, err: %s", f, err.Error()), ctx)
		return
	}
	response.Ok(ctx)
}
