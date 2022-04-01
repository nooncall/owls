package checker

import (
	"github.com/flipped-aurora/gin-vue-admin/server/service/tidb_or_mysql/task"
)

// Rule 评审规则元数据结构
type Rule struct {
	Item     string `json:"item"`     // 规则代号
	Name     string `json:"name"`     // 规则名称
	Summary  string `json:"summary"`  // 规则摘要
	Content  string `json:"content"`  // 规则解释
	Case     string `json:"case"`     // SQL示例
	Position int    `json:"position"` // 建议所处SQL字符位置，默认0表示全局建议
	Open     bool   `json:"open"`     // 是否禁用，

	CheckFunc func(*Rule, *Audit, *task.DBInfo) (pass bool, newSummary string, affectRows int) `json:"-"` // 函数名
}

var Rules []Rule

func init() {
	Rules = []Rule{
		{
			Item:      RuleCreate001,
			Name:      "injection check character set",
			Summary:   "表必须指定utf8mb4字符集,若指定collate则必须为utf8mb4_bin",
			CheckFunc: (*Rule).RuleCreateTableCharset,
			Open:      true,
		},
		{
			Item:      RuleCreate002,
			Name:      "injection check table comment",
			Summary:   "表必须要有注释",
			CheckFunc: (*Rule).RuleCreateTableComment,
			Open:      true,
		},
		{
			Item:      RuleCreate003,
			Name:      "injection check table primary key",
			Summary:   "表必须要有主键",
			CheckFunc: (*Rule).RuleCreateTableIndex,
			Open:      true,
		},
		{
			Item:      RuleCreate004,
			Name:      "injection check column comment",
			Summary:   "列必须要有注释",
			CheckFunc: (*Rule).RuleColCommentCheck,
			Open:      true,
		},
		{
			Item:      RuleCreate005,
			Name:      "injection check unique index name",
			Summary:   "Unique索引必须要以uniq_为前缀，比如uniq_ab",
			CheckFunc: (*Rule).RuleCreateTableUniqIndex,
			Open:      true,
		},
		{
			Item:      RuleCreate006,
			Name:      "injection check index name",
			Summary:   "普通索引必须要以idx_为前缀，比如idx_ab",
			CheckFunc: (*Rule).RuleCreateTableNormalIndex,
			Open:      true,
		},
		{
			Item:      RuleCreate007,
			Name:      "injection check index column amount",
			Summary:   "索引的列数不能超过3个",
			CheckFunc: (*Rule).RuleCreateTableIndexColNum,
			Open:      true,
		},
		{
			Item:      RuleCreate008,
			Name:      "injection check index amount",
			Summary:   "建表时最多有5个索引",
			CheckFunc: (*Rule).RuleCreateTableIndexNum,
			Open:      true,
		},
		{
			Item:      RuleCreate009,
			Name:      "injection check repeat index",
			Summary:   "不能有重复的索引，包括(索引名不同，字段相同；冗余索引，如(a),(a,b))",
			CheckFunc: (*Rule).RuleCreateTableDupIndex,
			Open:      true,
		},
		{
			Item:      RuleCreate010,
			Name:      "injection check column not null",
			Summary:   "每个列使用not null，text字段除外；text不能有默认值",
			CheckFunc: (*Rule).RuleCreateTableNotNullValue,
			Open:      true,
		},
		{
			Item:      RuleCreate011,
			Name:      "injection check column character set",
			Summary:   "不能设置列的字符集",
			CheckFunc: (*Rule).RuleCreateTableSetColCharset,
			Open:      true,
		},
		{
			Item:      RuleCreate012,
			Name:      "injection check joint index order",
			Summary:   "如果是联合索引，时间字段，不能放在第一列",
			CheckFunc: (*Rule).RuleCreateTableCoIndexOrder,
			Open:      true,
		},
		{
			Item:      RuleCreate013,
			Name:      "injection check joint index exclude",
			Summary:   "联合索引，不能包含主键或者唯一索引列",
			CheckFunc: (*Rule).RuleCreateTableCoIndexEx,
			Open:      true,
		},
		{
			Item:      RuleCreate014,
			Name:      "injection check index column character amount",
			Summary:   "索引列字符量>128拦截",
			CheckFunc: (*Rule).RuleCreateTableIndexLen,
			Open:      true,
		},
		{
			Item:      RuleCreate015,
			Name:      "injection check text column amount",
			Summary:   "表text字段不能超过3个",
			CheckFunc: (*Rule).RuleCreateTableTextColNum,
			Open:      true,
		},
		{
			Item:      RuleCreate016,
			Name:      "injection can't use mysql key words as column name",
			Summary:   "不能使用mysql关键字、保留字作为列名",
			CheckFunc: (*Rule).RuleCreateTableNotUseKeyWorld,
			Open:      true,
		},
		{
			Item:      RuleCreate017,
			Name:      "injection can't use int as primary key, bigInt expect",
			Summary:   "不能使用int作为主键，请换成bigInt",
			CheckFunc: (*Rule).RuleNotUseIntAsPrimaryKey,
			Open:      true,
		},
		{ // 改表时同样适用
			Item:      RuleCreate018,
			Name:      "injection varchar length no more than 1024",
			Summary:   "varchar的长度不能超过1024，建议使用text",
			CheckFunc: (*Rule).RuleVarcharLengthLimit,
			Open:      true,
		},
		{
			Item:      RuleCreate019,
			Name:      "auto increment start value limit 1w",
			Summary:   "自增主键起始值不能超过10000",
			CheckFunc: (*Rule).RuleAutoIncrementLimit,
			Open:      true,
		},
		{
			Item:      RuleAlter001,
			Name:      "injection check drop column",
			Summary:   "禁止删列操作",
			CheckFunc: (*Rule).RuleAlterTableDropColumn,
			Open:      true,
		},
		{
			Item:      RuleAlter002,
			Name:      "injection check drop or truncate table",
			Summary:   "禁止drop、truncate表操作",
			CheckFunc: (*Rule).RuleAlterTableDrop,
			Open:      true,
		},
		// 这个规则，只检查添加索引的场景。（后续可能会有删除索引的场景）
		{
			Item:      RuleAlter003,
			Name:      "injection affect rows",
			Summary:   "影响行数超过3000",
			CheckFunc: (*Rule).RuleAffectRows,
			Open:      true,
		},
		{
			Item:      RuleAlter004,
			Name:      "injection don't add multi index or column once",
			Summary:   "不支持同时创建多个索引或者列",
			CheckFunc: (*Rule).RuleBanAddMulti,
			Open:      true,
		},
		// 跟上一条rule加上关联，认为一次修改仅有一个变动的列
		{
			Item:      RuleAlter005,
			Name:      "injection unsupported modify type",
			Summary:   "不支持的类型更改",
			CheckFunc: (*Rule).RuleUnsupportedType,
			Open:      true,
		},
		{
			Item:      RuleDml001,
			Name:      "injection check where",
			Summary:   "必须有where条件",
			CheckFunc: (*Rule).RuleDMLTableNoWhere,
			Open:      true,
		},
		{
			Item:      RuleDml002,
			Name:      "injection check always true",
			Summary:   "永远为真的比较条件",
			CheckFunc: (*Rule).RuleMeaninglessWhere,
			Open:      true,
		},
		{
			Item:      RuleDml003,
			Name:      "injection check use joined table when update",
			Summary:   "不建议使用联表删除或更新",
			CheckFunc: (*Rule).RuleMultiDeleteUpdate,
			Open:      true,
		},
		{
			Item:      RuleDml004,
			Name:      "injection check insert column def",
			Summary:   "必须指定插入列表",
			CheckFunc: (*Rule).RuleInsertColDef,
			Open:      true,
		},
		{
			Item:      RuleDml005,
			Name:      "injection check insert value amount",
			Summary:   "插入列列表与值列表个数不相同",
			CheckFunc: (*Rule).RuleInsertColValueEqual,
			Open:      true,
		},
		{
			Item:      RuleDml006,
			Name:      "injection check table exist",
			Summary:   "目标库中表不存在",
			CheckFunc: (*Rule).RuleAlterTableExist,
			Open:      true,
		},
		{
			Item:      RuleDml007,
			Name:      "injection check column exist",
			Summary:   "目标库中表的字段不存在",
			CheckFunc: (*Rule).RuleAlterTableColumnExist,
			Open:      true,
		},
		{
			Item:      RuleDml008,
			Name:      "injection index match ",
			Summary:   "删改数据索引不完全匹配",
			CheckFunc: (*Rule).RuleDmlIndexMatch,
			Open:      true,
		},
		{
			Item:      RuleDml009,
			Name:      "injection affect no more than 100",
			Summary:   "删改数据影响行数超过100",
			CheckFunc: (*Rule).RuleDmlNoMoreThan100,
			Open:      true,
		},
	}
}

