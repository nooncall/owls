package sql_util

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"vitess.io/vitess/go/mysql"
	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/format"
	_ "github.com/pingcap/tidb/types/parser_driver"

	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

//取出数据的顺序同建表语句的顺序
type Column struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default *string
	Extra   string
}

func buildDelRollBackSql(column []Column, dataItems [][]string, tableName string) (string, error) {
	columnLen := len(column)

	var columnName, values []string
	for _, v := range column {
		columnName = append(columnName, v.Field)
	}

	for _, v := range dataItems {
		if len(v) != columnLen {
			return "", fmt.Errorf("rollback build values failed, length not match, \n column : %v , \n ,data: %s", column, v)
		}
		values = append(values, fmt.Sprintf("(%s)", strings.Join(addQuotation(v), ", ")))
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s;",
		tableName, strings.Join(columnName, ", "), strings.Join(values, ", "))

	logger.Infof("build rollback delete sql : %s ", sql)
	return sql, nil
}

func addQuotation(str []string) (resp []string) {
	for _, v := range str {
		resp = append(resp, fmt.Sprintf("'%s'", v))
	}
	return
}

func GetTableColumn(tableName string, db *sql.DB) (*[]Column, error) {
	sql := fmt.Sprintf("show columns from %s ", tableName)
	rows, err := db.Query(sql)
	if err != nil {
		logger.Warnf("get Table column info err:%v, sql content: %s", err, sql)
		return nil, err
	}
	defer rows.Close()

	var column []Column
	for rows.Next() {
		var col Column
		if err := rows.Scan(&col.Field, &col.Type, &col.Null, &col.Key, &col.Default, &col.Extra); err != nil {
			return nil, err
		}
		column = append(column, col)
	}
	logger.Infof("Table column is: %v", column)
	return &column, err
}

//仅dml用
func GetTableName(sql string) (string, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		logger.Errorf("get Table name,parse sql failed, sql : %s", sql)
		return "", err
	}
	var left ast.ResultSetNode
	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.UpdateStmt:
			//fmt.Printf("%#v\n", node)
			left = node.TableRefs.TableRefs.Left
		case *ast.DeleteStmt:
			left = node.TableRefs.TableRefs.Left
		}

		tSource, ok := left.(*ast.TableSource)
		if !ok {
			logger.Errorf("get Table name, assert failed, sql : %s", sql)
			return "", errors.New("get Table name failed")
		}
		tName, ok := tSource.Source.(*ast.TableName)
		if !ok {
			logger.Errorf("get Table name, assert failed, sql : %s", sql)
			return "", errors.New("get Table name failed")
		}
		return tName.Name.L, nil
	}
	return "", nil
}

// add `` to key words
// todo 有瑕疵，同时有关键字和value内包含关键字内容的情况（eg: xxxx  where `index` = 'index hello'; ），会出问题。待更换更好的方式。
func HandelKeyWorldForCondition(originWhere string) string {
	targetWhere := originWhere
	columns := GetCondition(originWhere)
	for _, v := range columns {
		if IsKeyWord(v) {
			targetWhere = strings.ReplaceAll(targetWhere, v, fmt.Sprintf("`%s`", v))
		}
	}

	return targetWhere
}

// 获取sql的条件，字符串处理
// 方式： 小写化，and前后有空格，空格分割，再依次进一步判断。
func GetCondition(sqlAfterWhere string) []string {
	sqlAfterWhere = strings.ToLower(sqlAfterWhere)
	strs := strings.Split(sqlAfterWhere, "and")
	var condition []string
	for _, v := range strs {
		if condi := getOneCondition(v); condi != "" {
			condition = append(condition, condi)
		}
	}
	logger.Infof("sql condition is : %v", condition)
	return condition
}

//检查=， >，<，>=，<=，like，in
// != ,not in 已经在上一步过滤掉
func getOneCondition(str string) string {
	switch {
	case strings.Contains(str, ">="):
		return strings.TrimSpace(strings.Split(str, ">=")[0])
	case strings.Contains(str, "<="):
		return strings.TrimSpace(strings.Split(str, "<=")[0])
	case strings.Contains(str, "="):
		return strings.TrimSpace(strings.Split(str, "=")[0])
	case strings.Contains(str, ">"):
		return strings.TrimSpace(strings.Split(str, ">")[0])
	case strings.Contains(str, "<"):
		return strings.TrimSpace(strings.Split(str, "<")[0])
	case strings.Contains(str, " like "):
		return strings.TrimSpace(strings.Split(str, "like")[0])
	case strings.Contains(str, " not in "):
		return strings.TrimSpace(strings.Split(str, "not in")[0])
	case strings.Contains(str, " in "):
		return strings.TrimSpace(strings.Split(str, "in")[0])
	case strings.Contains(str, " between "):
		return strings.TrimSpace(strings.Split(str, "between")[0])
	}

	return ""
}

