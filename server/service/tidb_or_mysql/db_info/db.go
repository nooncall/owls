package db_info

import (
	"database/sql"
	"fmt"

	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/task"
)

const defaultDBName = "information_schema"

type DBInfoTool struct {
}

var DBTool DBInfoTool

func (DBInfoTool) GetDBConn(dbName, clusterName string) (*task.DBInfo, error) {
	cluster, err := GetClusterByName(clusterName)
	if err != nil {
		return nil, fmt.Errorf("get cluster info err: %s", err.Error())
	}

	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, dbName))
	if err != nil {
		return nil, fmt.Errorf("open db conn err: %s", err.Error())
	}

	defaultDB, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, defaultDBName))
	if err != nil {
		return nil, fmt.Errorf("open db conn err: %s", err.Error())
	}

	return &task.DBInfo{DB: db, DefaultDB: defaultDB, DBName: dbName}, nil
}

// return dbs and mapping cluster
func ListAllDB(clusterName string) ([]string, error) {
	cluster, err := GetClusterByName(clusterName)
	if err != nil {
		return nil, err
	}

	return ListDbByCluster(cluster)
}

func ListDbByCluster(cluster *OwlCluster) ([]string, error) {
	/*deCryptoData, err := utils.AesDeCrypto(utils.ParseStringedByte(cluster.Pwd))
	if err != nil {
		return nil, err
	}

	cluster.Pwd = fmt.Sprintf("%s", deCryptoData)*/
	conn, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8", cluster.User, cluster.Pwd, cluster.Addr, defaultDBName))
	if err != nil {
		return nil, fmt.Errorf("open db_info conn err: %s", err.Error())
	}
	defer conn.Close()

	rows, err := conn.Query("show databases;")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dbs []string
	for rows.Next() {
		var dbName string
		if err = rows.Scan(&dbName); err != nil {
			return nil, err
		}
		dbs = append(dbs, dbName)
	}
	return dbs, nil
}
