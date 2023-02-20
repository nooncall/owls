package redis_dao

import (
	"gorm.io/gorm"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/service/redis"
	"github.com/nooncall/owls/go/utils"
)

type TaskDaoImpl struct {
}

var Task TaskDaoImpl

func (TaskDaoImpl) AddTask(db *gorm.DB, task *redis.RedisTask) (int64, error) {
	if err := db.Create(task).Error; err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (TaskDaoImpl) UpdateTask(db *gorm.DB, task *redis.RedisTask) error {
	return db.Model(task).Where("id = ?", task.ID).Updates(task).Error
}

func (TaskDaoImpl) ListRedisTaskByTaskID(db *gorm.DB, id int64) ([]redis.RedisTask, error) {
	var tasks []redis.RedisTask
	return tasks, db.Find(&tasks).Error
}

func (TaskDaoImpl) ListTask(db *gorm.DB, info request.SortPageInfo, isDBA bool, status []string) ([]redis.RedisTask, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db = db.Offset(offset)
	if info.Key != "" {
		fmtKey := "%" + info.Key + "%"
		db = db.Where("id = ? or status = ? or name like ? or creator like ?",
			fmtKey, fmtKey, fmtKey, fmtKey)
	}
	db = db.Where("status in (?) and creator = ?", status, info.Operator)

	var count int64
	if err := db.Model(&redis.RedisTask{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db.Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(utils.GenerateOrderField(info.OrderKey, info.Desc))
	} else {
		db = db.Order("ct desc")
	}

	var tasks []redis.RedisTask
	if err := db.Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, count, nil
}

func (TaskDaoImpl) GetTask(db *gorm.DB, id int64) (*redis.RedisTask, error) {
	var task redis.RedisTask
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &task, nil
}
