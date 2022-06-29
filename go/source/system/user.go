package system

import (
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/model/system"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var User = new(user)

type user struct{}

func (u *user) TableName() string {
	return "sys_users"
}

func (u *user) Initialize() error {
	entities := []system.SysUser{
		{UUID: uuid.NewV4(), Username: "dev", Password: global.GVA_CONFIG.Login.InitPassword, NickName: "Dev", HeaderImg: "https://pica.zhimg.com/v2-4b6c9c82868b5d7c26a39776ffc5bc98_r.jpg", AuthorityId: "888", Phone: "17611111111", Email: "333333333@qq.com"},
		{UUID: uuid.NewV4(), Username: "admin", Password: global.GVA_CONFIG.Login.InitPassword, NickName: "Admin", HeaderImg: "https://pica.zhimg.com/v2-16022df2e21a5f69838d094856f29024_1440w.jpg", AuthorityId: "887", Phone: "17611111112", Email: "633333333@qq.com"},
		{UUID: uuid.NewV4(), Username: "user", Password: global.GVA_CONFIG.Login.InitPassword, NickName: "User", HeaderImg: "https://pica.zhimg.com/v2-2bb8251e0eb4bf14612bf02939e6ced0_1440w.jpg", AuthorityId: "886", Phone: "17611111113", Email: "933333333@qq.com"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, u.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (u *user) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("username = ?", "a303176530").First(&system.SysUser{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
