package task

import (
	"errors"
	"fmt"

	tidb "github.com/pingcap/parser/ast"
	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/sql_util"
	"github.com/qingfeng777/owls/server/utils/logger"
)

type ReadResult struct {
	DataItems [][]string `json:"data_items"`
	Columns   []string   `json:"columns"`
}

// ListRollbackData ...
func ReadData(req *SqlParam) (*ReadResult, error) {
	if req.OriginSql == "" || req.ClusterName == "" || req.DBName == "" || req.BackupId < 1 {
		logger.Infof("check param failed, originSql : %s ,clusterName :%s ,DBName: %s, backupId: %d",
			req.OriginSql, req.ClusterName, req.DBName, req.BackupId)
		return nil, fmt.Errorf("get param failed, expert: originSql : %s, clusterName :%s ,DBName: %s, backupId: %d",
			req.OriginSql, req.ClusterName, req.DBName, req.BackupId)
	}


	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		logger.Errorf("get db_info conn err: %s", err.Error())
		return nil, err
	}
	defer dbInfo.CloseConn()

	var result ReadResult
	result.Columns, err = sql_util.GetSqlColumn(req.OriginSql)
	if err != nil {
		return nil, err
	}

	//使用db信息，执行读sql ，获取结果

	return &result, nil
}

func GetTableInfo(req *SqlParam) error {
	if req.TableName == "" || req.ClusterName == "" || req.DBName == "" {
		return fmt.Errorf("check param failed, tableName : %s ,clusterName :%s ,DBName: %s",
			req.TableName, req.ClusterName, req.DBName)
	}

	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		return err
	}
	defer dbInfo.CloseConn()
	// 怎么获取table info呢？ 适合边写边测；
	// 执行desc table 名，二维数组获取、数组+interface 获取结果

	stmtNodes, _, _ := sql_util.Parser.Parse(req.OriginSql, "", "")
	for _, tiStmt := range stmtNodes {
		switch tiStmt.(type) {
		case *tidb.UpdateStmt:
		case *tidb.DeleteStmt:
		default:
			logger.Warnf("rollback, operate type not found")
			return errors.New("sql operate type not support, only update, delete supported")
		}
	}

	if err != nil {
		logger.Errorf("rollback sql err: %s", err.Error())
		updateBackupStatus(ItemRollBackFailed, req.BackupId, req.Executor)
		return fmt.Errorf("rollback err: %s", err.Error())
	}
	updateBackupStatus(ItemRollBackSuccess, req.BackupId, req.Executor)

	return nil
}
