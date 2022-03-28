package task

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/pingcap/parser/ast"

	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/sql_util"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

type BackupDao interface {
	AddBackup(backup *OwlBackup) (int64, error)
	UpdateBackup(backup *OwlBackup) error
	GetBackupInfoById(id int64) (*OwlBackup, error)
}

var backupDao BackupDao

func SetBackupDao(impl BackupDao) {
	backupDao = impl
}

type OwlBackup struct {
	ID           int64  `json:"id" gorm:"column:id"`
	Data         string `json:"data" gorm:"column:data"`
	Ct           int64  `json:"ct" gorm:"column:ct"`
	RollbackTime int64  `json:"rollback_time" gorm:"column:rollback_time"`
	RollbackUser string `json:"rollback_user" gorm:"column:rollback_user"`
	IsRolledBack int    `json:"is_rolled_back" gorm:"column:is_rolled_back"`
}

const (
	splitFieldAsterisk = "*"
	splitRowNumberSign = "#"
	NUL                = ""
)
const (
	replaceAsterisk   = "0d65760d0deb56cb59a49ebbe0939cf3"
	replaceNumberSign = "c13c1fce3ca97e765860002118ca3bb8"
	replaceNUL        = "4252a939eaac0e4f457f296e60882ebe"
)

func backup(db *sql.DB, taskType, sql string) (execBackup bool, backupId int64, err error) {
	if taskType != DML ||
		!needBackup(sql) {
		logger.Infof("not a update or delete operate , don't backup")
		return
	}
	execBackup = true
	logger.Infof("start backup data ...")

	selectSql, tableName, err := getSqlInfo(sql)
	if err != nil {
		logger.Errorf("convert sql to select err: %s", err.Error())
		return
	}
	isEmpty, backupId, err := fetchAndStoreBackupInfo(db, selectSql, tableName)
	if isEmpty {
		err = errors.New("backup nothing, condition match nothing")
	}
	return
}

func needBackup(sql string) bool {
	stmtNodes, _, _ := sql_util.Parser.Parse(sql, "", "")

	for _, tiStmt := range stmtNodes {
		switch tiStmt.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			return true
		default:
			return false
		}
	}
	return false
}

func fetchAndStoreBackupInfo(db *sql.DB, selectSql, tableName string) (isEmpty bool, backupId int64, err error) {
	rows, err := db.Query(selectSql)
	if err != nil {
		logger.Warnf("exec backup db_info exec err:%v", err)
		return
	}

	column, err := sql_util.GetTableColumn(tableName, db)
	if err != nil {
		return
	}

	dataStr := formatData(rows, column)
	if strings.TrimSpace(dataStr) == "" {
		isEmpty = true
		logger.Warnf("while backup, condition match nothing")
		return
	}

	backupId, err = backupDao.AddBackup(&OwlBackup{
		Data: dataStr,
		Ct:   time.Now().Unix(),
	})
	return
}

func getSqlInfo(sql string) (selectSql string, tableName string, err error) {
	tableName, err = sql_util.GetTableName(sql)

	where := sql_util.GetSqlAfterWhere(sql)
	selectSql = fmt.Sprintf("select * from %s where %s", tableName, sql_util.HandelKeyWorldForCondition(where))
	logger.Infof("build backup select sql : %s ", selectSql)

	return
}

//map 乱序, 需要按列顺序存
// 用#分割行,用*分割字段
// 需要把字段中的 # * 空字符替换掉, 展示以及回滚的时候再替换回来
func formatData(row *sql.Rows, columns *[]sql_util.Column) string {
	defer row.Close()
	values, err := sql_util.ScanMap(row)
	if err != nil {
		logger.Errorf("format data scanMap rows failed : %s", err.Error())
		return ""
	}

	var resp string
	for _, rowMap := range values {
		rowStr := ""
		for _, column := range *columns {
			if fieldValue, ok := rowMap[column.Field]; ok {
				rowStr += splitFieldAsterisk + convertField(uint8ToString(fieldValue))
			} else {
				logger.Errorf("backup data format data error, column : %s not found . data : %v", column.Field, rowMap)
			}
		}
		if len(rowStr) >= 1 {
			rowStr = rowStr[1:]
		}
		resp += splitRowNumberSign + rowStr
	}
	if len(resp) >= 1 {
		resp = resp[1:]
	}
	return resp
}

func convertField(fieldStr string) string {
	if fieldStr == NUL {
		return replaceNUL
	}
	fieldStr = strings.ReplaceAll(fieldStr, splitRowNumberSign, replaceNumberSign)
	return strings.ReplaceAll(fieldStr, splitFieldAsterisk, replaceAsterisk)
}

func reverseConvertFields(str []string) (resp []string) {
	for _, v := range str {
		resp = append(resp, reverseConvertField(v))
	}
	return
}

func reverseConvertField(str string) string {
	if str == replaceNUL {
		return NUL
	}
	str = strings.ReplaceAll(str, replaceAsterisk, splitFieldAsterisk)
	return strings.ReplaceAll(str, replaceNumberSign, splitRowNumberSign)
}

func uint8ToString(inter interface{}) string {
	uint8Array, ok := inter.([]uint8)
	if !ok {
		if inter != nil {
			logger.Errorf("uint8 to string error : received interface isn't a uint8 slice, data: %v", inter)
		}
		return ""
	}

	var byteArray []byte
	for _, v := range uint8Array {
		byteArray = append(byteArray, byte(v))
	}
	return string(byteArray)
}

func isNum(str string) bool {
	for _, v := range str {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