const (
	RuleCreate001 = "Create.001"
	RuleCreate002 = "Create.002"
	RuleCreate003 = "Create.003"
	RuleCreate004 = "Create.004"
	RuleCreate005 = "Create.005"
	RuleCreate006 = "Create.006"
	RuleCreate007 = "Create.007"
	RuleCreate008 = "Create.008"
	RuleCreate009 = "Create.009"
	RuleCreate010 = "Create.010"
	RuleCreate011 = "Create.011"
	RuleCreate012 = "Create.012"
	RuleCreate013 = "Create.013"
	RuleCreate014 = "Create.014"
	RuleCreate015 = "Create.015"
	RuleCreate016 = "Create.016"
	RuleCreate017 = "Create.017"
	RuleCreate018 = "Create.018"
	RuleCreate019 = "Create.019"

	RuleAlter001 = "Alter.001"
	RuleAlter002 = "Alter.002"
	RuleAlter003 = "Alter.003"
	RuleAlter004 = "Alter.004"
	RuleAlter005 = "Alter.005"

	RuleDml001 = "DML.001"
	RuleDml002 = "DML.002"
	RuleDml003 = "DML.003"
	RuleDml004 = "DML.004"
	RuleDml005 = "DML.005"
	RuleDml006 = "DML.006"
	RuleDml007 = "DML.007"
	RuleDml008 = "DML.008"
	RuleDml009 = "DML.009"
)

var breakRules []string = []string{
	RuleDml008,
	RuleAlter004,
}

func IsBreakRule(name string) bool {
	for _, v := range breakRules {
		if v == name {
			return true
		}
	}
	return false
}
