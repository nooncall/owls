package tidb_or_mysql

import (
	"time"

	"github.com/flipped-aurora/gin-vue-admin/server/config"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

type Pagination struct {
	Offset int    `json:"offset"`
	Limit  int    `json:"limit"`
	Key    string `json:"key"`

	Operator string `json:"operator"`
}

type ClockInf interface {
	Now() time.Time
	NowUnix() int64
}

var Clock ClockInf

func SetClock(impl ClockInf) {
	Clock = impl
}

type RealClock struct{}

func (RealClock) Now() time.Time                         { return time.Now() }
func (RealClock) NowUnix() int64                         { return time.Now().Unix() }
func (RealClock) After(d time.Duration) <-chan time.Time { return time.After(d) }

// todo, maybe better place ?
func InitConfigLog() {
	logger.InitLog(".", "test.log", "debug")
	config.InitConfig("../../config/config.yml")
}
