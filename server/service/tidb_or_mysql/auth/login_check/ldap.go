package login_check

import (
	"fmt"
	"sync"

	ldap "github.com/jtblin/go-ldap-client"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
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
	ldapConn = &ldap.LDAPClient{
		Base:         config.Conf.Login.LDAP.BaseDN,
		Host:         config.Conf.Login.LDAP.Host,
		Port:         config.Conf.Login.LDAP.Port,
		UseSSL:       config.Conf.Login.LDAP.UseSSL,
		BindDN:       config.Conf.Login.LDAP.BindDN,
		BindPassword: config.Conf.Login.LDAP.BindPwd,
		UserFilter:   "(uid=%s)",
		GroupFilter:  "(memberUid=%s)",
		Attributes:   []string{"givenName", "sn", "mail", "uid"},
	}
}
