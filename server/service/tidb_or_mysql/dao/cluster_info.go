package dao

import "github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/db_info"

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

func (ClusterImpl) ListCluster() ([]db_info.OwlCluster, error) {
	var clusters []db_info.OwlCluster
	return clusters, GetDB().Find(&clusters).Error
}
