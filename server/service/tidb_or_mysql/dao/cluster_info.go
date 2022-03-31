package dao

import (
	"fmt"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/db_info"
)

type ClusterImpl struct {
}

var Cluster ClusterImpl

func (ClusterImpl) AddCluster(cluster *db_info.OwlCluster) (int64, error) {
	err := GetDB().Create(cluster).Error
	return cluster.ID, err
}

func (ClusterImpl) UpdateCluster(cluster *db_info.OwlCluster) error {
	return GetDB().Model(cluster).Where("id = ?", cluster.ID).Updates(cluster).Error
}

func (ClusterImpl) DelCluster(id int64) error {
	return GetDB().Where("id = ?", id).Delete(&db_info.OwlCluster{}).Error
}

func (ClusterImpl) GetClusterByName(name string) (*db_info.OwlCluster, error) {
	var cluster db_info.OwlCluster
	return &cluster, GetDB().First(&cluster, "name = ?", name).Error
}

func (ClusterImpl) ListCluster(info request.SortPageInfo) ([]db_info.OwlCluster, error) {
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	db := GetDB().Offset(offset).Limit(limit)
	if info.OrderKey != "" {
		db = db.Order(generateOrderField(info.OrderKey, info.Desc))
	}
	if info.Key != "" {
		fmtKey := "%" + info.Key + "%"
		db = db.Where("name like ? or description like ? or addr like ?",
			fmtKey, fmtKey, fmtKey)
	}

	var clusters []db_info.OwlCluster
	return clusters, db.Find(&clusters).Error
}

func generateOrderField(key string, desc bool) string {
	if desc {
		return fmt.Sprintf("%s desc", key)
	}

	return fmt.Sprintf("%s asc", key)
}
