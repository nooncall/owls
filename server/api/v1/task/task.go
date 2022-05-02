package task

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/model/common/response"
	"github.com/qingfeng777/owls/server/service/auth/auth"
	"github.com/qingfeng777/owls/server/service/task"
	"github.com/qingfeng777/owls/server/utils"
)

type TaskApi struct{}

func (taskApi *TaskApi) ListReviewTask(ctx *gin.Context) {
	f := "ListReviewTask() -->"
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
	task, count, err := task.ListTask(page, task.ExecStatus())
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListReviewTask err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(response.PageResult{
		List:     task,
		Total:    count,
		PageSize: page.PageSize,
		Page:     page.Page,
	}, ctx)
}

func (taskApi *TaskApi) ListTask(ctx *gin.Context) {
	f := "ListTask() -->"
	var page request.SortPageInfo
	if err := ctx.BindJSON(&page); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	claims, _ := utils.GetClaims(ctx)
	page.Operator = claims.Username
	task, count, err := task.ListTask(page, task.ExecStatus())
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListReviewTask err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(response.PageResult{
		List:     task,
		Total:    count,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, ctx)
}

func (taskApi *TaskApi) ListHistoryTask(ctx *gin.Context) {
	f := "ListHistoryTask() -->"
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
	task, count, err := task.ListTask(page, task.HistoryStatus())
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListHistoryTask err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(response.PageResult{
		List:     task,
		Total:    count,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, ctx)
}

func (taskApi *TaskApi) GetTask(ctx *gin.Context) {
	f := "GetTask() -->"

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, get param failed :%s, id: %s ", f, err.Error(), idStr), ctx)
		return
	}

	task, err := task.GetTask(id, claims.Username)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: get task failed, err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(task, ctx)
}

func (taskApi *TaskApi) UpdateTask(ctx *gin.Context) {
	f := "UpdateTask()-->"
	var taskParam task.Task
	if err := ctx.BindJSON(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	taskParam.Executor = claims.Username
	if err := task.UpdateTask(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, update task failed :%s", f, err.Error()), ctx)
		return
	}
	response.Ok(ctx)
}

func (taskApi *TaskApi) AddTask(ctx *gin.Context) {
	f := "AddTask()-->"

	taskType := ctx.Query("type")
	if taskType == "" {
		response.FailWithMessage("need type param", ctx)
		return
	}

	var taskParam task.Task

	switch taskType {
	case task.Auth:
		taskParam.SubTask = auth.Auth{}
	}

	if err := ctx.BindJSON(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	taskParam.Creator = claims.Username
	id, err := task.AddTask(&taskParam)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, add task failed :%s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(id, ctx)
}