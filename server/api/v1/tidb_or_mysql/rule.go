package tidb_or_mysql

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/model/common/response"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/checker"
	"github.com/qingfeng777/owls/server/utils"
)

type RuleApi struct{}

func (ruleApi *RuleApi) LisRule(ctx *gin.Context) {
	var pageInfo request.SortPageInfo
	ctx.ShouldBindJSON(&pageInfo)
	if err := utils.Verify(pageInfo.PageInfo, utils.PageInfoVerify); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	rules, total := checker.ListRules(pageInfo)
	response.OkWithData(ListData{
		List:     rules,
		Total:    int64(total),
		PageSize: pageInfo.PageSize,
		Page:     pageInfo.Page,
	}, ctx)
}

func (ruleApi *RuleApi) UpdateRuleStatus(ctx *gin.Context) {
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