//注意这个返回的列是乱序的
func GetUpdateColumn(sql string) ([]string, error) {
	stmtNodes, _, err := getParser().Parse(sql, "", "")
	if err != nil {
		return nil, err
	}

	var cols []string
	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.UpdateStmt:
			for _, v := range node.List {
				cols = append(cols, v.Column.String())
			}
		default:
			errStr := fmt.Sprintf("get sql update cols err, not update operate, sql: %s", sql)
			logger.Errorf(errStr)
			return nil, errors.New(errStr)
		}
	}
	return cols, nil
}

type ColsSeq struct {
	index int
	field string
}

// columnInfos：所有的顺序的列信息， column： 待排序列信息
// 注意这个返回值里有空值，为了对应其原来的顺序和位置
// eg : {"id","","age"}
func SortColOnOriginSeq(columnInfos []Column, column []string) []string {
	length := len(columnInfos)
	resp := make([]string, length, length)
	for _, colName := range column {
		for index, columnInfo := range columnInfos {
			if columnInfo.Field == colName {
				resp[index] = colName
			}
		}
	}
	return resp
}

func getNotNullSortedCols(cols []string) (resp []string) {
	for _, v := range cols {
		if v != "" {
			resp = append(resp, v)
		}
	}
	return resp
}

func GetTablePrimaryInfo(columnInfos []Column) (primaryCol []string, index []int) {
	for i, v := range columnInfos {
		if v.Key == "PRI" {
			primaryCol = append(primaryCol, v.Field)
			index = append(index, i)
		}
	}

	return
}

type WriterBuffer struct {
	Condition strings.Builder
}

func (v *WriterBuffer) Write(p []byte) (n int, err error) {
	v.Condition.Write(p)
	return
}

var Parser *parser.Parser = parser.New()

//返回值为where后面的部分，不包括where
func GetSqlAfterWhere(sql string) string {
	stmtNodes, _, err := Parser.Parse(sql, "", "")
	if err != nil {
		logger.Errorf("get sql condition err: %s", err.Error())
	}

	writer := &WriterBuffer{}
	ctx := format.NewRestoreCtx(format.RestoreStringSingleQuotes|format.RestoreKeyWordLowercase, writer)
	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.UpdateStmt:
			node.Where.Restore(ctx)
		case *ast.DeleteStmt:
			node.Where.Restore(ctx)
		default:
			logger.Errorf("get sql after err, not supported type. sql: %s", sql)
		}
	}
	return writer.Condition.String()
}

func SplitMultiSql(sql string) (resp []string, err error) {
	stmtNodes, _, err := getParser().Parse(sql, "utf8mb4", "")
	if err != nil {
		logger.Infof("parser sql error : %s", err.Error())
		return nil, fmt.Errorf("there is an error in your SQL syntax: %s", err.Error())
	}

	for _, v := range stmtNodes {
		resp = append(resp, deleteSpecifyCharAtHead(v.Text()))
	}

	// vitess can't parse multi sql once
	for i, v := range resp {
		_, err = sqlparser.Parse(v)
		if err != nil {
			logger.Infof("vitess parser sql error : %s", err.Error())
			return nil, fmt.Errorf("there is an error in your SQL syntax: %s, at num: %d", err.Error(), i)
		}
	}
	return
}

var sqlParser *parser.Parser

func getParser() *parser.Parser {
	if sqlParser != nil {
		return sqlParser
	}
	sqlParser = parser.New()
	return sqlParser
}

const (
	Semicolon   = `;`
	Space       = ` `
	DoubleSpace = `  `
	Table       = `	`
	NewLine     = `
`
)

func deleteSpecifyCharAtHead(str string) string {
	if len(str) < 1 {
		return str
	}
	head := str[:1]
	if head == Semicolon || head == Space || head == NewLine {
		return deleteSpecifyCharAtHead(str[1:])
	}
	return str
}

func replaceSpecifyChar(str string) string {
	if len(str) < 1 {
		return str
	}
	if !strings.Contains(str, DoubleSpace) &&
		!strings.Contains(str, Table) &&
		!strings.Contains(str, NewLine) {
		return str
	}
	str = strings.ReplaceAll(str, NewLine, Space)
	str = strings.ReplaceAll(str, Table, Space)
	return replaceSpecifyChar(strings.ReplaceAll(str, DoubleSpace, Space))
}

const KeyJoinChar = "+"

