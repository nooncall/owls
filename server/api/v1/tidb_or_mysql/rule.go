package tidb_or_mysql

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"

	"github.com/gin-gonic/gin"

	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/checker"
)

type RuleApi struct {}

func (ruleApi *RuleApi)LisRule(ctx *gin.Context) {
	rules := checker.ListRules()
	response.OkWithData(ListData{
		List:   rules,
		Total:  int64(len(rules)),
		More:   false,
		Offset: 0,
	}, ctx)
}

func (ruleApi *RuleApi)UpdateRuleStatus(ctx *gin.Context) {
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
