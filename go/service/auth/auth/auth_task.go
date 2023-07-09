package auth

import (
	"context"

	"github.com/nooncall/owls/go/service/task"
)

func (a Auth) AddTask(ctx context.Context, cluster string, db int, parentTaskID int64) (int64, bool, error) {
	// a.ParentTaskID = parentTaskID
	count, err := authDao.AddAuth(&a)
	return count, true, err
}

// give auth by set status
func (a Auth) ExecTask(ctx context.Context, startSubTaskId, taskId int64, cluster, db string) error {
	a.Status = StatusPass
	return authDao.UpdateAuth(&a)
}

func (a Auth) UpdateTask(action string) error {
	switch action {
	case task.ActionCancel:
		a.Status = StatusCancelApply
	case task.ActionApproval:
		a.Status = StatusPass
	case task.Reject:
		a.Status = StatusReject
	case task.ActionUpdate:
	}
	return authDao.UpdateAuth(&a)
}

func (a Auth) ListTask(parentTaskID int64) (interface{}, error) {
	panic("implement me")
}

func (a Auth) GetTask(id int64) (interface{}, error) {
	return authDao.GetAuth(id)
}

func (a Auth) GetExecWaitTask() ([]interface{}, int64, error) {
	panic("implement me")
}
