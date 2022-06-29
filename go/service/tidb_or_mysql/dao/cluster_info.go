package dao

import (
	"gorm.io/gorm"

	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/db_info"
	"github.com/nooncall/owls/go/utils"
)

type ClusterImpl struct {
}

var Cluster ClusterImpl

func (ClusterImpl) AddCluster(db *gorm.DB, cluster *db_info.OwlCluster) (int64, error) {
	err := db.Create(cluster).Error
	return cluster.ID, err
}

func (ClusterImpl) UpdateCluster(db *gorm.DB, cluster *db_info.OwlCluster) error {
	return db.Model(cluster).Where("id = ?", cluster.ID).Updates(cluster).Error
}

func (ClusterImpl) DelCluster(db *gorm.DB, id int64) error {
	return db.Where("id = ?", id).Delete(&db_info.OwlCluster{}).Error
}

func (ClusterImpl) GetClusterByName(db *gorm.DB, name string) (*db_info.OwlCluster, error) {
	var cluster db_info.OwlCluster
	return &cluster, db.First(&cluster, "name = ?", name).Error
}

func (ClusterImpl) ListCluster(db *gorm.DB, info request.SortPageInfo) ([]db_info.OwlCluster, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db = db.Offset(offset).Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(utils.GenerateOrderField(info.OrderKey, info.Desc))
	}
	if info.Key != "" {
		fmtKey := "%" + info.Key + "%"
		db = db.Where("name like ? or description like ? or addr like ?",
			fmtKey, fmtKey, fmtKey)
	}

	var clusters []db_info.OwlCluster
	return clusters, db.Find(&clusters).Error
}

func (ClusterImpl) ListAllCluster(db *gorm.DB) ([]db_info.OwlCluster, error) {
	var clusters []db_info.OwlCluster
	return clusters, db.Find(&clusters).Error
}
