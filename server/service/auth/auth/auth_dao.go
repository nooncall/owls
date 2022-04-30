package auth

import (
	"gorm.io/gorm"

	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/utils"
)

type authTaskDaoImpl struct {
}

var authTaskDao authTaskDaoImpl

func GetDB() *gorm.DB {
	// todo, refactor to config
	return global.GVA_DB.Debug()
}

func (authTaskDaoImpl) AddAuthTask(AuthTask *AuthTask) (int64, error) {
	tx := GetDB().Begin()
	if err := tx.Create(AuthTask).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return AuthTask.ID, tx.Commit().Error
}

func (authTaskDaoImpl) UpdateAuthTask(AuthTask *AuthTask) error {
	return GetDB().Model(AuthTask).Where("id = ?", AuthTask.ID).Updates(AuthTask).Error
}

func (authTaskDaoImpl) ListAuthTask(info request.SortPageInfo, isDBA bool, status []string) ([]AuthTask, int64, error) {
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
	if err := db.Model(&AuthTask{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db.Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(utils.GenerateOrderField(info.OrderKey, info.Desc))
	} else {
		db = db.Order("ct desc")
	}

	var AuthTasks []AuthTask
	if err := db.Find(&AuthTasks).Error; err != nil {
		return nil, 0, err
	}

	return AuthTasks, count, nil
}

func (authTaskDaoImpl) GetAuthTask(id int64) (*AuthTask, error) {
	var AuthTask AuthTask
	return &AuthTask, GetDB().First(&AuthTask, "id = ?", id).Error
}
