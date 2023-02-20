package redis

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/model/common/response"
	"github.com/nooncall/owls/go/service/redis"
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
