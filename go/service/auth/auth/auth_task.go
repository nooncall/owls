package auth

import (
	"context"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/service/task"
)

func (a Auth) AddTask() (int64, error) {
	return authDao.AddAuth(&a)
}

// give auth by set status
func (a Auth) ExecTask(ctx context.Context, taskId int64) error {
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

func (a Auth) ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) (interface{}, int64, error) {
	panic("implement me")
}

func (a Auth) GetTask(id int64) (interface{}, error) {
	return authDao.GetAuth(id)
}

func (a Auth) GetExecWaitTask() ([]interface{}, int64, error) {
	panic("implement me")
}
