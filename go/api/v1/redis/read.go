package redis

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/api/v1/tidb_or_mysql"
	"github.com/nooncall/owls/go/model/common/response"
	"github.com/nooncall/owls/go/service/redis"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/checker"
)

type ReadApi struct{}

func (readApi *ReadApi) ReadData(ctx *gin.Context) {
	f := "redis.ReadData() -->"

	var req redis.Params
	if err := ctx.BindJSON(&req); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	if req.Cluster == "" {
		response.FailWithMessage(fmt.Sprintf("%s, cluster not selected ", f), ctx)
		return
	}

	response.OkWithData(redis.ExecQuery(ctx.Request.Context(), &req), ctx)
}

func (readApi *ReadApi) ListRule(ctx *gin.Context) {
	f := "redis.ListRule() -->"

	readRules, err := redis.GetReadWhitelist(ctx)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get read white list failed :%s ", f, err.Error()), ctx)
		return
	}
	writeRules, err := redis.GetWriteWhitelist(ctx)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get write white list failed :%s ", f, err.Error()), ctx)
		return
	}

	response.OkWithData(tidb_or_mysql.ListData{
		List: map[string]interface{}{
			"read":  readRules,
			"write": writeRules,
		},
		PageSize: 30, // todo store to db
		Page:     1,
		Total:    30,
	}, ctx)
}

func (ruleApi *ReadApi) UpdateRuleStatus(ctx *gin.Context) {
	f := "redis.UpdateRuleStatus()-->"

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
