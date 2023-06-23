package auth

import "github.com/nooncall/owls/go/model/common/request"

type AuthAble interface {
	AuthCluster() ([]string, error)
	AuthDB(cluster string) ([]string, error)
	GetManagers() []string
}

type Auth struct {
	ID           int64  `json:"id" gorm:"column:id"`
	UserId       uint   `json:"user_id" gorm:"user_id"`
	ParentTaskID int64  `json:"parent_task_id" gorm:"parent_task_id"`
	Username     string `json:"username" gorm:"username"`
	DataType     string `json:"data_type"`
	Cluster      string `json:"cluster"` // 单个集群名
	DB           string `json:"db"`      // 单个db名
	Status       string `json:"status"`
}

type authTool struct {
	dataType, username string
}

func GetAuthTool(dataType, username string) authTool {
	return authTool{dataType, username}
}

// todo, 都是list，根据条件获取权限结果
func (a authTool) AuthCluster() ([]string, error) {
	var authTasks []Auth
	if err := GetDB().Find(&authTasks, "username = ? and data_type = ?", a.username, a.dataType).Error; err != nil {
		return nil, err
	}

	var result []string
	for _, auth := range authTasks {
		result = append(result, auth.Cluster)
	}
	return result, nil
}

func (a authTool) AuthDB(cluster string) ([]string, error) {
	var authTasks []Auth
	if err := GetDB().Find(&authTasks, "username = ? and data_type = ? and cluster = ?",
		a.username, a.dataType, cluster).Error; err != nil {
		return nil, err
	}

	var result []string
	for _, auth := range authTasks {
		result = append(result, auth.Cluster)
	}
	return result, nil
}

func (authTool) GetManagers() []string {
	panic("implement me")
}

// todo, cluster 管理要从tidb 中迁移出来。

func ListDataType() []string {
	return Types()
}

func ListAuth(info request.SortPageInfo) ([]Auth, int64, error) {
	return authDao.ListAuth(info, []string{StatusPass})
}

func DelAuth(id int64) error {
	return authDao.DelAuth(id)
}

func AddAuth(authTask *Auth) (int64, error) {
	return authDao.AddAuth(authTask)
}
