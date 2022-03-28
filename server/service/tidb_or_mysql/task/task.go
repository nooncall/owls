package task

import (
	"errors"
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql"
	"strings"
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/sql_util"
)

type OwlTask struct {
	ID            int64  `json:"id" gorm:"column:id"`
	Name          string `json:"name" gorm:"column:name"`
	Status        string `json:"status" gorm:"column:status"`
	Creator       string `json:"creator" gorm:"column:creator"`
	Reviewer      string `json:"reviewer" gorm:"column:reviewer"`
	Executor      string `json:"executor" gorm:"column:executor"`
	ExecInfo      string `json:"exec_info" gorm:"column:exec_info"`
	RejectContent string `json:"reject_content" gorm:"column:reject_content"`
	Ct            int64  `json:"ct" gorm:"column:ct"`
	Ut            int64  `json:"ut" gorm:"column:ut"`
	Et            int64  `json:"et" gorm:"column:et"`
	Ft            int64  `json:"ft" gorm:"column:ft"`

	SubTasks  []OwlSubtask  `json:"sub_tasks" gorm:"-"`
	ExecItems []OwlExecItem `json:"exec_items" gorm:"-"`
	ExecItem  *OwlExecItem  `json:"exec_item" gorm:"-"`
	EditAuth  *EditAuth     `json:"edit_auth" gorm:"-"`

	StatusName string `json:"status_name" gorm:"-"`
	Action     Action `json:"action" gorm:"-"`
}

type TaskDao interface {
	AddTask(task *OwlTask) (int64, error)
	UpdateTask(task *OwlTask) error
	ListTask(pagination *tidb_or_mysql.Pagination, isDBA bool, status []ItemStatus) ([]OwlTask, int64, error)
	ListHistoryTask(page *tidb_or_mysql.Pagination, isDBA bool) ([]OwlTask, int64, error)
	GetTask(id int64) (*OwlTask, error)
	GetExecWaitTask() ([]OwlTask, int, error)
}

var taskDao TaskDao

func SetTaskDao(impl TaskDao) {
	taskDao = impl
}

type SubTaskDao interface {
	UpdateItem(item *OwlExecItem) error
	DelItem(item *OwlExecItem) error
	UpdateItemByBackupId(item *OwlExecItem) error
}

var subTaskDao SubTaskDao

func SetSubTaskDao(impl SubTaskDao) {
	subTaskDao = impl
}

func AddTask(task *OwlTask) (int64, error) {
	reviewer, err := AuthTool.GetReviewer(task.Creator)
	if err != nil {
		return 0, err
	}
	task.Reviewer = reviewer

	// 拆分
	for idx, subTask := range task.SubTasks {
		newSubTask, err := splitSubTaskExecItems(&subTask)
		if err != nil {
			return 0, err
		}
		task.SubTasks[idx] = *newSubTask
	}

	if err := checkExecItemNum(task); err != nil {
		return 0, err
	}

	// check sql
	checkPass := true
	for idx, subTask := range task.SubTasks {
		dbInfo, err := dbTool.GetDBConn(subTask.DbName, subTask.ClusterName)
		if err != nil {
			return 0, err
		}

		for itemIdx, item := range subTask.ExecItems {
			pass, suggestion, affectRow := checker.SqlCheck(item.SQLContent, "", "", dbInfo)
			if affectRow > 0 {
				task.SubTasks[idx].ExecItems[itemIdx].AffectRows = affectRow
			}

			task.SubTasks[idx].ExecItems[itemIdx].Status = ItemCheckPass
			if !pass {
				checkPass = false
				task.SubTasks[idx].ExecItems[itemIdx].RuleComments = suggestion
				task.SubTasks[idx].ExecItems[itemIdx].Status = ItemCheckFailed
			}
		}

		dbInfo.CloseConn()
	}

	if checkPass {
		task.Status = CheckPass
	} else {
		task.Status = CheckFailed
	}
	task.Ct = time.Now().Unix()

	// add task
	return taskDao.AddTask(task)
}

func splitSubTaskExecItems(subTask *OwlSubtask) (*OwlSubtask, error) {
	var items []OwlExecItem
	for _, execItem := range subTask.ExecItems {
		sqls, err := sql_util.SplitMultiSql(execItem.SQLContent)
		if err != nil {
			return nil, err
		}
		for _, v := range sqls {
			newItem := execItem
			newItem.SQLContent = v
			items = append(items, newItem)
		}
	}
	subTask.ExecItems = items
	return subTask, nil
}

