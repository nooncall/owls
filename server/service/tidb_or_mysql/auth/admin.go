package auth

import (
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/admin"
)

type AdminAuthToolImpl struct {
}

var AdminAuthService AdminAuthToolImpl

func (AdminAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	admins, _, err := admin.ListAdmin(&tidb_or_mysql.Pagination{Limit: 10})
	if err != nil {
		return "", err
	}

	var resp string
	for i, v := range admins {
		if i == 0 {
			resp += v.Username
		} else {
			resp += "," + v.Username
		}
	}
	return resp, nil
}

func (AdminAuthToolImpl) IsDba(userName string) (isDba bool, err error) {
	return true, nil
	// return admin.IsAdmin(userName)
}
