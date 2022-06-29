package login

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"

	ldap "github.com/jtblin/go-ldap-client"
	"github.com/qingfeng777/owls/server/global"
	"github.com/qingfeng777/owls/server/model/system"
	uuid "github.com/satori/go.uuid"
)

type LdapUser struct {
}

func LdapUserImpl() *LdapUser {
	return &LdapUser{}
}

func (userService *LdapUser) Register(u system.SysUser) (err error, userInter system.SysUser) {
	panic("implement me")
	return err, u
}

func (userService *LdapUser) Login(u *system.SysUser) (err error, userInter *system.SysUser) {
	if !global.Initialized() {
		return fmt.Errorf("db not init"), nil
	}

	// ldap 登录成功后，同步到本地
	ok, _, err := getLdapConn().Authenticate(u.Username, u.Password)
	if err != nil {
		return fmt.Errorf("authenticating user err %s: %+v ", u.Username, err), nil
	}
	if !ok {
		return fmt.Errorf("authenticating failed for user %s", u.Username), nil
	}

	var user system.SysUser
	err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			return err, nil
		}

		u.UUID = uuid.NewV4()
		u.AuthorityId = "886" // todo, refactor to const
		if err = global.GVA_DB.Create(&u).Error; err != nil {
			return err, nil
		}

		err = global.GVA_DB.Where("username = ?", u.Username).Preload("Authorities").Preload("Authority").First(&user).Error
		if err != nil {
			return err, nil
		}
	}

	var am system.SysMenu
	ferr := global.GVA_DB.First(&am, "name = ? AND authority_id = ?", user.Authority.DefaultRouter, user.AuthorityId).Error
	if errors.Is(ferr, gorm.ErrRecordNotFound) {
		user.Authority.DefaultRouter = "404"
	}

	return err, &user
}

var ldapConn *ldap.LDAPClient

func getLdapConn() *ldap.LDAPClient {
	once.Do(setLdapConn)
	return ldapConn
}

var once sync.Once

func setLdapConn() {
	login := global.GVA_CONFIG.Login
	ldapConn = &ldap.LDAPClient{
		Base:         login.Ldap.BaseDn,
		Host:         login.Ldap.Host,
		Port:         login.Ldap.Port,
		UseSSL:       login.Ldap.UseSll,
		BindDN:       login.Ldap.BaseDn,
		BindPassword: login.Ldap.BindPwd,
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
}