func checkExecItemNum(task *OwlTask) error {
	num := 0
	for _, v := range task.SubTasks {
		for range v.ExecItems {
			num++
		}
	}
	if num > config.Conf.Server.NumOnceLimit {
		return fmt.Errorf("exec too many sql once, sqlNum: %d, limit: %d", num, config.Conf.Server.NumOnceLimit)
	}
	return nil
}

func refreshTaskStatus(taskId int64, et, ft int64, executor, execInfo string) error {
	task, err := taskDao.GetTask(taskId)
	if err != nil {
		return err
	}

	status := Executing
	if task.ExecItems[len(task.ExecItems)-1].Status != ItemCheckPass {
		failed := false
		for _, v := range task.ExecItems {
			if v.Status != ItemSuccess {
				status = ExecFailed
				failed = true
			}
		}
		if !failed {
			status = ExecSuccess
		}
	}

	return taskDao.UpdateTask(&OwlTask{ID: taskId, Status: status, Et: et, Ft: ft, Executor: executor, Ut: time.Now().Unix(), ExecInfo: execInfo})
}

func UpdateTask(task *OwlTask) error {
	dbTask, err := taskDao.GetTask(task.ID)
	if err != nil {
		return err
	}

	isReviewer := strings.Contains(dbTask.Reviewer, task.Executor)
	isDba, err := AuthTool.IsDba(task.Executor)
	if err != nil {
		return err
	}

	//对于执行变更，检查权限
	if task.Action == BeginAt || task.Action == SkipAt || (dbTask.Status == DBAPass && task.Action == Progress) {
		if !(isReviewer && allIsDmlTask(task)) && !isDba {
			return errors.New("auth invalid")
		}
	}

	task.Ut = time.Now().Unix()

	switch task.Action {
	case EditItem:
		if err := subTaskDao.UpdateItem(task.ExecItem); err != nil {
			return err
		}
		return recheckTask(task.ID, task.Executor)
	case DelItem:
		if err := subTaskDao.DelItem(task.ExecItem); err != nil {
			return err
		}
		return recheckTask(task.ID, task.Executor)
	case DoCancel:
		return doCancel(task, dbTask, isDba)
	case SkipAt, BeginAt:
		return Exec(task, dbTask)
	case Progress:
		return ProgressEdit(task, dbTask)
	case DoReject:
		return taskDao.UpdateTask(&OwlTask{
			ID:            task.ID,
			Status:        Reject,
			Executor:      task.Executor,
			RejectContent: task.RejectContent,
			Ut:            time.Now().Unix(),
		})
	default:
		return fmt.Errorf("action type not found, action: %s", task.Action)
	}
}

func recheckTask(id int64, operator string) error {
	task, err := taskDao.GetTask(id)
	if err != nil {
		return err
	}

	return checkTask(task)
}

func checkTask(task *OwlTask) error {
	checkPass := true
	for _, subTask := range task.SubTasks {
		dbInfo, err := dbTool.GetDBConn(subTask.DbName, subTask.ClusterName)
		if err != nil {
			return err
		}

		for _, item := range subTask.ExecItems {
			pass, suggestion, affectRow := checker.SqlCheck(item.SQLContent, "", "", dbInfo)
			if affectRow > 0 {
				item.AffectRows = affectRow
			}
			if !pass {
				checkPass = false
				item.Status = ItemCheckFailed
				item.RuleComments = suggestion
			} else {
				item.Status = ItemCheckPass
				item.RuleComments = " "
			}

			if err := subTaskDao.UpdateItem(&item); err != nil {
				return err
			}
		}

		dbInfo.CloseConn()
	}

	if checkPass {
		task.Status = CheckPass
	} else {
		task.Status = CheckFailed
	}

	return taskDao.UpdateTask(task)
}

func doCancel(task, dbTask *OwlTask, isDba bool) error {
	switch {
	case dbTask.Creator == task.Executor:
		task.Status = Cancel
	case isDba:
		task.Status = ExecCancel
	default:
		return errors.New("no auth to do cancel")
	}

	return taskDao.UpdateTask(&OwlTask{
		ID:       task.ID,
		Status:   task.Status,
		Executor: task.Executor,
	})
}

