package admin

import (
	"fmt"
	"time"

	"github.com/nooncall/owls/go/service/tidb_or_mysql"

	"gorm.io/gorm"
)

type OwlAdmin struct {
	ID          int64  `json:"id" gorm:"column:id"`
	Username    string `json:"username" gorm:"username"`
	Description string `json:"description" gorm:"column:description"`

	Ct      int64  `json:"ct" gorm:"column:ct"`
	Creator string `json:"creator" gorm:"creator"`
}

type AdminDao interface {
	AddAdmin(admin *OwlAdmin) (int64, error)
	GetAdmin(username string) (*OwlAdmin, error)
	ListAdmin(pagination *tidb_or_mysql.Pagination) ([]OwlAdmin, int64, error)
	DelAdmin(id int64) error
}

var adminDao AdminDao

func SetAdminDao(impl AdminDao) {
	adminDao = impl
}

func AddAdmin(admin *OwlAdmin) (int64, error) {
	// add admin
	admin.Ct = time.Now().Unix()
	return adminDao.AddAdmin(admin)
}

func ListAdmin(pagination *tidb_or_mysql.Pagination) ([]OwlAdmin, int64, error) {
	return nil, 0, nil
	// todo, del or refactor
	// return adminDao.ListAdmin(pagination)
}

func DelAdmin(id int64) error {
	return adminDao.DelAdmin(id)
}

func IsAdmin(username string) (bool, error) {
	if username == "admin" {
		return true, nil
	}

	_, err := adminDao.GetAdmin(username)
	if gorm.ErrRecordNotFound == err {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("get admin %s err", username)
	}

	return true, nil
}
