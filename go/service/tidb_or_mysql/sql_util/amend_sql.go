package sql_util

import (
	"fmt"

	"github.com/pingcap/parser/ast"
	"github.com/qingfeng777/owls/server/utils/logger"
)

func AddLimit(sql string) string {
	f := "AddLimit-->: "
	stmtNodes, _, err := getParser().Parse(sql, "", "")
	if err != nil {
		logger.Errorf("%sparse sql err: %v", f, err)
		return sql
	}

	for _, tiStmt := range stmtNodes {
		switch node := tiStmt.(type) {
		case *ast.SelectStmt:
			if node.Limit == nil {
				return fmt.Sprintf("%s limit 10", sql)
			}
		default:
			return sql
		}
	}

	return sql
}
