package auth

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/model/common/response"
	"github.com/nooncall/owls/go/service/auth/auth"
	"github.com/nooncall/owls/go/utils"
)

// list data type

// 列表，删除，添加

type AuthApi struct{}

func (AuthApi *AuthApi) ListAuth(ctx *gin.Context) {
	f := "ListHistoryAuth() -->"

	var page request.SortPageInfo
	if err := ctx.BindJSON(&page); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	page.Operator = claims.Username
	auth, count, err := auth.ListAuth(page)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListHistoryAuth err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(response.PageResult{
		List:     auth,
		Total:    count,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, ctx)
}

func (AuthApi *AuthApi) ListDataType(ctx *gin.Context) {
	types := auth.ListDataType()
	response.OkWithData(response.PageResult{
		List:     types,
		Total:    int64(len(types)),
		Page:     0,
		PageSize: 20,
	}, ctx)
}

func (AuthApi *AuthApi) DelAuth(ctx *gin.Context) {
	f := "DelAuth() -->"

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr), ctx)
		return
	}

	err = auth.DelAuth(id)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: del Auth failed, err: %s", f, err.Error()), ctx)
		return
	}

	response.Ok(ctx)
}

func (AuthApi *AuthApi) AddAuth(ctx *gin.Context) {
	f := "AddAuth()-->"

	var AuthParam auth.Auth
	if err := ctx.BindJSON(&AuthParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	id, err := auth.AddAuth(&AuthParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, add Auth failed :%s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(id, ctx)
}
