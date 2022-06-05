package system

import (
	"github.com/pkg/errors"
	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/system"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

var User = new(user)

type user struct{}

func (u *user) TableName() string {
	return "sys_users"
}

func (u *user) Initialize() error {
	//todo, refactor to config;
	entities := []system.SysUser{
		{UUID: uuid.NewV4(), Username: "dev", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "Dev", HeaderImg: "https://qmplusimg.henrongyi.top/gva_header.jpg", AuthorityId: "888", Phone: "17611111111", Email: "333333333@qq.com"},
		{UUID: uuid.NewV4(), Username: "admin", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "Admin", HeaderImg: "https://qmplusimg.henrongyi.top/gva_header.jpg", AuthorityId: "887", Phone: "17611111112", Email: "633333333@qq.com"},
		{UUID: uuid.NewV4(), Username: "user", Password: "e10adc3949ba59abbe56e057f20f883e", NickName: "User", HeaderImg: "https:///qmplusimg.henrongyi.top/1572075907logo.png", AuthorityId: "886", Phone: "17611111113", Email: "933333333@qq.com"},
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
