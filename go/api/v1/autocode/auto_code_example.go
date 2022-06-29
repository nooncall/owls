package autocode

import (
	"github.com/gin-gonic/gin"
	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/autocode"
	autocodeReq "github.com/qingfeng777/owls/server/model/autocode/request"
	"github.com/qingfeng777/owls/server/model/common/response"
	"github.com/qingfeng777/owls/server/service"
	"github.com/qingfeng777/owls/server/utils"
	"go.uber.org/zap"
)

type AutoCodeExampleApi struct{}

var autoCodeExampleService = service.ServiceGroupApp.AutoCodeServiceGroup.AutoCodeExampleService

// @Tags AutoCodeExample
// @Summary 创建AutoCodeExample
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body autocode.AutoCodeExample true "AutoCodeExample模型"
// @Success 200 {object} response.Response{msg=string} "创建AutoCodeExample"
// @Router /autoCodeExample/createAutoCodeExample [post]
func (autoCodeExampleApi *AutoCodeExampleApi) CreateAutoCodeExample(c *gin.Context) {
	var autoCodeExample autocode.AutoCodeExample
	_ = c.ShouldBindJSON(&autoCodeExample)
	if err := autoCodeExampleService.CreateAutoCodeExample(autoCodeExample); err != nil {
		global.GVA_LOG.Error("创建失败!", zap.Error(err))
		response.FailWithMessage("创建失败", c)
	} else {
		response.OkWithMessage("创建成功", c)
	}
}

// @Tags AutoCodeExample
// @Summary 删除AutoCodeExample
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body autocode.AutoCodeExample true "AutoCodeExample模型"
// @Success 200 {object} response.Response{msg=string} "删除AutoCodeExample"
// @Router /autoCodeExample/deleteAutoCodeExample [delete]
func (autoCodeExampleApi *AutoCodeExampleApi) DeleteAutoCodeExample(c *gin.Context) {
	var autoCodeExample autocode.AutoCodeExample
	_ = c.ShouldBindJSON(&autoCodeExample)
	if err := autoCodeExampleService.DeleteAutoCodeExample(autoCodeExample); err != nil {
		global.GVA_LOG.Error("删除失败!", zap.Error(err))
		response.FailWithMessage("删除失败", c)
	} else {
		response.OkWithMessage("删除成功", c)
	}
}

// @Tags AutoCodeExample
// @Summary 更新AutoCodeExample
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body autocode.AutoCodeExample true "更新AutoCodeExample"
// @Success 200 {object} response.Response{msg=string} "更新AutoCodeExample"
// @Router /autoCodeExample/updateAutoCodeExample [put]
func (autoCodeExampleApi *AutoCodeExampleApi) UpdateAutoCodeExample(c *gin.Context) {
	var autoCodeExample autocode.AutoCodeExample
	_ = c.ShouldBindJSON(&autoCodeExample)
	if err := autoCodeExampleService.UpdateAutoCodeExample(&autoCodeExample); err != nil {
		global.GVA_LOG.Error("更新失败!", zap.Error(err))
		response.FailWithMessage("更新失败", c)
	} else {
		response.OkWithMessage("更新成功", c)
	}
}

// @Tags AutoCodeExample
// @Summary 用id查询AutoCodeExample
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query autocode.AutoCodeExample true "用id查询AutoCodeExample"
// @Success 200 {object} response.Response{data=map[string]interface{},msg=string} "用id查询AutoCodeExample"
// @Router /autoCodeExample/findAutoCodeExample [get]
func (autoCodeExampleApi *AutoCodeExampleApi) FindAutoCodeExample(c *gin.Context) {
	var autoCodeExample autocode.AutoCodeExample
	_ = c.ShouldBindQuery(&autoCodeExample)
	if err := utils.Verify(autoCodeExample, utils.IdVerify); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if err, reAutoCodeExample := autoCodeExampleService.GetAutoCodeExample(autoCodeExample.ID); err != nil {
		global.GVA_LOG.Error("查询失败!", zap.Error(err))
		response.FailWithMessage("查询失败", c)
	} else {
		response.OkWithDetailed(gin.H{"reAutoCodeExample": reAutoCodeExample}, "查询成功", c)
	}
}

// @Tags AutoCodeExample
// @Summary 分页获取AutoCodeExample列表
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data query autocodeReq.AutoCodeExampleSearch true "页码, 每页大小, 搜索条件"
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取AutoCodeExample列表,返回包括列表,总数,页码,每页数量"
// @Router /autoCodeExample/getAutoCodeExampleList [get]
func (autoCodeExampleApi *AutoCodeExampleApi) GetAutoCodeExampleList(c *gin.Context) {
	var pageInfo autocodeReq.AutoCodeExampleSearch
	_ = c.ShouldBindQuery(&pageInfo)
	if err, list, total := autoCodeExampleService.GetAutoCodeExampleInfoList(pageInfo); err != nil {
		global.GVA_LOG.Error("获取失败!", zap.Error(err))
		response.FailWithMessage("获取失败", c)
	} else {
		response.OkWithDetailed(response.PageResult{
			List:     list,
			Total:    total,
			Page:     pageInfo.Page,
			PageSize: pageInfo.PageSize,
		}, "获取成功", c)
	}
}
