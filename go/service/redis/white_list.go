package redis

import (
	"context"
	"errors"

	"github.com/nooncall/owls/go/utils"
	"gorm.io/gorm"
)

type WhiteList struct {
	ID         int64  `json:"id" gorm:"column:id"`
	Cmd        string `json:"cmd" gorm:"column:cmd"`
	CheckType  string `json:"check_type" gorm:"column:check_type"` //exist;limit_num,xx;only_string
	CmdType    string `json:"cmd_type"  gorm:"column:cmd_type"`
	Creator    string `json:"creator" gorm:"column:creator"`
	CreateTime int64  `json:"create_time" gorm:"column:create_time"`
}

type WhiteListDao interface {
	AddWhiteList(ctx context.Context, db *gorm.DB, WhiteList *WhiteList) (int64, error)
	DelWhiteList(ctx context.Context, db *gorm.DB, id int64) error
	ListWhiteList(ctx context.Context, db *gorm.DB, cmdType string) ([]WhiteList, error)
}

var whiteListDao WhiteListDao

func SetWhiteListDao(impl WhiteListDao) {
	whiteListDao = impl
}

func ListReadWhiteList(ctx context.Context, db *gorm.DB) ([]WhiteList, error) {
	return whiteListDao.ListWhiteList(ctx, GetDB(), "read")
}

func ListWriteWhiteList(ctx context.Context, db *gorm.DB) ([]WhiteList, error) {
	return whiteListDao.ListWhiteList(ctx, GetDB(), "write")
}

func DelWhiteList(ctx context.Context, db *gorm.DB, id int64) error {
	return whiteListDao.DelWhiteList(ctx, GetDB(), id)
}

func AddWhiteList(ctx context.Context, db *gorm.DB, WhiteList *WhiteList) (int64, error) {
	if WhiteList == nil || (WhiteList.CmdType != "read" && WhiteList.CmdType != "write") {
		return 0, errors.New("params error, nil or type error. ")
	}

	return whiteListDao.AddWhiteList(ctx, GetDB(), WhiteList)
}

func getReadWhiteCmd() map[string]string {
	return map[string]string{
		"get":       CheckTypeExist,
		"mget":      CheckTypeExist,
		"hget":      CheckTypeExist,
		"hmget":     CheckTypeExist,
		"lrange":    "limit_num,100",
		"zrange":    "limit_num,100",
		"sismember": CheckTypeExist,
		"scard":     CheckTypeExist,
		"zcard":     CheckTypeExist,
		"hscan":     CheckTypeExist,
		"ttl":       CheckTypeExist,
		"type":      CheckTypeExist,
		"hlen":      CheckTypeExist,
		"exists":    CheckTypeExist,
		"sscan":     CheckTypeExist,
	}
}

func getWriteWhiteCmd() map[string]string {
	return map[string]string{
		"set":    CheckTypeExist,
		"mset":   CheckTypeExist,
		"hset":   CheckTypeExist,
		"hmset":  CheckTypeExist,
		"hdel":   CheckTypeExist,
		"del":    CheckTypeExist,
		"zrem":   CheckTypeExist,
		"srem":   CheckTypeExist,
		"sadd":   CheckTypeExist,
		"zadd":   CheckTypeExist,
		"incrby": CheckTypeExist,
	}
}

func getAllWhiteCmd() map[string]string {
	return utils.MergeStringMaps(getReadWhiteCmd(), getWriteWhiteCmd())
}

func getReadCmdType(key string) string {
	return getReadWhiteCmd()[key]
}

func getWriteCmdType(key string) string {
	return getWriteWhiteCmd()[key]
}
