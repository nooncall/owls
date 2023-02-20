package redis

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/nooncall/owls/go/utils"
	"github.com/nooncall/owls/go/utils/logger"
)

const (
	CheckTypeExist      = "exist"
	CheckTypeNum        = "limit_num"
	CheckTypeOnlyString = "only_string"
)

const CheckPass = "pass"

// config eg: key: cmd, val : CheckType,LenLimit
type PermitCmd struct {
	Cmd       string `json:"cmd"`        // key; 去匹配用户输入命令的第一个字段
	CheckType string `json:"check_type"` // val-1; 0 just match, 1 check nums , 2 only string
	LenLimit  int64  `json:"len_limit"`  // val-2;
}

type Check struct {
}

var checker Check

func GetAllWhitelist(ctx context.Context) ([]PermitCmd, error) {
	var rules []PermitCmd
	for k, v := range getAllWhiteCmd() {
		checkType, limit := getCheckTypeAndLenLimit(ctx, v)
		rules = append(rules, PermitCmd{
			Cmd:       k,
			CheckType: checkType,
			LenLimit:  limit,
		})
	}

	return rules, nil
}

func (Check) CheckReadCmd(ctx context.Context, cmd, cluster string, db int) (bool, string, error) {
	cmd = utils.DelUselessSpace(cmd)
	checkV := getCmdType(strings.ToLower(getCmdPrefix(cmd)))
	if len(checkV) < 1 {
		return false, "此命令不允许使用", nil
	}

	checkType, limit := getCheckTypeAndLenLimit(ctx, checkV)
	switch checkType {
	case CheckTypeExist:
		return true, CheckPass, nil
	case CheckTypeNum:
		return checkSpecifyCmdAffectedNums(ctx, cmd, cluster, limit, db)
	default:
		return false, fmt.Sprintf("exec cmd check err, type not found, cmd: %s, checkType: %s", cmd, checkType), nil
	}
}

func (Check) CheckWriteCmd(ctx context.Context, cmd, prefix, service string, usePrefix bool) (bool, string, error) {
	cmd = utils.DelUselessSpace(cmd)
	checkV := getCmdType(strings.ToLower(getCmdPrefix(cmd)))
	if len(checkV) < 1 {
		return false, "此命令不允许使用", nil
	}

	checkType, _ := getCheckTypeAndLenLimit(ctx, checkV)
	switch checkType {
	case CheckTypeExist:
		return true, CheckPass, nil
	case CheckTypeOnlyString:
		return checkOnlyString(ctx, cmd, prefix, service, usePrefix)
	default:
		logger.Infof("exec cmd check err, type not found, cmd: %s, checkType: %s", cmd, checkType)
		return false, "", fmt.Errorf("exec cmd check err, type not found, cmd: %s, checkType: %s", cmd, checkType)
	}
}

func checkSpecifyCmdAffectedNums(ctx context.Context, cmd, cluster string, limit int64, db int) (bool, string, error) {
	cmdSplit := strings.Split(cmd, " ")
	if len(cmdSplit) < 2 {
		return false, "", fmt.Errorf("check cmd faild, err cmd. cmd: %s", cmd)
	}

	var affectNum int64
	var err error
	redisCli, err := NewRedisCli(cluster, db)
	if err != nil {
		return false, "", err
	}

	switch strings.ToLower(cmdSplit[0]) {
	case "hgetall":
		affectNum, err = redisCli.HLen(ctx, cmdSplit[1]).Result()
		if err != nil {
			return false, "", fmt.Errorf("while chek, get hlen err. key: %s, err: %s", cmdSplit[0], err.Error())
		}
	case "smembers":
		affectNum, err = redisCli.SCard(ctx, cmdSplit[1]).Result()
		if err != nil {
			return false, "", fmt.Errorf("while chek, get scard err. key: %s, err: %s", cmdSplit[0], err.Error())
		}
	case "zrange", "lrange":
		if len(cmdSplit) < 4 {
			return false, "", fmt.Errorf("check cmd zrange faild, err cmd. cmd: %s", cmd)
		}
		if cmdSplit[3] == "-1" {
			return false, "此命令不允许使用-1作为边界", nil
		}
		affectNum, err = utils.SubTwoStrNum(cmdSplit[3], cmdSplit[2])
		if err != nil {
			return false, "", fmt.Errorf("parse cmd range faild, err cmd. cmd: %s", cmd)
		}
	}

	if affectNum <= limit {
		return true, CheckPass, nil
	}
	logger.Infof("affect num big than limit, limit : %d, affectNum: %d", limit, affectNum)
	return false, fmt.Sprintf("命令影响数据量超限，限制值:%d", limit), nil
}

// cmd eg: del hi hello
func checkOnlyString(ctx context.Context, cmd, cluster string, db int) (bool, string, error) {
	cmdSplit := strings.Split(cmd, " ")
	if len(cmdSplit) < 2 {
		return false, "", fmt.Errorf("check cmd faild, err cmd. cmd: %s", cmd)
	}

	redis, err := NewRedisCli(cluster, db)
	if err != nil {
		return false, "", err
	}
	for _, key := range cmdSplit[1:] {
		keyType, err := redis.Type(ctx, key).Result()
		if err != nil {
			return false, "", fmt.Errorf("while chek, get type err. key: %s, err: %s", key, err.Error())
		}

		if keyType != "string" {
			logger.Infof("only string check failed, key: %s, keyType:%s", key, keyType)
			return false, fmt.Sprintf("此命令仅可对string使用"), nil
		}
	}
	return true, CheckPass, nil
}

func getCheckTypeAndLenLimit(ctx context.Context, checkV string) (string, int64) {
	strs := strings.Split(checkV, ",")
	if len(strs) < 2 {
		return strs[0], 0
	}

	lenLimit, err := strconv.ParseInt(strs[1], 10, 64)
	if err != nil {
		logger.Errorf("parse len limit err, config check vale: %s", checkV)
		return strs[0], 0
	}

	return strs[0], lenLimit
}

func getCmdPrefix(cmd string) string {
	return strings.Split(cmd, " ")[0]
}
