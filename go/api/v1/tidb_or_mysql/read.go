package tidb_or_mysql

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/model/common/response"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/task"
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

	info, err := task.GetTableInfo(&req)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: get table info failed, err: %s", f, err.Error()), ctx)
		return
	}
	response.OkWithData(info, ctx)
}
