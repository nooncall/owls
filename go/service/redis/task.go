package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/utils/logger"
	"gorm.io/gorm"
)

type SubTask interface {
	AddTask() (int64, error)
	ExecTask() error
	UpdateTask(action string) error
	ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) ([]interface{}, int64, error)
	GetTask(id int64) (interface{}, error)
	GetExecWaitTask() ([]interface{}, int64, error)
}
type SubTask2 interface { // need
	AddTask() (int64, error)
	ExecTask(taskId int64) error
	UpdateTask(action string) error
	ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) ([]interface{}, int64, error)
	GetTask(id int64) (interface{}, error)
	GetExecWaitTask() ([]interface{}, int64, error)
}

type RedisTask struct {
	ID       int64  `json:"id" gorm:"column:id"`
	TaskID   int64  `json:"task_id" gorm:"column:task_id"`
	Cmd      string `json:"cmd" gorm:"column:cmd"`
	ExecInfo string `json:"exec_info" gorm:"column:exec_info"`
	Cluster  string `json:"cluster"`
	DB       string `json:"db"`
}

type RedisTaskDao interface {
	AddTask(db *gorm.DB, task *RedisTask) (int64, error)
	UpdateTask(db *gorm.DB, task *RedisTask) error
	ListRedisTaskByTaskID(db *gorm.DB, id int64) ([]RedisTask, error)
	ListTask(db *gorm.DB, pageInfo request.SortPageInfo, isDBA bool, status []string) ([]RedisTask, int64, error)
	GetTask(db *gorm.DB, id int64) (*RedisTask, error)
}

var redisTaskDao RedisTaskDao

func SetRedisTaskDao(impl RedisTaskDao) {
	redisTaskDao = impl
}

func (r *RedisTask) AddTask() (int64, error) {
	// split and store， others as normal
	cmds := strings.Split(r.Cmd, ";")
	tx := GetDB().Begin()

	for _, v := range cmds {
		if v == "" {
			continue
		}

		if _, err := redisTaskDao.AddTask(tx, &RedisTask{
			ID:     r.ID,
			TaskID: r.TaskID,
			Cmd:    v,
		}); err != nil {
			tx.Rollback()
			return 0, err
		}

	}

	return 0, tx.Commit().Error
}

func (r *RedisTask) ExecTask(ctx context.Context, taskId int64) error {
	tasks, err := redisTaskDao.ListRedisTaskByTaskID(GetDB(), taskId)
	if err != nil {
		return fmt.Errorf("while exec task, get task err: %v", err)
	}

	intDB, err := strconv.Atoi(r.DB)

	// exec, // 假设都是独立的，更常见的场景。
	var failed bool
	for _, v := range tasks {
		resp, err := exec(ctx, v.Cmd, r.Cluster, intDB)
		if err != nil {
			failed = true
			logger.Warnf("exec redis task failed, taskId: %d, resp: %v, err: %v", v.ID, resp, err)
		}
		v.ExecInfo = fmt.Sprintf("%v", resp)
		if err := v.UpdateTask(""); err != nil {
			return err
		}
	}
	if failed {
		return errors.New("exec failed")
	}
	return nil
}

func (r *RedisTask) UpdateTask(action string) error {
	return redisTaskDao.UpdateTask(GetDB(), r)
}

func (r *RedisTask) ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) (interface{}, int64, error) {
	return redisTaskDao.ListTask(GetDB(), pageInfo, isDBA, status)
}
func (r *RedisTask) GetTask(id int64) (interface{}, error) {
	return redisTaskDao.GetTask(GetDB(), id)
}
