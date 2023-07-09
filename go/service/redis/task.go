package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/service/task"
	"github.com/nooncall/owls/go/utils/logger"
	"gorm.io/gorm"
)

type RedisTask struct {
	ID        int64  `json:"id" gorm:"column:id"`
	TaskID    int64  `json:"task_id" gorm:"column:task_id"`
	Cmd       string `json:"cmd" gorm:"column:cmd"`
	ExecInfo  string `json:"exec_info" gorm:"column:exec_info"`
	CheckInfo string `json:"check_info" gorm:"column:check_info"`
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

func (r *RedisTask) AddTask(ctx context.Context, cluster string, db int, parentTaskID int64) (int64, bool, error) {
	// split and store， others as normal
	cmds := strings.Split(r.Cmd, ";")
	tx := GetDB().Begin()
	checkPass := true

	for _, v := range cmds {
		if v == "" {
			continue
		}

		pass, checkMsg, err := checker.CheckWriteCmd(ctx, v, cluster, db)
		if err != nil {
			return 0, checkPass, err
		}
		if !pass {
			checkPass = pass
		}

		if _, err := redisTaskDao.AddTask(tx, &RedisTask{
			TaskID:    parentTaskID,
			CheckInfo: checkMsg,
			Cmd:       v,
		}); err != nil {
			tx.Rollback()
			return 0, checkPass, err
		}
	}

	return 0, checkPass, tx.Commit().Error
}

func (r *RedisTask) ExecTask(ctx context.Context, startSubTaskId, taskId int64, cluster, db string) error {
	tasks, err := redisTaskDao.ListRedisTaskByTaskID(GetDB(), taskId)
	if err != nil {
		return fmt.Errorf("while exec task, get task err: %v", err)
	}

	database, err := strconv.Atoi(db)
	if err != nil {
		return err
	}

	// exec, 假设都是独立的，更常见的场景。
	start := false
	if startSubTaskId <= 0 {
		start = true
	}

	var failed bool
	for _, v := range tasks {
		if v.ID == startSubTaskId {
			start = true
		}
		if !start {
			continue
		}

		resp, err := exec(ctx, v.Cmd, cluster, database)
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

func (r *RedisTask) ListTask(parentTaskID int64) (interface{}, error) {
	rTasks, err := redisTaskDao.ListRedisTaskByTaskID(GetDB(), parentTaskID)
	if err != nil {
		return nil, err
	}

	var iTasks []task.SubTask
	for i := range rTasks {
		iTasks = append(iTasks, task.SubTask(&rTasks[i]))
	}
	return iTasks, nil
}

func (r *RedisTask) ListPageTask(pageInfo request.SortPageInfo, isDBA bool, status []string) (interface{}, int64, error) {
	return redisTaskDao.ListTask(GetDB(), pageInfo, isDBA, status)
}
func (r *RedisTask) GetTask(id int64) (interface{}, error) {
	return redisTaskDao.GetTask(GetDB(), id)
}
