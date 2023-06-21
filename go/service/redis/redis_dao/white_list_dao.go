package redis_dao

import (
	"context"

	"github.com/nooncall/owls/go/service/redis"
	"gorm.io/gorm"
)

type WhiteListDao struct{}

func NewWhiteListDao() redis.WhiteListDao {
	return &WhiteListDao{}
}

func (dao *WhiteListDao) AddWhiteList(ctx context.Context, db *gorm.DB, WhiteList *redis.WhiteList) (int64, error) {
	result := db.Create(WhiteList)
	if result.Error != nil {
		return 0, result.Error
	}
	return WhiteList.ID, nil
}

func (dao *WhiteListDao) DelWhiteList(ctx context.Context, db *gorm.DB, id int64) error {
	result := db.Delete(&redis.WhiteList{}, id)
	return result.Error
}

func (dao *WhiteListDao) ListWhiteList(ctx context.Context, db *gorm.DB, cmdType string) ([]redis.WhiteList, error) {
	var whiteLists []redis.WhiteList
	result := db.Where("cmd_type = ?", cmdType).Find(&whiteLists)
	if result.Error != nil {
		return nil, result.Error
	}
	return whiteLists, nil
}
