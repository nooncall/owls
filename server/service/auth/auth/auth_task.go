package auth

import "github.com/qingfeng777/owls/server/model/common/request"

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
