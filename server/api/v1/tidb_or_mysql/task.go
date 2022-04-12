package tidb_or_mysql

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

type TaskApi struct{}

func (taskApi *TaskApi) ListReviewTask(ctx *gin.Context) {
	f := "ListReviewTask() -->"
	var page request.SortPageInfo
	if err := ctx.BindJSON(&page); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s ", f, err.Error()), ctx)
		return
	}

	page.Operator = ctx.MustGet("user").(string)
	task, count, err := task.ListTask(page, task.ExecStatus())
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListReviewTask err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(ListData{
		List:     setTypeNameForTasks(task),
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
	task, count, err := task.ListTask(page, task.SubmitStatus())
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListReviewTask err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(ListData{
		List:     setTypeNameForTasks(task),
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

	page.Operator = ctx.MustGet("user").(string)
	task, count, err := task.ListHistoryTask(page)
	if err != nil {
		response.FailWithMessage(fmt.Sprintf("%s: list ListHistoryTask err: %s", f, err.Error()), ctx)
		return
	}

	response.OkWithData(ListData{
		List:     setTypeNameForTasks(task),
		Total:    count,
		Page:     page.Page,
		PageSize: page.PageSize,
	}, ctx)
}

func (taskApi *TaskApi) GetTask(ctx *gin.Context) {
	f := "GetTask() -->"

	claims, err := utils.GetClaims(ctx)
	if err != nil{
		response.FailWithMessage("get user err: " + err.Error(), ctx)
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

	response.OkWithData(setTypeNameForTask(task), ctx)
}

func (taskApi *TaskApi) UpdateTask(ctx *gin.Context) {
	f := "UpdateTask()-->"
	var taskParam task.OwlTask
	if err := ctx.BindJSON(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil{
		response.FailWithMessage("get user err: " + err.Error(), ctx)
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
	var taskParam task.OwlTask
	if err := ctx.BindJSON(&taskParam); err != nil {
		response.FailWithMessage(fmt.Sprintf("%s, parse param failed :%s", f, err.Error()), ctx)
		return
	}

	if err := task.CheckTaskType(&taskParam); err != nil {
		response.FailWithMessage(err.Error(), ctx)
		return
	}

	claims, err := utils.GetClaims(ctx)
	if err != nil{
		response.FailWithMessage("get user err: " + err.Error(), ctx)
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

// todo, go run this
func ExecWaitTask() {
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("The execWaitTask goroutine panic, err:%s", err)
			// keep the goroutine running
			ExecWaitTask()
		}
	}()

	for {
		<-time.After(time.Minute)

		waitTasks, _, err := task.GetExecWaitTask()
		if err != nil {
			logger.Errorf("the goroutine get exec wait tasks err:%v", err)
		}
		for _, waitTask := range waitTasks {
			countDown := waitTask.Et - time.Now().Unix()
			if countDown <= 0 {
				waitTask.Action = task.Progress
				if err := task.ExecTaskDirectly(&waitTask, &waitTask); err != nil {
					logger.Errorf("while exec task in cron err: %s", err.Error())
				}
			}
		}
	}
}

func setTypeNameForTasks(tasks []task.OwlTask) []task.OwlTask {
	for i, _ := range tasks {
		setTypeNameForTask(&tasks[i])
	}
	return tasks
}

func setTypeNameForTask(oneTask *task.OwlTask) *task.OwlTask {
	for i, v := range oneTask.ExecItems {
		oneTask.ExecItems[i].TaskType = task.GetTypeName(v.TaskType)
	}
	return oneTask
}
