package checker

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/mysql"
	"github.com/pingcap/parser/types"

	"vitess.io/vitess/go/vt/sqlparser"

	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/sql_util"
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/task"
	"github.com/flipped-aurora/gin-vue-admin/server/utils/logger"
)

// RuleOK OK
func (q *Rule) RuleOK(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	return true, q.Summary, 0
}

//RuleCreateTableCharset Create.001
func (q *Rule) RuleCreateTableCharset(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var charSetIsMB4, collateIsNotMB4 bool
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for k := range node.Options {
					//must be utf8mb4
					if node.Options[k].Tp == ast.TableOptionCharset && node.Options[k].StrValue == "utf8mb4" {
						charSetIsMB4 = true
					}
					//utf8mb4_bin or no value,
					if node.Options[k].Tp == ast.TableOptionCollate && node.Options[k].StrValue != "utf8mb4_bin" {
						collateIsNotMB4 = true
					}
				}
				if charSetIsMB4 && !collateIsNotMB4 {
					return true, q.Summary, 0
				}
				return false, q.Summary, 0
			}
		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableComment Create.002
func (q *Rule) RuleCreateTableComment(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for k := range node.Options {
					if node.Options[k].Tp == ast.TableOptionComment && node.Options[k].StrValue != "" {
						return true, q.Summary, 0
					}
				}
				return false, q.Summary, 0
			}
		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableIndex Create.003
func (q *Rule) RuleCreateTableIndex(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for k := range node.Constraints {
					if node.Constraints[k].Tp == ast.ConstraintPrimaryKey && len(node.Constraints[k].Keys) > 0 {
						return true, q.Summary, 0
					}
				}
				return false, q.Summary, 0
			}
		}
	}
	return true, q.Summary, 0
}

