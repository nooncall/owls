package task

import (
	"fmt"
	"github.com/qingfeng777/owls/server/model/common/request"
	"time"
)

// 新建task 模块，用以封装统一的task、审批流

// 1, 设计数据结构
// 2, 权限，也使用新的task
// 3, todo， 改写 审批模块，使用新的task； 实现功能为主，代码优化功能强大后再做。
// 4, 暂时先不管审批task的流程，直接新做一套新的，

// todo, 接口层，先用一个参数指定类型，set 接口filed，再unmarshel

type Task struct {
	ID            int64  `json:"id" gorm:"column:id"`
	Name          string `json:"name" gorm:"column:name"`
	Status        string `json:"status" gorm:"column:status"`
	Creator       string `json:"creator" gorm:"column:creator"`
	Reviewer      string `json:"reviewer" gorm:"column:reviewer"`
	Executor      string `json:"executor" gorm:"column:executor"`
	Description   string `json:"description"`
	ExecInfo      string `json:"exec_info" gorm:"column:exec_info"`
	SubTaskID     int64  `json:"sub_task_id"`
	SubTaskType   string `json:"sub_task_type"`
	RejectContent string `json:"reject_content" gorm:"column:reject_content"`
	Ct            int64  `json:"ct" gorm:"column:ct"`
	Ut            int64  `json:"ut" gorm:"column:ut"`
	Et            int64  `json:"et" gorm:"column:et"`
	Ft            int64  `json:"ft" gorm:"column:ft"`

	SubTask SubTask `json:"sub_task" gorm:"-"`

	StatusName string `json:"status_name" gorm:"-"`
	Action     string `json:"action" gorm:"-"`
}

type SubTask interface {
	AddTask() (int64, error)
	ExecTask() error
	UpdateTask() error
	ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) ([]interface{}, int64, error)
	GetTask(id int64) (interface{}, error)
	GetExecWaitTask() ([]interface{}, int64, error)
}

func AddTask(task *Task) (int64, error) {
	task.Ct = time.Now().Unix()
	subId, err := task.SubTask.AddTask()
	if err != nil{
		return 0, err
	}

	task.SubTaskID, task.Status = subId, WaitApproval
	return taskDao.AddTask(task)
}

func UpdateTask(task *Task) error {
	if err := task.SubTask.UpdateTask(); err != nil {
		return err
	}

	return taskDao.UpdateTask(task)
}

func ListTask(pageInfo request.SortPageInfo, status []string, subTask SubTask) ([]Task, int64, error) {
	tasks, count, err := taskDao.ListTask(pageInfo, true, status)
	if err != nil {
		return nil, 0, err
	}

	for i, v := range tasks {
		if tasks[i].SubTask, err = getSubTask(v, subTask); err != nil {
			return nil, 0, err
		}
	}

	return tasks, count, nil
}

func GetTask(id int64, operator string, subTask SubTask) (*Task, error) {
	task, err := taskDao.GetTask(id)
	if err != nil {
		return nil, err
	}

	task.SubTask, err = getSubTask(*task, subTask)
	return task, nil
}

func getSubTask(task Task, subTask SubTask) (SubTask, error) {
	var err error

	switch task.SubTaskType {
	case Auth:
		subTaskNoType, err := subTask.GetTask(task.SubTaskID)
		if err != nil {
			return nil, err
		}

		if typeTask, ok := subTaskNoType.(SubTask); ok {
			return typeTask, nil
		}
	default:
		return nil, fmt.Errorf("sub task type err: %s", task.SubTaskType)
	}

	return nil, fmt.Errorf("get subTask for %d err: %v", task.SubTaskID, err)
}
