package task

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/model/common/response"
	request2 "github.com/qingfeng777/owls/server/model/system/request"
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

	subTask, _, err := genSubType(ctx)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	page.Operator = claims.Username
	task, count, err := task.ListTask(page, task.ExecStatus(), subTask)
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

	subTask, _, err := genSubType(ctx)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	claims, _ := utils.GetClaims(ctx)
	page.Operator = claims.Username
	task, count, err := task.ListTask(page, task.ExecStatus(), subTask)
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

	subTask, _, err := genSubType(ctx)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}

	page.Operator = claims.Username
	task, count, err := task.ListTask(page, task.HistoryStatus(), subTask)
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

	subTask, _, err := genSubType(ctx)
	if err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	task, err := task.GetTask(id, claims.Username, subTask)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: get task failed, err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(task, ctx)
}

func (taskApi *TaskApi) UpdateTask(ctx *gin.Context) {
	f := "UpdateTask()-->"

	var taskParam TaskParam
	if err := ctx.BindJSON(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}
	fillSubTask(&taskParam, nil)

	if err := task.UpdateTask(&taskParam.Task); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, update task failed :%s", f, err.Error()), ctx)
		return
	}
	response.Ok(ctx)
}

func (taskApi *TaskApi) AddTask(ctx *gin.Context) {
	f := "AddTask()-->"

	var taskParam TaskParam
	if err := ctx.BindJSON(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil {
		response.FailWithMessage("get user err: "+err.Error(), ctx)
		return
	}
	fillSubTask(&taskParam, claims)

	id, err := task.AddTask(&taskParam.Task)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, add task failed :%s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(id, ctx)
}

func genSubType(ctx *gin.Context) (task.SubTask, string, error) {
	taskType := ctx.Query("type")
	if taskType == "" {
		return nil, "", errors.New("need type param")
	}

	switch taskType {
	case task.Auth:
		return auth.Auth{}, task.Auth, nil
	default:
		return nil, "", fmt.Errorf("sub task type not found: %s", taskType)
	}
}

type TaskParam struct {
	Task task.Task `json:"task"`
	Auth auth.Auth `json:"auth"`
}

func fillSubTask(param *TaskParam, claims *request2.CustomClaims) {
	switch param.Task.SubTaskType {
	case task.Auth:
		if claims != nil {
			param.Auth.UserId = claims.ID
			param.Auth.Username = claims.Username
		}
		param.Task.SubTask = param.Auth
	}

	if claims != nil{
		param.Task.Creator = claims.Username
	}
}