func IsSubKey(keyA, keyB string) bool {
	var short, long string
	if len(keyA) < len(keyB) {
		short, long = keyA, keyB
	} else {
		short, long = keyB, keyA
	}

	subKeys := strings.Split(long, KeyJoinChar)
	for _, v := range subKeys {
		if v == short {
			return true
		}
	}
	return false
}

// tidb 暂时不支持'SELECT * FROM mysql.help_keyword;', 支持后可以替换一下
var keyWords = map[string]struct{}{
	"ADD": {}, "ADMIN": {}, "ALL": {}, "ALTER": {}, "ANALYZE": {}, "AND": {}, "AS": {}, "ASC": {},

	"BETWEEN": {}, "BIGINT": {}, "BINARY": {}, "BLOB": {}, "BOTH": {}, "BUCKETS": {}, "BUILTINS": {}, "BY": {},

	"CANCEL": {}, "CASCADE": {}, "CASE": {}, "CHANGE": {}, "CHAR": {}, "CHARACTER": {}, "CHECK": {}, "CMSKETCH": {}, "COLLATE": {},
	"COLUMN": {}, "CONSTRAINT": {}, "CONVERT": {}, "CREATE": {}, "CROSS": {}, "CURRENT_DATE": {}, "CURRENT_ROLE": {}, "CURRENT_TIME": {}, "CURRENT_TIMESTAMP": {},
	"CURRENT_USER": {},

	"DATABASE": {}, "DATABASES": {}, "DAY_HOUR": {}, "DAY_MICROSECOND": {}, "DAY_MINUTE": {}, "DAY_SECOND": {}, "DDL": {}, "DECIMAL": {}, "DEFAULT": {},
	"DELAYED": {}, "DELETE": {}, "DEPTH": {}, "DESC": {}, "DESCRIBE": {}, "DISTINCT": {}, "DISTINCTROW": {}, "DIV": {}, "DOUBLE": {},
	"DRAINER": {}, "DROP": {}, "DUAL": {},

	"ELSE": {}, "ENCLOSED": {}, "ESCAPED": {}, "EXCEPT": {}, "EXISTS": {}, "EXPLAIN": {},

	"FALSE": {}, "FLOAT": {}, "FOR": {}, "FORCE": {}, "FOREIGN": {}, "FROM": {}, "FULLTEXT": {},

	"GENERATED": {}, "GRANT": {}, "GROUP": {}, "HAVING": {}, "HIGH_PRIORITY": {}, "HOUR_MICROSECOND": {}, "HOUR_MINUTE": {}, "HOUR_SECOND": {}, "IF": {},

	"IGNORE": {}, "IN": {}, "INDEX": {}, "INFILE": {}, "INNER": {}, "INSERT": {}, "INT": {}, "INT1": {}, "INT2": {},
	"INT3": {}, "INT4": {}, "INT8": {}, "INTEGER": {}, "INTERVAL": {}, "INTO": {}, "IS": {},

	"JOB": {}, "JOBS": {}, "JOIN": {}, "KEY": {}, "KEYS": {}, "KILL": {},

	"LEADING": {}, "LEFT": {}, "LIKE": {}, "LIMIT": {}, "LINEAR": {}, "LINES": {}, "LOAD": {}, "LOCALTIME": {}, "LOCALTIMESTAMP": {},
	"LOCK": {}, "LONG": {}, "LONGBLOB": {}, "LONGTEXT": {}, "LOW_PRIORITY": {},

	"MATCH": {}, "MAXVALUE": {}, "MEDIUMBLOB": {}, "MEDIUMINT": {}, "MEDIUMTEXT": {}, "MINUTE_MICROSECOND": {}, "MINUTE_SECOND": {}, "MOD": {},

	"NATURAL": {}, "NODE_ID": {}, "NODE_STATE": {}, "NOT": {}, "NO_WRITE_TO_BINLOG": {}, "NULL": {}, "NUMERIC": {},

	"ON": {}, "OPTIMISTIC": {}, "OPTIMIZE": {}, "OPTION": {}, "OPTIONALLY": {}, "OR": {}, "ORDER": {}, "OUTER": {}, "OUTFILE": {},

	"PARTITION": {}, "PESSIMISTIC": {}, "PRECISION": {}, "PRIMARY": {}, "PROCEDURE": {}, "PUMP": {},

	"RANGE": {}, "READ": {}, "REAL": {}, "REFERENCES": {}, "REGEXP": {}, "REGION": {}, "REGIONS": {}, "RELEASE": {}, "RENAME": {},
	"REPEAT": {}, "REPLACE": {}, "REQUIRE": {}, "RESTRICT": {}, "REVOKE": {}, "RIGHT": {}, "RLIKE": {}, "ROW": {},

	"SAMPLES": {}, "SECOND_MICROSECOND": {}, "SELECT": {}, "SET": {}, "SHOW": {}, "SMALLINT": {}, "SPATIAL": {}, "SPLIT": {}, "SQL": {},
	"SQL_BIG_RESULT": {}, "SQL_CALC_FOUND_ROWS": {}, "SQL_SMALL_RESULT": {}, "SSL": {}, "STARTING": {}, "STATS": {}, "STATS_BUCKETS": {}, "STATS_HEALTHY": {}, "STATS_HISTOGRAMS": {},
	"STATS_META": {}, "STORED": {}, "STRAIGHT_JOIN": {},

	"TABLE": {}, "TERMINATED": {}, "THEN": {}, "TIDB": {}, "TIFLASH": {}, "TINYBLOB": {}, "TINYINT": {}, "TINYTEXT": {}, "TO": {},
	"TOPN": {}, "TRAILING": {}, "TRIGGER": {}, "TRUE": {},

	"UNION": {}, "UNIQUE": {}, "UNLOCK": {}, "UNSIGNED": {}, "UPDATE": {}, "USAGE": {}, "USE": {}, "USING": {}, "UTC_DATE": {},
	"UTC_TIME": {}, "UTC_TIMESTAMP": {},

	"VALUES": {}, "VARBINARY": {}, "VARCHAR": {}, "VARCHARACTER": {}, "VARYING": {}, "VIRTUAL": {},

	"WHEN": {}, "WHERE": {}, "WIDTH": {}, "WITH": {}, "WRITE": {}, "XOR": {}, "YEAR_MONTH": {}, "ZEROFILL": {},
}

