package login_check

import (
	"fmt"
	"sync"

	ldap "github.com/jtblin/go-ldap-client"
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/utils/logger"
)

type LoginServiceImpl struct {
}

var LoginService LoginServiceImpl

func (LoginServiceImpl) Login(userName, pwd string) error {
	ok, _, err := getLdapConn().Authenticate(userName, pwd)
	if err != nil {
		return fmt.Errorf("authenticating user err %s: %+v ", userName, err)
	}
	if !ok {
		return fmt.Errorf("authenticating failed for user %s", "username")
	}
	logger.Infof("user: %s Login", userName)
	return nil
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
