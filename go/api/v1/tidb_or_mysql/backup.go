package tidb_or_mysql

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/model/common/response"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/task"
	"github.com/nooncall/owls/go/utils"
)

type BackupApi struct{}

func (backupApi *BackupApi) ListRollbackData(ctx *gin.Context) {
	f := "GetTask() -->"

	var req task.SqlParam
	if err := ctx.BindJSON(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	rollBackData, err := task.ListRollbackData(&req)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: get rollBackData failed, err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(rollBackData, ctx)
}

func (backupApi *BackupApi) Rollback(ctx *gin.Context) {
	f := "Rollback()-->"
	var req task.SqlParam
	if err := ctx.BindJSON(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	req.Executor = claims.Username
	if err := task.Rollback(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: rollback failed, err: %s", f, err.Error()), ctx)
		return
	}
	response.Ok(ctx)
}
