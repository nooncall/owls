package dao

import (
	"gorm.io/gorm"

	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/task"
)

type TaskDaoImpl struct {
}

var Task TaskDaoImpl

func (TaskDaoImpl) AddTask(task *task.OwlTask) (int64, error) {
	tx := GetDB().Begin()
	if err := tx.Create(task).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	for _, subTask := range task.SubTasks {
		subTask.TaskID = task.ID
		if err := tx.Create(&subTask).Error; err != nil {
			tx.Rollback()
			return 0, err
		}

		for _, item := range subTask.ExecItems {
			item.SubtaskID = subTask.ID
			item.TaskID = task.ID
			if err := tx.Create(&item).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
		}

	}

	return task.ID, tx.Commit().Error
}

func (TaskDaoImpl) UpdateTask(task *task.OwlTask) error {
	return GetDB().Model(task).Where("id = ?", task.ID).Updates(task).Error
}

func (TaskDaoImpl) ListTask(info request.SortPageInfo, isDBA bool, status []task.ItemStatus) ([]task.OwlTask, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := GetDB().Offset(offset)
	if info.Key != "" {
		fmtKey := "%" + info.Key + "%"
		db = db.Where("id like ? or name like ? or status like ? or creator like ?",
			fmtKey, fmtKey, fmtKey, fmtKey)
	}
	db = db.Where("status in (?) and creator = ?", status, info.Operator)

	var count int64
	if err := db.Model(&task.OwlTask{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db.Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(generateOrderField(info.OrderKey, info.Desc))
	}else {
		db = db.Order("ct desc")
	}

	var tasks []task.OwlTask
	if err := db.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	for idx, taskV := range tasks {
		formattedItems, _, err := getTaskExecItems(GetDB(), &taskV)
		if err != nil {
			return nil, 0, err
		}

		tasks[idx].ExecItems = formattedItems
	}

	return tasks, count, nil
}

func getTaskExecItems(db *gorm.DB, taskP *task.OwlTask) ([]task.OwlExecItem, []task.OwlSubtask, error) {
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

func (TaskDaoImpl) GetTask(id int64) (*task.OwlTask, error) {
	var task task.OwlTask
	if err := GetDB().First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}

	formattedItems, subTasks, err := getTaskExecItems(GetDB(), &task)
	if err != nil {
		return nil, err
	}

	task.SubTasks = subTasks
	task.ExecItems = formattedItems
	return &task, nil
}

const ExecWaitTasksLimit = 5

func (TaskDaoImpl) GetExecWaitTask() ([]task.OwlTask, int64, error) {
	condition := "status in (?)"

	var count int64
	if err := GetDB().Model(&task.OwlTask{}).Where(condition,
		task.ExecWait).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var tasks []task.OwlTask
	if err := GetDB().Order("et asc").Limit(ExecWaitTasksLimit).
		Find(&tasks, condition, task.ExecWait).Error; err != nil {
		return nil, 0, err
	}

	for idx, taskV := range tasks {
		formattedItems, subTasks, err := getTaskExecItems(GetDB(), &taskV)
		if err != nil {
			return nil, 0, err
		}

		tasks[idx].ExecItems = formattedItems
		tasks[idx].SubTasks = subTasks
	}

	return tasks, count, nil
}
