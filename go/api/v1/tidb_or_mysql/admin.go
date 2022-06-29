package tidb_or_mysql

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/nooncall/owls/go/model/common/response"
	service "github.com/nooncall/owls/go/service/tidb_or_mysql"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/admin"
)

type AdminApi struct{}

func (adminApi *AdminApi) AddAdmin(ctx *gin.Context) {
	f := "AddAdmin()-->"
	var adminParam admin.OwlAdmin
	if err := ctx.BindJSON(&adminParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	adminParam.Creator = ctx.MustGet("user").(string)
	id, err := admin.AddAdmin(&adminParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, add admin failed :%s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(id, ctx)
}

func (adminApi *AdminApi) ListAdmin(ctx *gin.Context) {
	f := "ListAdmin() -->"
	var page service.Pagination
	if err := ctx.BindJSON(&page); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	page.Operator = ctx.MustGet("user").(string)
	admins, count, err := admin.ListAdmin(&page)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list admin err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(ListData{
		List:   admins,
		Total:  count,
		More:   count > int64(page.Offset+page.Limit),
		Offset: page.Offset,
	}, ctx)
}

func (adminApi *AdminApi) DelAdmin(ctx *gin.Context) {
	f := "DelAdmin()-->"

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr), ctx)
		return
	}

	if err := admin.DelAdmin(id); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s,del admin failed :%s ", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}
