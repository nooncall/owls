package checker

import (
	"strconv"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"

	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/sql_util"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

/**
对以下两种操作做判断
ALTER TABLE testalter_tbl MODIFY c CHAR(10);
ALTER TABLE testalter_tbl CHANGE i j BIGINT;
*/
func unsupportedTypeChange(sql string, tiStmt ast.StmtNode, info *task.DBInfo) bool {
	old, new, tableName, newFieldType, newFieldLen, notChange := getInfoByParser(sql, tiStmt)
	if notChange {
		logger.Infof("while get Field info by parser err: stmt don't have value ")
		return true
	}

	column, err := sql_util.GetTableColumn(tableName, info.DB)
	if err != nil {
		logger.Errorf("get column info from db_info err: %s", err.Error())
		return true
	}

	//modify 的话对原字段操作，new即old
	if old == "" {
		old = new
	}

	return IsBanned(*column, old, newFieldType, newFieldLen)
}

/**
判断是否为禁止的操作
规则要求: alter change/modify操作， 禁止以下动作：
1）int/bigint --> varchar
2) bigint --> int
3) varchar --> bigint/int
4) varchar(n) --> varchar(m) 其中n>m
*/
/**
sql 类型和mysql.type 对应关系
varchar  varchar
int      long
bigint   longlong
*/
func IsBanned(column []sql_util.Column, old string, newFieldType byte, NewFieldLen int) bool {
	var oldCol *sql_util.Column
	for _, v := range column {
		if strings.ToLower(v.Field) == old {
			oldCol = &v
			break
		}
	}
	if oldCol == nil {
		logger.Warnf("while check is banned change, get oldCol err,old: %s cols : %v", old, column)
		return false
	}

	logger.Infof("ban change check info: old: %s, newFieldType: %v, newFiledLen: %d, oldCol: %v",
		old, newFieldType, NewFieldLen, oldCol)

	switch {
	case strings.Contains(oldCol.Type, "bigint"):
		if newFieldType == mysql.TypeVarchar || newFieldType == mysql.TypeLong {
			return true
		}
	case strings.Contains(oldCol.Type, "int"):
		if newFieldType == mysql.TypeVarchar {
			return true
		}
	case strings.Contains(oldCol.Type, "varchar"):
		if newFieldType == mysql.TypeLonglong || newFieldType == mysql.TypeLong {
			return true
		}
		if newFieldType == mysql.TypeVarchar && NewFieldLen < getLenFromType(oldCol.Type) {
			return true
		}
	}

	return false
}

func getLenFromType(str string) int {
	if !strings.Contains(str, "(") || !strings.Contains(str, ")") {
		logger.Errorf("get field len from type err, type : %s", str)
		return -1
	}
	numStr := str[strings.Index(str, "(")+1 : strings.Index(str, ")")]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		logger.Errorf("get field len from type err: %s", err.Error())
	}
	return num
}

func getInfoByParser(sql string, tiStmt ast.StmtNode) (old, new, tableName string,
	newFieldType byte, newFieldLen int, notChange bool) {
	node := tiStmt.(*ast.AlterTableStmt)
	if len(node.Specs) < 1 || len(node.Specs[0].NewColumns) < 1 {
		logger.Warnf("while get Field info by parser err: stmt.NewColumns don't have value ")
		notChange = true
		return
	}

	if node.Specs[0].Tp != ast.AlterTableModifyColumn &&
		node.Specs[0].Tp != ast.AlterTableChangeColumn {
		notChange = true
		return
	}

	if node.Specs[0].OldColumnName != nil {
		old = node.Specs[0].OldColumnName.Name.L
	}
	tableName = node.Table.Name.L
	new = node.Specs[0].NewColumns[0].Name.Name.L
	newFieldType = node.Specs[0].NewColumns[0].Tp.Tp
	newFieldLen = node.Specs[0].NewColumns[0].Tp.Flen
	logger.Infof("get sql info by parser, tableName: %s,old:%s, new: %s, newFieldType: %v, newFieldLen: %d",
		tableName, old, new, newFieldType, newFieldLen)
	return
}
