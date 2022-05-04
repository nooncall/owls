package tidb_or_mysql

import (
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/pkg/errors"

	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/db_info"
	"github.com/qingfeng777/owls/server/utils"
)

var Cluster = new(cluster)

type cluster struct{}

func (u *cluster) TableName() string {
	return "owl_clusters"
}

func (u *cluster) Initialize() error {
	conf := global.GVA_CONFIG.Mysql
	cryptoData, err := utils.AesCrypto([]byte(conf.Password))
	if err != nil {
		return fmt.Errorf("crypto password err: %v", err)
	}

	Pwd := utils.StringifyByteDirectly(cryptoData)
	entities := []db_info.OwlCluster{
		{Name: "self cluster", Description: "cluster this project using", Addr: conf.Path + conf.Port, User: conf.Username, Pwd: Pwd, Ct: time.Now().Unix(), Operator: "init"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, u.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (u *cluster) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.First(&db_info.OwlCluster{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