func IsKeyWord(name string) bool {
	name = strings.ToUpper(name)
	if _, ok := keyWords[name]; ok {
		return true
	}
	return false
}

const varcharLimit = 1024

func VarcharLengthTooLong(nodes []ast.StmtNode) bool {
	for _, tiStmt := range nodes {
		switch node := tiStmt.(type) {
		case *ast.AlterTableStmt:
			for _, v := range node.Specs {
				for _, col := range v.NewColumns {
					if col.Tp.Tp == mysql.TypeVarchar && col.Tp.Flen > varcharLimit {
						return true
					}
				}
			}
		case *ast.CreateTableStmt:
			for _, v := range node.Cols {
				if v.Tp.Tp == mysql.TypeVarchar && v.Tp.Flen > varcharLimit {
					return true
				}
			}
		default:
			return false
		}
	}
	return false
}

func SinglePrimaryKeyIsInt(nodes []ast.StmtNode) bool {
	for _, tiStmt := range nodes {
		switch node := tiStmt.(type) {
		case *ast.CreateTableStmt:
			primaryKeyColName := ""
			for _, v := range node.Constraints {
				if v.Tp == ast.ConstraintPrimaryKey && len(v.Keys) == 1 {
					primaryKeyColName = v.Keys[0].Column.Name.L
				}
			}

			if primaryKeyColName == "" {
				return false
			}

			for _, v := range node.Cols {
				if v.Name.Name.L == primaryKeyColName &&
					(v.Tp.Tp == mysql.TypeShort ||
						v.Tp.Tp == mysql.TypeLong ||
						v.Tp.Tp == mysql.TypeTiny ||
						v.Tp.Tp == mysql.TypeInt24) {
					return true
				}
			}
		default:
			return false
		}
	}
	return false
}

//Rows defines methods that scanner needs, which database/sql.Rows already implements
type Rows interface {
	Close() error

	Columns() ([]string, error)

	Next() bool

	Scan(dest ...interface{}) error
}

func ScanMap(rows Rows) ([]map[string]interface{}, error) {
	return resolveDataFromRows(rows)
}

func resolveDataFromRows(rows Rows) ([]map[string]interface{}, error) {
	if nil == rows {
		return nil, errors.New("[scanner]: rows is nil")
	}
	columns, err := rows.Columns()
	if nil != err {
		return nil, err
	}
	length := len(columns)
	var result []map[string]interface{}
	//unnecessary to put below into rows.Next loop,reduce allocating
	values := make([]interface{}, length)
	for i := 0; i < length; i++ {
		values[i] = new(interface{})
	}
	for rows.Next() {
		err = rows.Scan(values...)
		if nil != err {
			return nil, err
		}
		mp := make(map[string]interface{})
		for idx, name := range columns {
			//mp[name] = reflect.ValueOf(values[idx]).Elem().Interface()
			mp[name] = *(values[idx].(*interface{}))
		}
		result = append(result, mp)
	}
	return result, nil
}
