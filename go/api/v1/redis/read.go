package redis

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/model/common/response"
	"github.com/nooncall/owls/go/service/redis"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/checker"
	"github.com/nooncall/owls/go/utils"
)

type ReadApi struct{}

func (readApi *ReadApi) ReadData(ctx *gin.Context) {
	f := "redis.ReadData() -->"

	var req redis.Params
	if err := ctx.BindJSON(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(redis.ExecQuery(ctx.Request.Context(), &req), ctx)
}

func (readApi *ReadApi) ListRule(ctx *gin.Context) {
	var pageInfo request.SortPageInfo
	ctx.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	redis.list

	rules, total := redis.ListRules(pageInfo)
	response.OkWithData(ListData{
		List:     rules,
		Total:    int64(total),
		PageSize: pageInfo.PageSize,
		Page:     pageInfo.Page,
	}, ctx)
}

func (ruleApi *ReadApi) UpdateRuleStatus(ctx *gin.Context) {
	f := "UpdateRuleStatus()-->"

	params := struct {
		Name   string `json:"name" binding:"required"`
		Action string `json:"action" binding:"required"`
	}{}
	if err := ctx.BindJSON(&params); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	if err := checker.UpdateRuleStatus(params.Name, params.Action); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,update RuleStatus failed :%s ", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}
