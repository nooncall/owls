package sql_util

import (
	"fmt"
	"strings"

	"github.com/nooncall/owls/go/utils/logger"
	"github.com/pingcap/parser/ast"
)

type readSql struct {
	sql string
}

func NewReadSql(sql string) *readSql {
	return &readSql{sql}
}

func (readSql *readSql) String() string {
	return readSql.sql
}

func (readSql *readSql) appendLimit() {
	readSql.sql = fmt.Sprintf("%s limit 10", readSql.sql)
}

func (readSql *readSql) Trim() {
	readSql.sql = strings.TrimSpace(readSql.sql)
	readSql.sql = strings.TrimRight(readSql.sql, ";")
}

func (readSql *readSql) SetLimitResult() string {
	f := "AddLimitToSelect-->: "

	readSql.Trim()
	stmtNodes, _, err := getParser().Parse(readSql.sql, "", "")
	if err != nil {
		logger.Errorf("%sparse sql err: %v", f, err)
		return readSql.sql
	}

	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.SelectStmt:
			if node.Limit == nil {
				readSql.appendLimit()
				return readSql.String()
			}
		default:
			return readSql.String()
		}
	}

	return readSql.String()
}
