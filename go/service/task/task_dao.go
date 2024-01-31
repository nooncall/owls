package task

import (
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/utils"
	"gorm.io/gorm"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/admin"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/task"
)

type taskDaoImpl struct {
}

var taskDao taskDaoImpl

func GetDB() *gorm.DB {
	// todo, refactor to config
	return global.GVA_DB.Debug()
}

func (taskDaoImpl) AddTask(task *Task) (int64, error) {
	tx := GetDB().Begin()
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return task.ID, tx.Commit().Error
}

func (taskDaoImpl) UpdateTask(task *Task) error {
	return GetDB().Model(task).Where("id = ?", task.ID).Updates(task).Error
}

func (taskDaoImpl) ListTask(info request.SortPageInfo, isDBA bool, status []task.ItemStatus, subType string) ([]Task, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := GetDB().Offset(offset)
	db.Where("sub_task_type = ?", subType)
	if info.Key != "" {
		fmtKey := "%" + info.Key + "%"
		db = db.Where("id like ? or name like ? or status like ? or creator like ?",
			fmtKey, fmtKey, fmtKey, fmtKey)
	}
	db = db.Where("status in (?)", status)

	// check admin
	isAdmin, err := admin.IsAdmin(info.Operator)
	if err != nil {
		return nil, 0, err
	}
	if !isAdmin {
		db.Where("creator = ?", info.Operator)
	}

	var count int64
	if err := db.Model(&Task{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db.Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(utils.GenerateOrderField(info.OrderKey, info.Desc))
	} else {
		db = db.Order("ct desc")
	}

	var tasks []Task
	if err := db.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

func getTaskExecItems(db *gorm.DB, taskP *Task) ([]task.OwlExecItem, []task.OwlSubtask, error) {
	var formattedItems []task.OwlExecItem
	var subTasks []task.OwlSubtask
	if err := db.Find(&subTasks, "task_id = ?", taskP.ID).Error; err != nil {
		return nil, nil, err
	}

	for idx, subTask := range subTasks {
		var items []task.OwlExecItem
		if err := db.Find(&items, "subtask_id = ?", subTask.ID).Error; err != nil {
			return nil, nil, err
		}
		subTasks[idx].ExecItems = items

		for _, v := range items {
			v.DBName = subTask.DbName
			v.ClusterName = subTask.ClusterName
			v.TaskType = subTask.TaskType
			formattedItems = append(formattedItems, v)
		}
	}

	return formattedItems, subTasks, nil
}

func (taskDaoImpl) GetTask(id int64) (*Task, error) {
	var task Task
	return &task, GetDB().First(&task, "id = ?", id).Error
}

const ExecWaitTasksLimit = 5

func (taskDaoImpl) GetExecWaitTask() ([]Task, int64, error) {
	condition := "status in (?)"

	var count int64
	if err := GetDB().Model(&Task{}).Where(condition,
		task.ExecWait).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var tasks []Task
	if err := GetDB().Order("et asc").Limit(ExecWaitTasksLimit).
		Find(&tasks, condition, task.ExecWait).Error; err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}
