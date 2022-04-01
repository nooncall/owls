package checker

import (
	"fmt"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

// Audit 待评审的SQL结构体，由原SQL和其对应的抽象语法树组成
type Audit struct {
	Query  string              // 查询语句
	Stmt   sqlparser.Statement // 通过Vitess解析出的抽象语法树
	TiStmt []ast.StmtNode      // 通过TiDB解析出的抽象语法树
}

type CheckerService struct {
}

var Checker CheckerService

func (CheckerService) SqlCheck(sql, charset, collation string, info *task.DBInfo) (pass bool, suggestion string, affectRow int) {
	audit, err := NewAudit(sql, charset, collation)
	if err != nil {
		return false, fmt.Sprintf("sql解析错误：%s", err.Error()), 0
	}

	pass = true
	for _, v := range Rules {
		if !v.Open {
			continue
		}

		onePass, summary, oneAffectRow := v.CheckFunc(&v, audit, info)
		if affectRow < oneAffectRow {
			affectRow = oneAffectRow
		}
		if !onePass {
			pass = false
			suggestion = summary
			break
		}
	}

	return pass, suggestion, affectRow
}

func (CheckerService) ListRules() interface{} {
	return Rules
}

// NewQuery4Audit return a struct for Audit
func NewAudit(sql, charset, collation string) (*Audit, error) {
	q := &Audit{Query: sql}
	// vitess 语法解析不上报，以 tidb parser 为主
	q.Stmt, _ = sqlparser.Parse(sql)

	// tdib parser 语法解析
	var warns []error
	var err error
	q.TiStmt, warns, err = parser.New().Parse(sql, charset, collation)
	if len(warns) > 0 {
		logger.Warn("parse sql warning: ", utils.ErrsJoin("; ", warns))
	}
	return q, err
}
