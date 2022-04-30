package auth

import "github.com/qingfeng777/owls/server/model/common/request"

type AuthTask struct {
	ID       int64  `json:"id" gorm:"column:id"`
	UserId   int64  `json:"user_id" gorm:"user_id"`
	DataType string `json:"data_type"`
	Cluster  string `json:"cluster"` // 单个集群名
	DB       string `json:"db"`      // 多个db名，逗号分隔 eg： test,dev
	Status   string `json:"status"`
}

func (a AuthTask) AddTask() (int64, error) {
	return authTaskDao.AddAuthTask(&a)
}

const authStatusOn = "ON"

// give auth by set status
func (a AuthTask) ExecTask() error {
	a.Status = authStatusOn
	return authTaskDao.UpdateAuthTask(&a)
}

func (a AuthTask) UpdateTask() error {
	return authTaskDao.UpdateAuthTask(&a)
}

func (a AuthTask) ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) ([]interface{}, int64, error) {
	panic("implement me")
}

func (a AuthTask) GetTask(id int64) (interface{}, error) {
	return authTaskDao.GetAuthTask(id)
}

func (a AuthTask) GetExecWaitTask() ([]interface{}, int64, error) {
	panic("implement me")
}
