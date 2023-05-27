package redis

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nooncall/owls/go/utils"
)

func ExecReadTask(ctx context.Context, params *Params) (interface{}, error) {
	pass, msg, err := checker.CheckReadCmd(ctx, params.Cmd, params.Cluster, params.DB)
	if err != nil {
		return nil, err
	}
	if !pass {
		return nil, errors.New(msg)
	}

	return exec(ctx, params.Cmd, params.Cluster, params.DB)
}

func load(ctx context.Context, key interface{}) (value interface{}, err error) {
	return "", nil
}

func exec(ctx context.Context, cmd, cluster string, db int) (resp interface{}, err error) {
	//todo, redis as a param
	redisCli, err := NewRedisCli(cluster, db)
	if err != nil {
		return nil, err
	}

	cmd = utils.DelUselessSpace(cmd)
	cmdSplit := strings.Split(cmd, " ")
	if len(cmdSplit) < 2 {
		return nil, fmt.Errorf("while exec cmd err: wrong cmd, cmd: %s", cmd)
	}

	switch strings.ToLower(cmdSplit[0]) {
	case "mset":
		var pairs []interface{}
		for _, v := range cmdSplit[1:] {
			pairs = append(pairs, v)
		}
		cmdResult := redisCli.MSet(ctx, pairs)
		return cmdResult.Val(), cmdResult.Err()
	// multi continuous keys
	case "mget":
		cmdResult := redisCli.Do(ctx, utils.StringArrayToInterfaceArray(cmdSplit)...)
		return cmdResult.Val(), cmdResult.Err()
	// one key
	default:
		var othersParams []interface{}
		if len(cmdSplit) >= 3 {
			for _, v := range cmdSplit[2:] {
				othersParams = append(othersParams, v)
			}
		}
		cmdResult := redisCli.Do(ctx, utils.StringArrayToInterfaceArray(cmdSplit)...)
		return cmdResult.Val(), cmdResult.Err()
	}
}

func filterNilStr(data interface{}) string {
	str := fmt.Sprintf("%v", data)
	return strings.ReplaceAll(str, "<nil>", "")
}
