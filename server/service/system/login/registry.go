package login

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/system"
	"github.com/qingfeng777/owls/server/utils"
	uuid "github.com/satori/go.uuid"
)

type RegistryUser struct {
}

func RegistryUserImpl() *RegistryUser {
	return &RegistryUser{}
}

func (userService *RegistryUser) Register(u system.SysUser) (err error, userInter system.SysUser) {
	var user system.SysUser
	if !errors.Is(global.GVA_DB.Where("username = ?", u.Username).First(&user).Error, gorm.ErrRecordNotFound) { // 判断用户名是否注册
		return errors.New("用户名已注册"), userInter
	}
	// 否则 附加uuid 密码md5简单加密 注册
	u.Password = utils.MD5V([]byte(u.Password))
	u.UUID = uuid.NewV4()
	u.AuthorityId = "886"
	err = global.GVA_DB.Create(&u).Error
	return err, u
}

func (userService *RegistryUser) Login(u *system.SysUser) (err error, userInter *system.SysUser) {
	if !global.Initialized() {
		return fmt.Errorf("db not init"), nil
	}

	var user system.SysUser
	u.Password = utils.MD5V([]byte(u.Password))
	err = global.GVA_DB.Where("username = ? AND password = ?", u.Username, u.Password).Preload("Authorities").Preload("Authority").First(&user).Error
	if err == nil {
		var am system.SysMenu
		ferr := global.GVA_DB.First(&am, "name = ? AND authority_id = ?", user.Authority.DefaultRouter, user.AuthorityId).Error
		if errors.Is(ferr, gorm.ErrRecordNotFound) {
			user.Authority.DefaultRouter = "404"
		}
	}
	return err, &user
}
