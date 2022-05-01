package auth

import "github.com/qingfeng777/owls/server/model/common/request"

type Auth struct {
	ID       int64  `json:"id" gorm:"column:id"`
	UserId   int64  `json:"user_id" gorm:"user_id"`
	DataType string `json:"data_type"`
	Cluster  string `json:"cluster"` // 单个集群名
	DB       string `json:"db"`      // 多个db名，逗号分隔 eg： test,dev
	Status   string `json:"status"`
}

func (a Auth) AddTask() (int64, error) {
	return authDao.AddAuth(&a)
}

const authStatusOn = "ON"

// give auth by set status
func (a Auth) ExecTask() error {
	a.Status = authStatusOn
	return authDao.UpdateAuth(&a)
}

func (a Auth) UpdateTask() error {
	return authDao.UpdateAuth(&a)
}

func (a Auth) ListTask(pageInfo request.SortPageInfo, isDBA bool, status []string) ([]interface{}, int64, error) {
	panic("implement me")
}

func (a Auth) GetTask(id int64) (interface{}, error) {
	return authDao.GetAuth(id)
}

func (a Auth) GetExecWaitTask() ([]interface{}, int64, error) {
	panic("implement me")
}
