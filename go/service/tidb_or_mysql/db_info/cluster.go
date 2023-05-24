package db_info

import (
	"gorm.io/gorm"
	"time"

	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/model/common/request"
	"github.com/nooncall/owls/go/service/auth/auth"
	"github.com/nooncall/owls/go/service/tidb_or_mysql"
	"github.com/nooncall/owls/go/utils"
)

// todo, list support type, backend and ui
type OwlCluster struct {
	ID          int64  `json:"id" gorm:"column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Description string `json:"description" gorm:"column:description"`
	Addr        string `json:"addr" gorm:"column:addr"` //ip : port
	User        string `json:"user" gorm:"column:user"`
	Pwd         string `json:"pwd" gorm:"column:pwd"`
	CType       string `json:"ctype" gorm:"column:ctype"`

	Ct       int64  `json:"ct" gorm:"column:ct"`
	Ut       int64  `json:"ut" gorm:"column:ut"`
	Operator string `json:"operator" gorm:"column:operator"`
}

type ClusterDao interface {
	AddCluster(db *gorm.DB, cluster *OwlCluster) (int64, error)
	UpdateCluster(db *gorm.DB, cluster *OwlCluster) error
	DelCluster(db *gorm.DB, id int64) error
	GetClusterByName(db *gorm.DB, clusterName string) (*OwlCluster, error)
	ListCluster(db *gorm.DB, pageInfo request.SortPageInfo) ([]OwlCluster, error)
	ListAllCluster(db *gorm.DB) ([]OwlCluster, error)
}

var clusterDao ClusterDao

func SetClusterDao(impl ClusterDao) {
	clusterDao = impl
}

func AddCluster(cluster *OwlCluster) (int64, error) {
	cryptoData, err := utils.AesCrypto([]byte(cluster.Pwd))
	if err != nil {
		return 0, err
	}
	cluster.Ct = time.Now().Unix()

	cluster.Pwd = utils.StringifyByteDirectly(cryptoData)
	return clusterDao.AddCluster(global.GetDB(), cluster)
}

func UpdateCluster(cluster *OwlCluster) error {
	if cluster.Pwd == pwdReplace {
		cluster.Pwd = ""
	}

	if cluster.Pwd != "" {
		cryptoData, err := utils.AesCrypto([]byte(cluster.Pwd))
		if err != nil {
			return err
		}

		cluster.Pwd = utils.StringifyByteDirectly(cryptoData)
	}
	cluster.Ut = tidb_or_mysql.Clock.NowUnix()

	return clusterDao.UpdateCluster(global.GetDB(), cluster)
}

func DelCluster(id int64) error {
	return clusterDao.DelCluster(global.GetDB(), id)
}

func GetClusterByName(name string) (*OwlCluster, error) {
	cluster, err := clusterDao.GetClusterByName(global.GetDB(), name)
	if err != nil {
		return nil, err
	}

	deCryptoData, err := utils.AesDeCrypto(utils.ParseStringedByte(cluster.Pwd))
	if err != nil {
		return nil, err
	}
	cluster.Pwd = string(deCryptoData)
	return cluster, nil
}

func ListCluster(pageInfo request.SortPageInfo) ([]OwlCluster, error) {
	return clusterDao.ListCluster(global.GetDB(), pageInfo)
}

const pwdReplace = "******"

func ListClusterForUI(pageInfo request.SortPageInfo) ([]OwlCluster, error) {
	clusters, err := ListCluster(pageInfo)
	if err != nil {
		return nil, err
	}

	for i, _ := range clusters {
		clusters[i].Pwd = pwdReplace
	}
	return clusters, nil
}

func ListClusterName(userId uint, filter bool) ([]string, error) {
	clusters, err := clusterDao.ListAllCluster(global.GetDB())
	if err != nil {
		return nil, err
	}

	var result []string
	for _, v := range clusters {
		result = append(result, v.Name)
	}

	if !filter || !global.GVA_CONFIG.DBFilter.ReadNeedAuth {
		return result, nil
	}

	return auth.FilterCluster(result, userId, auth.DB), nil
}
