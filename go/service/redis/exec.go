package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.pri.ibanyu.com/middleware/seaweed/xlog"
	pub "gitlab.pri.ibanyu.com/server/redis_manager/pub.git/idl/grpc/redis_manager"
	"gitlab.pri.ibanyu.com/server/redis_manager/service.git/util"
)

// begin write task

// 执行写任务
//1 总体执行，2 指定单个开始执行所有
func ExecWriteTask(ctx context.Context, req *pub.ModifyTaskRequest) error {
	task, err := taskModel.GetTaskById(ctx, req.TaskId)
	if err != nil {
		return fmt.Errorf("get task by id err: %s", err.Error())
	}

	if task.TaskType == TASK_TYPE_APOLLO {
		return ExecWriteApolloTask(ctx, req)
	}

	startId, err := getExecStartId(ctx, req.Modify, task.SubTask, req.SubTask)
	if err != nil {
		return err
	}

	// mean needn't exec task
	if startId < 0 {
		err = taskModel.UpdateTaskStatus(ctx, SysExecSucceed, task.Id)
		if err != nil {
			xlog.Errorf(ctx, "update task status to succeed err, err: %s", err.Error())
		}
		return err
	}

	jump := true
	failed := false
	for _, subTask := range task.SubTask {
		if subTask.Id != startId && jump {
			continue
		}
		jump = false

		//exec task
		resp, err := exec(ctx, subTask.Cmd, task.Service, task.Prefix, task.UsePrefix)
		if err != nil {
			failed = true
			subTask.Status = ItemExecFailed
			xlog.Warnf(ctx, "exec sub task failed, taskId: %d, resp: %v", subTask.Id, resp)
			subTask.ExecResult = fmt.Sprintf("%s", err.Error())
		} else {
			subTask.Status = ItemExecSucceed
			subTask.ExecResult = fmt.Sprintf("%v", resp)
		}

		subTask.ExecuteAt = time.Now().Unix()
		err = subTaskModel.UpdateSubTaskExecCmd(ctx, subTask)
		if err != nil {
			xlog.Errorf(ctx, "after exec ,update subTask err, err: %s", err.Error())
		}

		if failed {
			err := taskModel.UpdateTask(ctx, req, SysExecFailed)
			if err != nil {
				xlog.Errorf(ctx, "update task status to failed err, err: %s", err.Error())
			}
			break
		}
	}

	if !failed {
		err = taskModel.UpdateTask(ctx, req, SysExecSucceed)
		if err != nil {
			xlog.Errorf(ctx, "update task status to succeed err, err: %s", err.Error())
		}
	}
	return err
}

func getExecStartId(ctx context.Context, execType pub.Modify, subTasks []*pub.SubTask, targetSubT *pub.SubTask) (int64, error) {
	// 返回 0 代表出错，eg:execType not found
	switch execType {
	case pub.Modify_Execute:
		// 如果现在的状态不是执行失败或者系统审核成功，直接更新为执行成功状态
		for _, v := range subTasks {
			if v.Status == ItemExecFailed || v.Status == ItemSysJudgeSucceed {
				return v.Id, nil
			}
		}
		return -1, nil
	case pub.Modify_ExecuteBeginAt:
		return targetSubT.Id, nil
	case pub.Modify_ExecuteSkipAt:
		find := false
		for _, v := range subTasks {
			if find {
				return v.Id, nil
			}
			if v.Id == targetSubT.Id {
				find = true
				err := subTaskModel.UpdateSubTaskStatus(ctx, v.Id, ItemExecSkipped)
				if err != nil {
					xlog.Errorf(ctx, "update task status to skip failed, taskId: %d", v.Id)
				}
			}
		}

		//跳过的是最后一个，则不执行
		if find {
			return -1, nil
		} else {
			return 0, fmt.Errorf("execute skip task, target not found, targeId: %d", targetSubT.Id)
		}
	default:
		return 0, fmt.Errorf("execute task err, type not found, type: %d", execType)
	}
}

func ExecReadTask(ctx context.Context, params *Params) (interface{}, error) {
	pass, msg, err := checker.CheckReadCmd(ctx, params.cmd, params.cluster, params.db)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, errors.New(msg)
	}

	return exec(ctx, params.cmd, params.cluster, params.db)
}

func load(ctx context.Context, key interface{}) (value interface{}, err error) {
	return "", nil
}

func exec(ctx context.Context, cmd, cluster string, db int) (resp interface{}, err error) {
	//todo, redis as a param
	redisCli, err := NewRedisCli(cluster, db)
	if err != nil {
		return nil, err
	}

	cmd = util.DelUselessSpace(cmd)
	cmdSplit := strings.Split(cmd, " ")
	if len(cmdSplit) < 2 {
		return nil, fmt.Errorf("while exec cmd err: wrong cmd, cmd: %s", cmd)
	}

	switch strings.ToLower(cmdSplit[0]) {
	case "mset":
		var pairs []interface{}
		for _, v := range cmdSplit[1:] {
			pairs = append(pairs, v)
		}
		cmdResult := redisCli.MSet(ctx, pairs)
		return cmdResult.Val(), cmdResult.Err()
	// multi continuous keys
	case "mget":
		cmdResult := redisCli.Do(ctx, cmdSplit[0], cmdSplit[1:])
		return cmdResult.Val(), cmdResult.Err()
	// one key
	default:
		var othersParams []interface{}
		if len(cmdSplit) >= 3 {
			for _, v := range cmdSplit[2:] {
				othersParams = append(othersParams, v)
			}
		}
		cmdResult := redisCli.Do(ctx, cmdSplit[0], []string{cmdSplit[1]}, othersParams)
		return cmdResult.Val(), cmdResult.Err()
	}

	return resp, err
}

func filterNilStr(data interface{}) string {
	str := fmt.Sprintf("%v", data)
	return strings.ReplaceAll(str, "<nil>", "")
}
