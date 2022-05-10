package task

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

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

	backup, err := backupDao.GetBackupInfoById(req.BackupId)
	if err != nil {
		logger.Errorf("get backup info err: %s", err.Error())
		return nil, err
	}

	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		logger.Errorf("get db_info conn err: %s", err.Error())
		return &ReadResult{DataItems: dataSplit(backup.Data)}, nil
	}
	defer dbInfo.CloseConn()

	index, cols, err := getUpdateColsInfo(req.OriginSql, dbInfo.DB)
	if err != nil {
		logger.Warnf("get update index err:%+v", err.Error())
		// index 用于标记被更改的列，可以容忍
	}

	return &ReadResult{DataItems: dataSplit(backup.Data), Index: index, Columns: cols}, nil
}

//todo, 用来获取数据的列名，sql中应该可以拿到。
func getUpdateColsInfoBackForRead(originSql string, db *sql.DB) (index []int, cols []string, err error) {
	var tableName string
	var noIndex bool
	tableName, _ = sql_util.GetTableName(originSql)

	stmtNodes, _, _ := sql_util.Parser.Parse(originSql, "", "")
	for _, tiStmt := range stmtNodes {
		switch tiStmt.(type) {
		case *tidb.DeleteStmt:
			noIndex = true
		}
	}

	column, err := sql_util.GetTableColumn(tableName, db)
	if err != nil {
		return
	}
	for _, v := range *column {
		cols = append(cols, v.Field)
	}
	if noIndex {
		return
	}

	updateCols, err := sql_util.GetUpdateColumn(originSql)
	if err != nil {
		return
	}
	colsOnTabSeq := sql_util.SortColOnOriginSeq(*column, updateCols)
	for i, v := range colsOnTabSeq {
		if v != "" {
			index = append(index, i)
		}
	}
	return
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

	stmtNodes, _, _ := sql_util.Parser.Parse(req.OriginSql, "", "")
	for _, tiStmt := range stmtNodes {
		switch tiStmt.(type) {
		case *tidb.UpdateStmt:
			err = rollbackUpdate(req.OriginSql, backup.Data, dbInfo.DB)
		case *tidb.DeleteStmt:
			err = rollbackDel(req.OriginSql, backup.Data, dbInfo.DB)
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

