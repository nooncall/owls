package task

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	tidb "github.com/pingcap/parser/ast"

	"github.com/nooncall/owls/go/service/tidb_or_mysql/sql_util"
	"github.com/nooncall/owls/go/utils/logger"
)

type BackupDataResp struct {
	DataItems [][]string `json:"data_items"`
	Index     []int      `json:"index"`
	Columns   []string   `json:"columns"`
}

// ListRollbackData ...
func ListRollbackData(req *SqlParam) (*BackupDataResp, error) {
	if req.Sql == "" || req.ClusterName == "" || req.DBName == "" || req.BackupId < 1 {
		logger.Infof("check param failed, originSql : %s ,clusterName :%s ,DBName: %s, backupId: %d",
			req.Sql, req.ClusterName, req.DBName, req.BackupId)
		return nil, fmt.Errorf("get param failed, expert: originSql : %s, clusterName :%s ,DBName: %s, backupId: %d",
			req.Sql, req.ClusterName, req.DBName, req.BackupId)
	}

	backup, err := backupDao.GetBackupInfoById(req.BackupId)
	if err != nil {
		logger.Errorf("get backup info err: %s", err.Error())
		return nil, err
	}

	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		logger.Errorf("get db_info conn err: %s", err.Error())
		return &BackupDataResp{DataItems: dataSplit(backup.Data)}, nil
	}
	defer dbInfo.CloseConn()

	index, cols, err := getUpdateColsInfo(req.Sql, dbInfo.DB)
	if err != nil {
		logger.Warnf("get update index err:%+v", err.Error())
		// index 用于标记被更改的列，可以容忍
	}

	return &BackupDataResp{DataItems: dataSplit(backup.Data), Index: index, Columns: cols}, nil
}

func transIntToInt32(num []int) (resp []int32) {
	for _, v := range num {
		resp = append(resp, int32(v))
	}
	return
}

