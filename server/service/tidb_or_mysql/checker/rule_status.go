package checker

import (
	"fmt"
	"strings"

	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

const ruleClose = 1

type OwlRuleStatus struct {
	Name  string `json:"name" gorm:"column:name"`
	Close int    `json:"close" gorm:"column:close"`
}

type RuleStatusDao interface {
	ListAllStatus() ([]OwlRuleStatus, error)
	UpdateRuleStatus(ruleStatus *OwlRuleStatus) error
}

var ruleStatusDao RuleStatusDao

func SetRuleStatusDao(impl RuleStatusDao) {
	ruleStatusDao = impl
}

func InitRuleStatus() {
	ruleStatuses, err := ruleStatusDao.ListAllStatus()
	if err != nil {
		logger.Errorf("init rules, set rule status err: %s", err.Error())
		return
	}

	for _, ruleStatus := range ruleStatuses {
		for idx, _ := range Rules {
			if Rules[idx].Name == ruleStatus.Name && ruleStatus.Close == ruleClose {
				Rules[idx].Open = false
			}
		}
	}
}

func UpdateRuleStatus(name, action string) (err error) {
	switch action {
	case "close":
		if err = ruleStatusDao.UpdateRuleStatus(&OwlRuleStatus{Name: name, Close: 1}); err == nil {
			updateCacheRuleStatus(name, true)
		}
		return
	case "open":
		if err = ruleStatusDao.UpdateRuleStatus(&OwlRuleStatus{Name: name, Close: 0}); err == nil {
			updateCacheRuleStatus(name, false)
		}
		return
	default:
		return fmt.Errorf("update rule status action not defind, action: %s", action)
	}
}

func ListRules(info request.SortPageInfo) ([]Rule,int) {
	// todo, support sort
	limit := info.PageSize
	offset := info.PageSize * (info.Page - 1)

	result := Rules
	if info.Key != ""{
		result = filterRule(info.Key)
	}

	if offset > len(result){
		return  nil, len(result)
	}
	if offset + limit > len(result){
		limit = len(result) - offset
	}

	return result[offset:offset+limit], len(result)
}

func filterRule(key string) []Rule {
	var rules []Rule
	for _, v := range Rules{
		if strings.Contains(v.Name, key) || strings.Contains(v.Summary, key){
			rules = append(rules, v)
		}
	}
	return rules
}

func updateCacheRuleStatus(name string, close bool) {
	for i, v := range Rules {
		if v.Name == name {
			Rules[i].Open = !close
		}
	}
}
