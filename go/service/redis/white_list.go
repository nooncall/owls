package redis

import "github.com/nooncall/owls/go/utils"

func getReadWhiteCmd() map[string]string {
	return map[string]string{
		"get":       CheckTypeExist,
		"mget":      CheckTypeExist,
		"hget":      CheckTypeExist,
		"hmget":     CheckTypeExist,
		"lrange":    "limit_num,100",
		"zrange":    "limit_num,100",
		"sismember": CheckTypeExist,
		"scard":     CheckTypeExist,
		"zcard":     CheckTypeExist,
		"hscan":     CheckTypeExist,
		"ttl":       CheckTypeExist,
		"type":      CheckTypeExist,
		"hlen":      CheckTypeExist,
		"exists":    CheckTypeExist,
		"sscan":     CheckTypeExist,
	}
}

func getWriteWhiteCmd() map[string]string {
	return map[string]string{
		"set":    CheckTypeExist,
		"mset":   CheckTypeExist,
		"hset":   CheckTypeExist,
		"hmset":  CheckTypeExist,
		"hdel":   CheckTypeExist,
		"del":    CheckTypeExist,
		"zrem":   CheckTypeExist,
		"srem":   CheckTypeExist,
		"sadd":   CheckTypeExist,
		"zadd":   CheckTypeExist,
		"incrby": CheckTypeExist,
	}
}

func getAllWhiteCmd() map[string]string {
	return utils.MergeStringMaps(getReadWhiteCmd(), getWriteWhiteCmd())
}

func getCmdType(key string) string {
	return getAllWhiteCmd()[key]
}