// RuleColCommentCheck Create.004
func (q *Rule) RuleColCommentCheck(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, node := range audit.TiStmt {
		switch n := node.(type) {
		case *ast.CreateTableStmt:
			for _, c := range n.Cols {
				colComment := false
				for _, o := range c.Options {
					if o.Tp == ast.ColumnOptionComment {
						colComment = true
					}
				}
				if !colComment {
					return false, q.Summary, 0
				}
			}
		case *ast.AlterTableStmt:
			for _, s := range n.Specs {
				switch s.Tp {
				case ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn:
					for _, c := range s.NewColumns {
						colComment := false
						for _, o := range c.Options {
							if o.Tp == ast.ColumnOptionComment {
								colComment = true
							}
						}
						if !colComment {
							return false, q.Summary, 0
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableUniqIndex Create.005
func (q *Rule) RuleCreateTableUniqIndex(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var prefix = "uniq_"
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
						if !strings.HasPrefix(v.Name, prefix) {
							return false, q.Summary, 0
						}
					}
				}
			case *ast.AlterTableStmt:
				for _, v := range node.Specs {
					switch v.Tp {
					case ast.AlterTableAddConstraint:
						switch v.Constraint.Tp {
						case ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex:
							if !strings.HasPrefix(v.Constraint.Name, prefix) {
								return false, q.Summary, 0
							}
						}
					}
				}
			}

		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableNormalIndex Create.006
func (q *Rule) RuleCreateTableNormalIndex(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var prefix = "idx_"
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintIndex, ast.ConstraintKey:
						if !strings.HasPrefix(v.Name, prefix) {
							return false, q.Summary, 0
						}
					}
				}
			case *ast.AlterTableStmt:
				for _, v := range node.Specs {
					switch v.Tp {
					case ast.AlterTableAddConstraint:
						switch v.Constraint.Tp {
						case ast.ConstraintIndex, ast.ConstraintKey:
							if !strings.HasPrefix(v.Constraint.Name, prefix) {
								return false, q.Summary, 0
							}
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableIndexNum Create.007
func (q *Rule) RuleCreateTableIndexColNum(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var maxIndexColNum = 3
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintPrimaryKey:
						if len(v.Keys) > maxIndexColNum {
							return false, q.Summary, 0
						}
						if v.Refer != nil && len(v.Refer.IndexPartSpecifications) > maxIndexColNum {
							return false, q.Summary, 0
						}
					}
				}
			case *ast.AlterTableStmt:
				for _, v := range node.Specs {
					switch v.Tp {
					case ast.AlterTableAddConstraint:
						switch v.Constraint.Tp {
						case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintPrimaryKey:
							if len(v.Constraint.Keys) > maxIndexColNum {
								return false, q.Summary, 0
							}
							if v.Constraint.Refer != nil && len(v.Constraint.Refer.IndexPartSpecifications) > maxIndexColNum {
								return false, q.Summary, 0
							}
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableIndexNum Create.008
func (q *Rule) RuleCreateTableIndexNum(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var maxIndexNum = 5
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				var indexNum int
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintPrimaryKey:
						indexNum++
					}
				}
				if indexNum > maxIndexNum {
					return false, q.Summary, 0
				}
			}
		}
	}
	return true, q.Summary, 0
}

//RuleCreateTableDupIndex Create.009
func (q *Rule) RuleCreateTableDupIndex(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				var idxmap = make(map[string]bool)
				for _, v := range node.Constraints {
					var keys = make([]string, len(v.Keys))
					for k, v1 := range v.Keys {
						keys[k] = fmt.Sprintf("%v", v1.Column)
					}
					key := strings.Join(keys, sql_util.KeyJoinChar)
					if idxmap[key] {
						return false, q.Summary, 0
					} else {
						idxmap[key] = true
					}
				}
				for k := range idxmap {
					for k1 := range idxmap {
						if k == k1 {
							continue
						}
						if strings.HasPrefix(k, k1) && sql_util.IsSubKey(k, k1) {
							return false, q.Summary, 0
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleAddNotNullValue Create.010
func (q *Rule) RuleCreateTableNotNullValue(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, node := range audit.TiStmt {
		switch n := node.(type) {
		case *ast.CreateTableStmt:
			for _, c := range n.Cols {
				if c.Tp.Tp == mysql.TypeBlob || c.Tp.Tp == mysql.TypeTinyBlob || c.Tp.Tp == mysql.TypeMediumBlob || c.Tp.Tp == mysql.TypeLongBlob || c.Tp.Tp == mysql.TypeJSON {
					for _, o := range c.Options {
						if o.Tp == ast.ColumnOptionDefaultValue {
							return false, q.Summary, 0
						}
					}
					continue
				}
				var bSetNotNull bool
				for _, o := range c.Options {
					if o.Tp == ast.ColumnOptionNotNull {
						bSetNotNull = true
					}
				}
				if !bSetNotNull {
					return false, q.Summary, 0
				}
			}
		case *ast.AlterTableStmt:
			for _, s := range n.Specs {
				switch s.Tp {
				case ast.AlterTableAddColumns, ast.AlterTableChangeColumn, ast.AlterTableModifyColumn:
					for _, c := range s.NewColumns {
						if c.Tp.Tp == mysql.TypeBlob || c.Tp.Tp == mysql.TypeTinyBlob || c.Tp.Tp == mysql.TypeMediumBlob || c.Tp.Tp == mysql.TypeLongBlob || c.Tp.Tp == mysql.TypeJSON {
							for _, o := range c.Options {
								if o.Tp == ast.ColumnOptionDefaultValue {
									return false, q.Summary, 0
								}
							}
							continue
						}
						var bSetNotNull bool
						for _, o := range c.Options {
							if o.Tp == ast.ColumnOptionNotNull {
								bSetNotNull = true
							}
						}
						if !bSetNotNull {
							return false, q.Summary, 0
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableSetColCharset Create.011
func (q *Rule) RuleCreateTableSetColCharset(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for _, col := range node.Cols {
					if col.Tp == nil || col.Tp.Tp == mysql.TypeJSON {
						continue
					}
					if col.Tp.Charset != "" || col.Tp.Collate != "" {
						return false, q.Summary, 0
					}
					for _, v := range col.Options {
						if v.Tp == ast.ColumnOptionCollate {
							return false, q.Summary, 0
						}
					}
				}
			case *ast.AlterTableStmt:
				for _, spec := range node.Specs {
					switch spec.Tp {
					case ast.AlterTableAlterColumn, ast.AlterTableChangeColumn,
						ast.AlterTableModifyColumn, ast.AlterTableAddColumns:
						for _, col := range spec.NewColumns {
							if col.Tp == nil || col.Tp.Tp == mysql.TypeJSON {
								continue
							}
							if col.Tp.Charset != "" || col.Tp.Collate != "" {
								return false, q.Summary, 0
							}
							for _, v := range col.Options {
								if v.Tp == ast.ColumnOptionCollate {
									return false, q.Summary, 0
								}
							}
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableCoIndexOrder Create.012
func (q *Rule) RuleCreateTableCoIndexOrder(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				var colType = make(map[ast.ColumnName]*types.FieldType, len(node.Cols))
				for _, v := range node.Cols {
					colType[*v.Name] = v.Tp
				}
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintPrimaryKey:
						if len(v.Keys) <= 1 {
							continue
						}
						// colType[*v.Keys[0].Column] 为空，下面调 Tp 直接空指针了
						// colType[*v.Keys[0].Column] 之所以为空是因为有人传了不存在的key进来
						// 例如，KEY `idx_a_b` (`a`, `b`)，但是建表的过程中根本没有a这个列
						_, ok := colType[*v.Keys[0].Column]
						if !ok {
							return false, "key引用了不存在的字段: " + v.Keys[0].Column.String(), 0
						}
						if colType[*v.Keys[0].Column].Tp == mysql.TypeTimestamp || colType[*v.Keys[0].Column].Tp == mysql.TypeDatetime {
							return false, q.Summary, 0
						}

					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableCoIndexEx Create.013
func (q *Rule) RuleCreateTableCoIndexEx(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				var colType = make(map[string]bool)
				for _, v := range node.Constraints {
					if v.Tp == ast.ConstraintPrimaryKey || v.Tp == ast.ConstraintUniq {
						var keys []string
						for _, v1 := range v.Keys {
							keys = append(keys, fmt.Sprintf("%v", v1.Column))
						}
						key := strings.Join(keys, "_")
						colType[key] = true
					}
				}
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintIndex, ast.ConstraintKey:
						fmt.Printf("col type:%+v\n", colType)
						if len(v.Keys) <= 1 {
							continue
						}
						var keys []string
						for _, v1 := range v.Keys {
							keys = append(keys, fmt.Sprintf("%v", v1.Column))
						}
						key := strings.Join(keys, "_")
						fmt.Printf("key :%+v\n", key)

						for k1 := range colType {
							if strings.HasPrefix(key, k1) {
								return false, q.Summary, 0
							}
						}
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableIndexLen Create.014
func (q *Rule) RuleCreateTableIndexLen(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var lenth = 128
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				var indexCol = make(map[string]bool)
				for _, v := range node.Constraints {
					switch v.Tp {
					case ast.ConstraintIndex, ast.ConstraintKey, ast.ConstraintUniq, ast.ConstraintUniqKey, ast.ConstraintUniqIndex, ast.ConstraintPrimaryKey:
						for _, v1 := range v.Keys {
							col := fmt.Sprintf("%v", v1.Column)
							indexCol[col] = true
						}
					}
				}
				for _, v := range node.Cols {
					colName := fmt.Sprintf("%v", v.Name)
					if indexCol[colName] && v.Tp.Flen > lenth && v.Tp.Tp == mysql.TypeVarchar {
						return false, q.Summary, 0
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableTextColNum Create.015
func (q *Rule) RuleCreateTableTextColNum(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var maxTextColNum = 3
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				var textColNum int
				for _, v := range node.Cols {
					if v.Tp.Tp == mysql.TypeBlob {
						textColNum++
					}
				}
				if textColNum > maxTextColNum {
					return false, q.Summary, 0
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableTextColNum Create.016
func (q *Rule) RuleCreateTableNotUseKeyWorld(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.CreateTableStmt:
				for _, v := range node.Cols {
					if sql_util.IsKeyWord(v.Name.Name.O) {
						return false, q.Summary + " key: " + v.Name.Name.O, 0
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleCreateTableTextColNum Create.017
func (q *Rule) RuleNotUseIntAsPrimaryKey(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	if sql_util.SinglePrimaryKeyIsInt(audit.TiStmt) {
		return false, q.Summary, 0
	}
	return true, q.Summary, 0
}

// RuleCreateTableTextColNum Create.018
func (q *Rule) RuleVarcharLengthLimit(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	if sql_util.VarcharLengthTooLong(audit.TiStmt) {
		return false, q.Summary, 0
	}
	return true, q.Summary, 0
}

const autoIncrementLimit = 10000

func (q *Rule) RuleAutoIncrementLimit(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, node := range audit.TiStmt {
		switch n := node.(type) {
		case *ast.CreateTableStmt:
			for _, v := range n.Options {
				if v.Tp == ast.TableOptionAutoIncrement && v.UintValue > autoIncrementLimit {
					return false, q.Summary, 0
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleAlterTableDropColumn Alter.001
func (q *Rule) RuleAlterTableDropColumn(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch audit.Stmt.(type) {
	case *sqlparser.DDL:
		for _, tiStmt := range audit.TiStmt {
			switch node := tiStmt.(type) {
			case *ast.AlterTableStmt:
				for _, spec := range node.Specs {
					switch spec.Tp {
					case ast.AlterTableDropColumn:
						return false, q.Summary, 0
					}
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleAlterTableDrop Alter.002
func (q *Rule) RuleAlterTableDrop(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, tiStmt := range audit.TiStmt {
		switch tiStmt.(type) {
		case *ast.DropTableStmt, *ast.TruncateTableStmt:
			return false, q.Summary, 0
		}
	}
	return true, q.Summary, 0
}

// Alter.004
func (q *Rule) RuleBanAddMulti(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, tiStmt := range audit.TiStmt {
		switch tiStmt.(type) {
		case *ast.AlterTableStmt:
			node := tiStmt.(*ast.AlterTableStmt)
			// 仅判断了alter 语句下面有超过一个子项
			if len(node.Specs) > 1 {
				return false, q.Summary, 0
			}
		}
	}
	return true, q.Summary, 0
}

// Alter.005
func (q *Rule) RuleUnsupportedType(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, tiStmt := range audit.TiStmt {
		switch node := tiStmt.(type) {
		case *ast.AlterTableStmt:
			if len(node.Specs) < 1 {
				return true, q.Summary, 0
			}
			if node.Specs[0].Tp == ast.AlterTableChangeColumn || node.Specs[0].Tp == ast.AlterTableModifyColumn {
				if unsupportedTypeChange(audit.Query, tiStmt, info) {
					return false, q.Summary, 0
				}
			}
		}
	}
	return true, q.Summary, 0
}

//RuleDMLTableNoWhere DML.001
func (q *Rule) RuleDMLTableNoWhere(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	pass = true
	sqlparser.Walk(func(node sqlparser.SQLNode) (kontinue bool, err error) {
		switch n := node.(type) {
		case *sqlparser.Delete:
			if n.Where == nil {
				pass = false
				return false, nil
			}
		case *sqlparser.Update:
			if n.Where == nil {
				pass = false
				return false, nil
			}
		}
		return true, nil
	}, audit.Stmt)

	return pass, q.Summary, 0
}

// RuleMeaninglessWhere DML.002
func (q *Rule) RuleMeaninglessWhere(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	pass = true
	sqlparser.Walk(func(node sqlparser.SQLNode) (continueWalk bool, err error) {
		switch n := node.(type) {
		case *sqlparser.ComparisonExpr:
			factor := false
			switch n.Operator {
			case "!=", "<>":
				factor = true
			case "=", "<=>":
			default:
				return true, nil
			}
			var left []byte
			var right []byte
			// left
			switch l := n.Left.(type) {
			case *sqlparser.SQLVal:
				left = l.Val
			default:
				return true, nil
			}

			// right
			switch r := n.Right.(type) {
			case *sqlparser.SQLVal:
				right = r.Val
			default:
				return true, nil
			}

			// compare
			if (bytes.Equal(left, right) && !factor) || (!bytes.Equal(left, right) && factor) {
				pass = false
			}
			return false, nil
		}
		return true, nil
	}, audit.Stmt)
	return pass, q.Summary, 0
}

// RuleMultiDeleteUpdate DML.003
func (q *Rule) RuleMultiDeleteUpdate(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	pass = true
	switch audit.Stmt.(type) {
	case *sqlparser.Delete, *sqlparser.Update:
		sqlparser.Walk(func(node sqlparser.SQLNode) (continueWalk bool, err error) {
			switch node.(type) {
			case *sqlparser.JoinTableExpr:
				pass = false
				return false, nil
			case *sqlparser.SelectExprs:
				pass = false
				return false, nil
			}
			return true, nil
		}, audit.Stmt)
	}
	return pass, q.Summary, 0
}

// RuleInsertColDef DML.004
func (q *Rule) RuleInsertColDef(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	for _, tiStmt := range audit.TiStmt {
		switch n := tiStmt.(type) {
		case *ast.InsertStmt:
			if n.Columns == nil {
				return false, q.Summary, 0
			}
		}
	}
	return true, q.Summary, 0
}

// RuleInsertColValueEqual DML.005
func (q *Rule) RuleInsertColValueEqual(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	switch node := audit.Stmt.(type) {
	case *sqlparser.Insert:
		if node.Columns == nil {
			return false, q.Summary, 0
		}
		colLen := len(node.Columns)
		switch val := node.Rows.(type) {
		case sqlparser.Values:
			for k := range val {
				l := len(val[k])
				if l != colLen {
					return false, q.Summary, 0
				}
			}
		}
	}
	return true, q.Summary, 0
}

// RuleAlterTableExist DML.006
func (q *Rule) RuleAlterTableExist(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var tbs []string
	for _, tiStmt := range audit.TiStmt {
		switch n := tiStmt.(type) {
		case *ast.AlterTableStmt:
			tbs = append(tbs, n.Table.Name.String())
		case *ast.DropTableStmt, *ast.TruncateTableStmt:
			return false, q.Summary, 0
		case *ast.DeleteStmt:
			tbs = append(tbs, n.TableRefs.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
		case *ast.InsertStmt:
			tbs = append(tbs, n.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
		case *ast.UpdateStmt:
			switch n.TableRefs.TableRefs.Left.(type) {
			case *ast.Join:
				tbs = append(tbs, n.TableRefs.TableRefs.Left.(*ast.Join).Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
			default:
				tbs = append(tbs, n.TableRefs.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
			}
		default:
			return true, q.Summary, 0
		}
	}
	if len(tbs) != 1 {
		return false, q.Summary, 0
	}
	ts, err := getTableSysInfo(tbs[0], info)
	if err != nil || ts == nil {
		logger.Infof("get table sys info err:%+v, ts:%+v, task.DBInfo:%+v", err, ts, info)
		return false, q.Summary, 0
	}
	return true, q.Summary, 0
}

// RuleAlterTableColumnExist DML.007
func (q *Rule) RuleAlterTableColumnExist(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var cols []string
	for _, tiStmt := range audit.TiStmt {
		switch n := tiStmt.(type) {
		case *ast.AlterTableStmt:
			for _, v := range n.Specs {
				if v.Tp == ast.AlterTableAddColumns {
					continue
				}
				if v.OldColumnName != nil {
					cols = append(cols, v.OldColumnName.Name.String())
				} else {
					for _, v1 := range v.NewColumns {
						cols = append(cols, v1.Name.Name.String())
					}
				}
			}
		case *ast.InsertStmt:
			for _, v := range n.Columns {
				cols = append(cols, v.Name.String())
			}
		default:
			return true, q.Summary, 0
		}
	}

	//todo, rewrite
	var tbs []string
	for _, tiStmt := range audit.TiStmt {
		switch n := tiStmt.(type) {
		case *ast.AlterTableStmt:
			tbs = append(tbs, n.Table.Name.String())
		case *ast.DropTableStmt, *ast.TruncateTableStmt:
			return false, q.Summary, 0
		case *ast.DeleteStmt:
			tbs = append(tbs, n.TableRefs.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
		case *ast.InsertStmt:
			tbs = append(tbs, n.Table.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
		case *ast.UpdateStmt:
			switch n.TableRefs.TableRefs.Left.(type) {
			case *ast.Join:
				tbs = append(tbs, n.TableRefs.TableRefs.Left.(*ast.Join).Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
			default:
				tbs = append(tbs, n.TableRefs.TableRefs.Left.(*ast.TableSource).Source.(*ast.TableName).Name.String())
			}
		}
	}

	if len(tbs) != 1 {
		return false, q.Summary, 0
	}
	columns, err := readTableStruct(tbs[0], info)
	if err != nil {
		logger.Infof("Select table columns err:%+v", err.Error())
		return false, q.Summary, 0
	}

	for _, v := range cols {
		var exist bool
		for _, v1 := range columns {
			if v == v1.Name {
				exist = true
				break
			}
		}
		if !exist {
			logger.Infof("get table columns info err, cols:%+v, columns:%+v, task.DBInfo:%+v", cols, columns, info)
			return false, q.Summary, 0
		}
	}
	return true, q.Summary, 0
}

func (q *Rule) RuleAffectRows(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var tableName string
	for _, tiStmt := range audit.TiStmt {
		switch n := tiStmt.(type) {
		case *ast.AlterTableStmt:
			if !IsAlterIndexOperate(audit.Query) {
				return true, q.Summary, 0
			}
			tableName = n.Table.Name.String()
		case *ast.CreateIndexStmt:
			tableName = n.Table.Name.String()
		case *ast.DropIndexStmt:
			tableName = n.Table.Name.String()
		default:
			return true, q.Summary, 0
		}

		tableSysInfo, err := getTableSysInfo(tableName, info)
		if err == nil && tableSysInfo != nil && tableSysInfo.TableRows > 0 {
			affectRows = tableSysInfo.TableRows
		}
	}
	return true, q.Summary, affectRows
}

// IndexMatch DML008
func (q *Rule) RuleDmlIndexMatch(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var err error
	var match bool
	for _, tiStmt := range audit.TiStmt {
		switch tiStmt.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			match, err = IndexMach(info, audit.Query)
			if err != nil {
				continue
			}
		default:
			return true, q.Summary, affectRows
		}
		if !match {
			logger.Infof("sql check update %s, index not match", RuleDml008)
			return false, q.Summary, affectRows
		}
	}
	return true, q.Summary, 0
}

const _dmlMaxAllAffectRow = 100

// IndexMatch DML009
func (q *Rule) RuleDmlNoMoreThan100(audit *Audit, info *task.DBInfo) (pass bool, newSummary string, affectRows int) {
	var err error
	var affectRow int
	for _, tiStmt := range audit.TiStmt {
		switch tiStmt.(type) {
		case *ast.UpdateStmt, *ast.DeleteStmt:
			affectRow, err = getAffectRow(info, audit.Query)
			if err != nil {
				continue
			}
		default:
			return true, q.Summary, affectRows
		}
		logger.Infof("%s sql check dml affect row, affectRow: %d", RuleDml009, affectRow)

		if affectRow > 0 {
			affectRows = affectRow
		}

		if affectRows > _dmlMaxAllAffectRow {
			return false, q.Summary, affectRows
		}
	}
	return true, q.Summary, affectRows
}

type column struct {
	Name    string `bdb:"COLUMN_NAME"`
	Type    string `bdb:"COLUMN_TYPE"`
	Comment string `bdb:"COLUMN_COMMENT"`
}

func readTableStruct(table string, info *task.DBInfo) ([]column, error) {
	sql := fmt.Sprintf("select COLUMN_NAME, COLUMN_TYPE, COLUMN_COMMENT from %s where TABLE_NAME = '%s' and TABLE_SCHEMA = '%s'",
		"COLUMNS", table, info.DBName)
	rows, err := info.DefaultDB.Query(sql)
	if nil != err {
		return nil, err
	}
	defer rows.Close()
	var ts []column
	for rows.Next() {
		var col column
		if err = rows.Scan(&col.Name, &col.Type, &col.Comment); nil != err {
			return nil, err
		}
		ts = append(ts, col)
	}
	return ts, nil
}

type tableSysInfo struct {
	TableName string `bdb:"TABLE_NAME"`
	TableRows int    `bdb:"TABLE_ROWS"`
}

func getTableSysInfo(table string, info *task.DBInfo) (*tableSysInfo, error) {
	db := info.DB
	sqlContent := fmt.Sprintf("select TABLE_NAME,TABLE_ROWS from INFORMATION_SCHEMA.TABLES where TABLE_SCHEMA='%s' and TABLE_NAME='%s';", info.DBName, table)
	res, err := db.QueryContext(context.TODO(), sqlContent)
	if err != nil {
		logger.Infof("get table sys info err:%+v", err.Error())
		return nil, err
	}
	defer res.Close()

	for res.Next() {
		var ts tableSysInfo
		if err = res.Scan(&ts.TableName, &ts.TableRows); nil != err {
			return nil, err
		}

		if ts.TableName == table {
			return &ts, nil
		}
	}

	return nil, errors.New("table not found")
}

// index str : unique,primary ,fulltext + index/key
// alter语句，判断是不是对索引的操作
func IsAlterIndexOperate(sql string) bool {
	sql = strings.ToLower(sql)
	items := strings.Split(sql, " ")
	length := len(items)
	for i, v := range items {
		if v == "add" || v == "drop" {
			if i <= length-2 {
				next := i + 1
				if items[next] == "unique" ||
					items[next] == "fulltext" ||
					items[next] == "primary" ||
					items[next] == "index" ||
					items[next] == "key" {
					return true
				}
			}
		}
	}
	return false
}
