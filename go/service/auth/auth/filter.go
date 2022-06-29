package auth

import "github.com/nooncall/owls/go/utils/logger"

func FilterDB(dbs []string, userId uint, dataType, cluster string) []string {
	f := "FilterDB()-->: "
	auths, err := authDao.ListAuthForFilter(userId, StatusPass, DB)
	if err != nil {
		logger.Errorf("%s get auth err: %v", f, err)
		return dbs
	}

	var result []string
	for _, db := range dbs {
		for _, auth := range auths {
			if auth.Cluster == cluster && auth.DB == db && auth.DataType == dataType {
				result = append(result, db)
				break
			}
		}
	}

	return result
}

func FilterCluster(clusters []string, userId uint, dataType string) []string {
	f := "FilterCluster()-->: "
	auths, err := authDao.ListAuthForFilter(userId, StatusPass, DB)
	if err != nil {
		logger.Errorf("%s get auth err: %v", f, err)
		return clusters
	}

	var result []string
	for _, cluster := range clusters {
		for _, auth := range auths {
			if auth.Cluster == cluster && auth.DataType == dataType {
				result = append(result, cluster)
				break
			}
		}
	}

	return result
}
