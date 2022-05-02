package auth

import "github.com/qingfeng777/owls/server/model/common/request"

type AuthAble interface {
	AuthCluster() ([]string, error)
	AuthDB(cluster string) ([]string, error)
	GetManagers() []string
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
	if err := GetDB().Find(&authTasks, "username = ? and data_type = ?",a.username, a.dataType).Error; err != nil {
		return nil, err
	}

	var result []string
	for _, auth := range authTasks{
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
	for _, auth := range authTasks{
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
	return authDao.ListAuth(info, []string{authStatusOn})
}

func DelAuth(id int64) error {
	return authDao.DelAuth(id)
}

func AddAuth(authTask *Auth) (int64, error){
	return authDao.AddAuth(authTask)
}