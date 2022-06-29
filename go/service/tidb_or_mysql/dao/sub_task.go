package dao

import "github.com/nooncall/owls/go/service/tidb_or_mysql/task"

type SubTaskDaoImpl struct {
}

var SubTask SubTaskDaoImpl

func (SubTaskDaoImpl) UpdateItem(item *task.OwlExecItem) error {
	return GetDB().Model(item).Where("id = ?", item.ID).Updates(item).Error
}

func (SubTaskDaoImpl) DelItem(item *task.OwlExecItem) error {
	return GetDB().Model(item).Where("id = ?", item.ID).Delete(item).Error
}

func (SubTaskDaoImpl) UpdateItemByBackupId(item *task.OwlExecItem) error {
	return GetDB().Model(item).Where("backup_id = ?", item.BackupID).Updates(item).Error
}
