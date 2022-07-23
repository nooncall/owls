package tidb_or_mysql

import (
	"fmt"
	"github.com/nooncall/owls/go/model/system"
	"gorm.io/gorm"
	"time"

	"github.com/pkg/errors"

	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/service/tidb_or_mysql/db_info"
	"github.com/nooncall/owls/go/utils"
)

var Cluster = new(cluster)

type cluster struct{}

func (u *cluster) TableName() string {
	return "owl_clusters"
}

func (u *cluster) Initialize(initData *system.InitDBData) error {
	cryptoData, err := utils.AesCrypto([]byte(initData.Password))
	if err != nil {
		return fmt.Errorf("crypto password err: %v", err)
	}

	Pwd := utils.StringifyByteDirectly(cryptoData)
	entities := []db_info.OwlCluster{
		{Name: "self-cluster", Description: "cluster this project using", Addr: initData.Host, User: initData.UserName, Pwd: Pwd, Ct: time.Now().Unix(), Operator: "init"},
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
