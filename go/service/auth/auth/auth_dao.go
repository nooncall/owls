package auth

import (
	"gorm.io/gorm"

	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/common/request"
	"github.com/qingfeng777/owls/server/utils"
)

type authDaoImpl struct {
}

var authDao authDaoImpl

func GetDB() *gorm.DB {
	// todo, refactor to config
	return global.GVA_DB.Debug()
}

func (authDaoImpl) AddAuth(AuthTask *Auth) (int64, error) {
	tx := GetDB().Begin()
	if err := tx.Create(AuthTask).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return AuthTask.ID, tx.Commit().Error
}

func (authDaoImpl) DelAuth(id int64) error {
	return GetDB().Delete(&Auth{}, "id = ?", id).Error
}

func (authDaoImpl) UpdateAuth(AuthTask *Auth) error {
	return GetDB().Model(AuthTask).Where("id = ?", AuthTask.ID).Updates(AuthTask).Error
}

func (authDaoImpl) ListAuth(info request.SortPageInfo, status []string) ([]Auth, int64, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := GetDB().Offset(offset)
	if info.Key != "" {
		fmtKey := "%" + info.Key + "%"
		db = db.Where("id like ? or username like ? or cluster like ? or db like ?",
			fmtKey, fmtKey, fmtKey, fmtKey)
	}
	db = db.Where("status in (?)", status)

	var count int64
	if err := db.Model(&Auth{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	db.Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(utils.GenerateOrderField(info.OrderKey, info.Desc))
	} else {
		db = db.Order("id desc")
	}

	var AuthTasks []Auth
	if err := db.Find(&AuthTasks).Error; err != nil {
		return nil, 0, err
	}

	return AuthTasks, count, nil
}

func (authDaoImpl) GetAuth(id int64) (*Auth, error) {
	var AuthTask Auth
	return &AuthTask, GetDB().First(&AuthTask, "id = ?", id).Error
}

func (authDaoImpl) ListAuthForFilter(userId uint, status, dataType string) ([]Auth, error) {
	db := GetDB().Where("user_id = ? and status = ? and data_type = ?", userId, status, dataType)

	var authTasks []Auth
	if err := db.Find(&authTasks).Error; err != nil {
		return nil, err
	}

	return authTasks, nil
}
