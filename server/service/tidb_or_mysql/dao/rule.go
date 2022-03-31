package dao

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/checker"
	"gorm.io/gorm"
)

type RuleDaoImpl struct {
}

var Rule RuleDaoImpl

func (RuleDaoImpl) ListAllStatus() ([]checker.OwlRuleStatus, error) {
	var ruleStatus []checker.OwlRuleStatus
	return ruleStatus, GetDB().Find(&ruleStatus).Error
}

func (RuleDaoImpl) UpdateRuleStatus(ruleStatus *checker.OwlRuleStatus) error {
	err := GetDB().Where("name = ?", ruleStatus.Name).First(&checker.OwlRuleStatus{}).Error
	if err == gorm.ErrRecordNotFound{
		return GetDB().Create(ruleStatus).Error
	}
	if err != nil {
		return err
	}

	return GetDB().Model(ruleStatus).Where("name = ?", ruleStatus.Name).Updates(ruleStatus).Error
}
