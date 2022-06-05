package system

import (
	"github.com/pkg/errors"
	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/system"
	"gorm.io/gorm"
)

var UserAuthority = new(userAuthority)

type userAuthority struct{}

func (a *userAuthority) TableName() string {
	var entity system.SysUseAuthority
	return entity.TableName()
}

func (a *userAuthority) Initialize() error {
	entities := []system.SysUseAuthority{
		{SysUserId: 1, SysAuthorityAuthorityId: "888"},
		{SysUserId: 2, SysAuthorityAuthorityId: "887"},
		{SysUserId: 3, SysAuthorityAuthorityId: "886"},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil {
		return errors.Wrap(err, a.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (a *userAuthority) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("sys_user_id = ? AND sys_authority_authority_id = ?", 2, "888").First(&system.SysUseAuthority{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
