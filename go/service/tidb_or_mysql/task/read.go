package task

import (
	"errors"
	"fmt"

	"github.com/nooncall/owls/go/service/tidb_or_mysql/sql_util"
	"github.com/nooncall/owls/go/utils/logger"
)

type ReadResult struct {
	DataItems interface{} `json:"data_items"`
	Columns   []string    `json:"columns"`
}

// ListRollbackData ...
func ReadData(req *SqlParam) (*ReadResult, error) {
	if req.Sql == "" || req.ClusterName == "" || req.DBName == "" {
		logger.Infof("check param failed, originSql : %s ,clusterName :%s ,DBName: %s, backupId: %d",
			req.Sql, req.ClusterName, req.DBName, req.BackupId)
		return nil, fmt.Errorf("get param failed, expert: originSql : %s, clusterName :%s ,DBName: %s, backupId: %d",
			req.Sql, req.ClusterName, req.DBName, req.BackupId)
	}

	req.Sql = sql_util.AddLimit(req.Sql)

	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		logger.Errorf("get db_info conn err: %s", err.Error())
		return nil, err
	}
	defer dbInfo.CloseConn()

	rows, err := dbInfo.DB.Query(req.Sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result ReadResult
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	result.Columns = cols

	vals := make([]interface{}, len(cols))
	valsP := make([]interface{}, len(vals))
	//将接口转换为指针类型的接口
	for i := range vals {
		valsP[i] = &vals[i]
	}

	var data [][]interface{}
	for rows.Next() {
		if err := rows.Scan(valsP...); err != nil {
			return nil, err
		}

		for i := range vals {
			if v, ok := vals[i].([]byte); ok { //读取的数据是uint8类型的数组，需要转成byte类型的数组才好转换成其他
				vals[i] = string(v)
			}
		}
		data = append(data, vals)
	}
	result.DataItems = data

	return &result, nil
}

func GetTableInfo(req *SqlParam) (string, error) {
	if req.TableName == "" || req.ClusterName == "" || req.DBName == "" {
		return "", fmt.Errorf("check param failed, tableName : %s ,clusterName :%s ,DBName: %s",
			req.TableName, req.ClusterName, req.DBName)
	}

	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		return "", err
	}
	defer dbInfo.CloseConn()

	rows, err := dbInfo.DB.Query(fmt.Sprintf("show create table %s;", req.TableName))
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var table, tableInfo string
	for rows.Next() {
		if err = rows.Scan(&table, &tableInfo); err != nil {
			return "", err
		}
		return tableInfo, nil
	}
	return "", errors.New("get table info found nothing")
}
