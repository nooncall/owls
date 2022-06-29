package system

import (
	"github.com/gin-gonic/gin"
	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/common/response"
	"github.com/qingfeng777/owls/server/model/system/request"
	"go.uber.org/zap"
)

type AuthorityBtnApi struct{}

// @Tags AuthorityBtn
// @Summary 获取权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysAuthorityBtnReq true "菜单id, 角色id, 选中的按钮id"
// @Success 200 {object} response.Response{data=response.SysAuthorityBtnRes,msg=string} "返回列表成功"
// @Router /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) GetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	_ = c.ShouldBindJSON(&req)
	if err, res := authorityBtnService.GetAuthorityBtn(req); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithDetailed(res, "查询成功", c)
	}
}

// @Tags AuthorityBtn
// @Summary 设置权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body request.SysAuthorityBtnReq true "菜单id, 角色id, 选中的按钮id"
// @Success 200 {object} response.Response{msg=string} "返回列表成功"
// @Router /authorityBtn/getAuthorityBtn [post]
func (a *AuthorityBtnApi) SetAuthorityBtn(c *gin.Context) {
	var req request.SysAuthorityBtnReq
	_ = c.ShouldBindJSON(&req)
	if err := authorityBtnService.SetAuthorityBtn(req); err != nil {
		global.GVA_LOG.Error("分配失败!", zap.Error(err))
		response.FailWithMessage("分配失败", c)
	} else {
		response.OkWithMessage("分配成功", c)
	}
}

// @Tags AuthorityBtn
// @Summary 设置权限按钮
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Success 200 {object} response.Response{msg=string} "删除成功"
// @Router /authorityBtn/canRemoveAuthorityBtn [post]
func (a *AuthorityBtnApi) CanRemoveAuthorityBtn(c *gin.Context) {
	id := c.Query("id")
	if err := authorityBtnService.CanRemoveAuthorityBtn(id); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}