//返回被更改的列序号，以及列名
func getUpdateColsInfo(originSql string, db *sql.DB) (index []int, cols []string, err error) {
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

	updateCols, err := sql_util.GetSqlColumn(originSql)
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

type SqlParam struct {
	DBName      string `json:"db_name"`
	ClusterName string `json:"cluster_name"`
	TableName   string `json:"table_name"`
	Sql         string `json:"sql"`
	BackupId    int64  `json:"backup_id"`
	Executor    string `json:"executor"`
}

// 原sql，判断是删还是改
// 查出来备份数据
// 如果是删，直接插入回原来的表
// 如果是改,判断哪些字段，根据主键，把原来的字段set回去。
// 最后更改备份状态
func Rollback(req *SqlParam) error {
	if req.Sql == "" || req.ClusterName == "" || req.DBName == "" || req.BackupId < 1 {
		logger.Infof("check param failed, originSql : %s ,clusterName :%s ,DBName: %s, backupId: %d, executor: %s",
			req.Sql, req.ClusterName, req.DBName, req.BackupId, req.Executor)
		return fmt.Errorf("get param failed, expert: originSql : %s, clusterName :%s ,DBName: %s, backupId: %d, creator: %s",
			req.Sql, req.ClusterName, req.DBName, req.BackupId, req.Executor)
	}

	dbInfo, err := dbTool.GetDBConn(req.DBName, req.ClusterName)
	if err != nil {
		return err
	}
	defer dbInfo.CloseConn()

	backup, err := backupDao.GetBackupInfoById(req.BackupId)
	if err != nil {
		logger.Errorf("get backup info err: %s", err.Error())
		return err
	}

	stmtNodes, _, _ := sql_util.Parser.Parse(req.Sql, "", "")
	for _, tiStmt := range stmtNodes {
		switch tiStmt.(type) {
		case *tidb.UpdateStmt:
			err = rollbackUpdate(req.Sql, backup.Data, dbInfo.DB)
		case *tidb.DeleteStmt:
			err = rollbackDel(req.Sql, backup.Data, dbInfo.DB)
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

//更改taskItem 表的备份状态，更改backup表的记录信息
func updateBackupStatus(status ItemStatus, backupId int64, creator string) {
	err := subTaskDao.UpdateItemByBackupId(&OwlExecItem{
		BackupID:     backupId,
		BackupStatus: status,
	})

	if err != nil {
		logger.Errorf("update task item backup status failed, %s", err.Error())
	}

	if err := backupDao.UpdateBackup(&OwlBackup{
		ID:           backupId,
		RollbackTime: time.Now().Unix(),
		RollbackUser: creator,
		IsRolledBack: 1,
	}); err != nil {
		logger.Errorf("update backup record info err: %s", err.Error())
	}
}

//查出来原来的表的列信息，切分备份数据，拼sql， 执行加回去
func rollbackDel(originSql, data string, db *sql.DB) error {
	tableName, err := sql_util.GetTableName(originSql)
	if err != nil {
		return err
	}
	column, err := sql_util.GetTableColumn(tableName, db)
	if err != nil {
		return err
	}

	dataItems := dataSplit(data)
	rollbackSql, err := buildDelRollBackSql(*column, dataItems, tableName)
	if err != nil {
		return err
	}

	if result, err := db.Exec(rollbackSql); err != nil {
		return fmt.Errorf("exec update sql roll back failed : %s, %v", err.Error(), result)
	}

	return nil
}

//查出来原来表的列信息，找到改了哪些列，切分备份数据，拼sql，执行 改回去
func rollbackUpdate(originSql, data string, db *sql.DB) error {
	tableName, err := sql_util.GetTableName(originSql)
	if err != nil {
		return err
	}
	column, err := sql_util.GetTableColumn(tableName, db)
	if err != nil {
		return err
	}

	dataItems := dataSplit(data)
	rollbackSql, err := buildUpdateRollbackSql(*column, dataItems, tableName, originSql)
	if err != nil {
		return err
	}

	for _, v := range rollbackSql {
		if result, err := db.Exec(v); err != nil {
			return fmt.Errorf("exec update sql roll back failed : %s, %v", err.Error(), result)
		}
	}
	return nil
}

func buildDelRollBackSql(column []sql_util.Column, dataItems [][]string, tableName string) (string, error) {
	columnLen := len(column)

	var columnName, values []string
	for _, v := range column {
		columnName = append(columnName, v.Field)
	}

	for _, v := range dataItems {
		if len(v) != columnLen {
			return "", errors.New(fmt.Sprintf("rollback build values failed, length not match, \n column : %s , \n ,data: %s", column, v))
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(addQuotation(v), ", ")))
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		tableName, strings.Join(columnName, ", "), strings.Join(values, ", "))

	logger.Infof("build rollback delete sql : %s ", sql)
	return sql, nil
}

//找到改了哪些列，并给他们按照表里查出来的顺序信息排序；找到主键；切分备份数据
// 除了原sql的where，再把数据自己的主键给加上作为条件
func buildUpdateRollbackSql(column []sql_util.Column, dataItems [][]string, tableName, oldSql string) ([]string, error) {
	updateCols, err := sql_util.GetSqlColumn(oldSql)
	if err != nil {
		return nil, err
	}

	updatedColsOriginSeq := sql_util.SortColOnOriginSeq(column, updateCols)
	priCol, priIndex := sql_util.GetTablePrimaryInfo(column)

	sqlHead := fmt.Sprintf("UPDATE %s SET ", tableName)
	sqlWhereOrigin := "where " + sql_util.GetSqlAfterWhere(oldSql)
	var updateRollbackSql []string
	var set, primaryCondition string
	for _, row := range dataItems {
		// 每条数据一条sql
		set = buildSet(row, updatedColsOriginSeq)
		primaryCondition = buildPrimaryCondition(row, priCol, priIndex)
		updateRollbackSql = append(updateRollbackSql, sqlHead+set+handleFollowingSemicolon(sqlWhereOrigin)+primaryCondition)
	}

	logger.Infof("build rollback update sql : %v ", updateRollbackSql)
	return updateRollbackSql, nil
}

func handleFollowingSemicolon(subSql string) string {
	if strings.Contains(subSql, sql_util.Semicolon) {
		return sql_util.Space + strings.ReplaceAll(subSql, sql_util.Semicolon, sql_util.Space)
	}
	return sql_util.Space + subSql + sql_util.Space
}

func buildSet(row, updatedColsOriginSeq []string) string {
	set := ""
	for i, v := range updatedColsOriginSeq {
		if v != "" {
			set += fmt.Sprintf("%s = '%s',", v, row[i])
		}
	}
	// 去掉最后一个逗号, 后面是开区间
	set = set[:len(set)-1]

	return set
}

func buildPrimaryCondition(row, priCol []string, priIndex []int) string {
	primaryCondition := ""

	for fieldIndex, fieldVal := range row {
		for i := 0; i < len(priCol); i++ {
			//找到主键列
			if fieldIndex == priIndex[i] {
				// 不去重,因为主键可能用的是范围条件,而有重复条件没什么影响
				primaryCondition += fmt.Sprintf("AND %s = '%s'", priCol[i], fieldVal)
			}
		}
	}
	return primaryCondition
}

func dataSplit(data string) (resp [][]string) {
	dataRows := strings.Split(data, splitRowNumberSign)
	for _, v := range dataRows {
		dataItems := strings.Split(v, splitFieldAsterisk)
		resp = append(resp, reverseConvertFields(dataItems))
	}
	return resp
}

func addQuotation(str []string) (resp []string) {
	for _, v := range str {
		resp = append(resp, fmt.Sprintf("'%s'", v))
	}
	return
}
