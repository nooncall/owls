package auth

import (
	"errors"
	"strings"

	"github.com/qingfeng777/owls/server/global"
)

type ConfAuthToolImpl struct {
}

var ConfAuthService ConfAuthToolImpl

func (ConfAuthToolImpl) GetReviewer(userName string) (reviewerName string, err error) {
	return strings.Join(global.GVA_CONFIG.DBFilter.Reviewers, ","), nil
}

func (ConfAuthToolImpl) IsDba(userName string) (isDba bool, err error) {
	if len(global.GVA_CONFIG.DBFilter.Reviewers) < 1 {
		return false, errors.New("dba members not config")
	}
	for _, v := range global.GVA_CONFIG.DBFilter.Reviewers {
		if v == userName {
			return true, nil
		}
	}

	return false, nil
}
