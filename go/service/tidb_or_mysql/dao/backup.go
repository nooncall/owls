package dao

import (
	"github.com/nooncall/owls/go/service/tidb_or_mysql/task"
)

type BackupImpl struct {
}

var BackupDAO BackupImpl

func (BackupImpl) AddBackup(backup *task.OwlBackup) (int64, error) {
	err := GetDB().Create(backup).Error
	return backup.ID, err
}

func (BackupImpl) UpdateBackup(backup *task.OwlBackup) error {
	return GetDB().Model(backup).Where("id = ?", backup.ID).Updates(backup).Error
}

func (BackupImpl) GetBackupInfoById(id int64) (*task.OwlBackup, error) {
	var backup task.OwlBackup
	return &backup, GetDB().Where("id = ?", id).First(&backup).Error
}
