package sql_util

import (
	"fmt"
	"strings"

	"github.com/nooncall/owls/go/utils/logger"
	"github.com/pingcap/parser/ast"
)

func AddLimitToSelect(sql string) string {
	f := "AddLimitToSelect-->: "
	stmtNodes, _, err := getParser().Parse(sql, "", "")
	if err != nil {
		logger.Errorf("%sparse sql err: %v", f, err)
		return sql
	}

	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.SelectStmt:
			if node.Limit == nil {
				return appendLimit(sql)
			}
		default:
			return sql
		}
	}

	return sql
}

func appendLimit(sql string) string {
	sql = strings.TrimSpace(sql)
	if string(sql[len(sql)-1]) == ";" {
		sql = sql[:len(sql)-1]
	}

	return fmt.Sprintf("%s limit 10", sql)
}