func ProgressEdit(task, dbTask *OwlTask) error {
	switch dbTask.Status {
	case CheckPass:
		task.Status = ReviewPass
	case DBAPass, ReviewPass, ExecWait:
		return Exec(task, dbTask)
	default:
		// new cron task
		if task.Et > time.Now().Unix() {
			return Exec(task, dbTask)
		}
		return fmt.Errorf("progress failed, task status invalid, status: %s", dbTask.Status)
	}

	return taskDao.UpdateTask(&OwlTask{
		ID:       task.ID,
		Status:   task.Status,
		Ut:       time.Now().Unix(),
		Executor: task.Executor,
	})
}

func ListTask(pagination *tidb_or_mysql.Pagination, status []ItemStatus) ([]OwlTask, int64, error) {
	isDba, err := AuthTool.IsDba(pagination.Operator)
	if err != nil {
		return nil, 0, err
	}
	tasks, count, err := taskDao.ListTask(pagination, isDba, status)
	if err != nil {
		return nil, 0, err
	}

	for i, v := range tasks {
		tasks[i].StatusName = StatusName(v.Status)
		tasks[i].EditAuth = GetTaskOperateAuth(false, v.Creator == pagination.Operator, strings.Contains(v.Reviewer, pagination.Operator), isDba, &v)
	}

	return tasks, count, nil
}

func ListHistoryTask(pagination *tidb_or_mysql.Pagination) ([]OwlTask, int64, error) {
	isDba, err := AuthTool.IsDba(pagination.Operator)
	if err != nil {
		return nil, 0, err
	}
	tasks, count, err := taskDao.ListHistoryTask(pagination, isDba)
	if err != nil {
		return nil, 0, err
	}

	for i, v := range tasks {
		tasks[i].StatusName = StatusName(v.Status)
		//tasks[i].EditAuth = GetTaskOperateAuth(false, v.Creator == pagination.Operator, strings.Contains(v.Reviewer, pagination.Operator), isDba, &v)
	}

	return tasks, count, nil
}

func GetTask(id int64, operator string) (*OwlTask, error) {
	isDba, err := AuthTool.IsDba(operator)
	if err != nil {
		return nil, err
	}

	task, err := taskDao.GetTask(id)
	if err != nil {
		return nil, err
	}
	task.SubTasks = nil
	task.StatusName = StatusName(task.Status)

	//task.ExecItems = fmtExecItemFromOneTask(task)
	task.EditAuth = GetTaskOperateAuth(true, operator == task.Creator, strings.Contains(task.Reviewer, operator), isDba, task)
	return task, nil
}

func GetExecWaitTask() ([]OwlTask, int, error) {
	tasks, count, err := taskDao.GetExecWaitTask()
	if err != nil {
		return nil, 0, err
	}

	for i, v := range tasks {
		tasks[i].StatusName = StatusName(v.Status)
	}

	return tasks, count, nil
}

func CheckTaskType(task *OwlTask) error {
	for _, subTask := range task.SubTasks {
		taskType := subTask.TaskType
		for _, execItem := range subTask.ExecItems {
			if err := checkTaskType(execItem.SQLContent, taskType); err != nil {
				return err
			}
		}
	}
	return nil
}

func checkTaskType(sql string, taskType TaskType) error {
	sql = strings.TrimSpace(sql)
	parts := strings.SplitN(sql, " ", 2)
	if len(parts) == 0 {
		return fmt.Errorf("sql error")
	}

	var curType TaskType
	switch strings.ToLower(parts[0]) {
	case "create":
		// create index、create unique index 属于改表
		parts[1] = strings.TrimSpace(parts[1])
		if strings.EqualFold(strings.SplitN(parts[1], " ", 2)[0], "index") ||
			strings.EqualFold(strings.SplitN(parts[1], " ", 2)[0], "unique") {
			curType = DDLUpdate
		} else {
			curType = DDLCreate
		}
	case "insert", "delete", "update":
		curType = DML
	default:
		curType = DDLUpdate
	}
	if curType != taskType {
		return fmt.Errorf("task type error")
	}
	return nil
}
