package task

import "fmt"

type OwlSubtask struct {
	ID          int64  `json:"id" gorm:"column:id"`
	TaskID      int64  `json:"task_id" gorm:"column:task_id"`
	TaskType    string `json:"task_type" gorm:"column:task_type"`
	DbName      string `json:"db_name" gorm:"column:db_name"`
	ClusterName string `json:"cluster_name" gorm:"column:cluster_name"`

	ExecItems []OwlExecItem `json:"exec_items" gorm:"-"`
}

type OwlExecItem struct {
	ID           int64  `json:"id" gorm:"column:id"`
	TaskID       int64  `json:"task_id" gorm:"column:task_id"`
	SubtaskID    int64  `json:"subtask_id" gorm:"column:subtask_id"`
	SQLContent   string `json:"sql_content" gorm:"column:sql_content"`
	Remark       string `json:"remark" gorm:"column:remark"`
	AffectRows   int    `json:"affect_rows" gorm:"column:affect_rows"`
	RuleComments string `json:"rule_comments" gorm:"column:rule_comments"`
	Status       string `json:"status" gorm:"column:status"`
	ExecInfo     string `json:"exec_info" gorm:"column:exec_info"`
	BackupInfo   string `json:"backup_info" gorm:"column:backup_info"`
	BackupStatus string `json:"backup_status" gorm:"column:backup_status"`
	BackupID     int64  `json:"backup_id" gorm:"column:backup_id"`

	Ut int64 `json:"ut" gorm:"column:ut"`
	Et int64 `json:"et" gorm:"column:et"`

	DBName      string `json:"db_name" gorm:"-"`
	ClusterName string `json:"cluster_name" gorm:"-"`
	TaskType    string `json:"task_type" gorm:"-"`
}

type Status = string

const (
	//顺序递进
	CheckFailed Status = "checkFailed"
	CheckPass          = "checkPass"
	ReviewPass         = "reviewPass"
	DBAPass            = "dbaPass"
	ExecWait           = "execWait"
	Executing          = "executing"
	ExecFailed         = "execFailed"
	ExecSuccess        = "execSuccess"

	//终止状态
	Reject     Status = "reject"
	Cancel            = "cancel"
	ExecCancel        = "execCancel"

	//子项状态
	SkipExec Status = "skipExec"
)

func StatusName(status Status) string {
	switch status {
	case CheckFailed:
		return "系统检测失败"
	case CheckPass:
		return "待审核"
	case ReviewPass, ExecWait:
		return "待执行"
	case Executing:
		return "执行中"
	case ExecFailed:
		return "执行失败"
	case ExecSuccess:
		return "执行成功"
	case Reject:
		return "驳回"
	case Cancel:
		return "撤销"
	case ExecCancel:
		return "撤销执行"
	default:
		return fmt.Sprintf("未知状态:%s", status)
	}
}

type ItemStatus = string

const (
	ItemFailed      ItemStatus = "failed"
	ItemCheckFailed            = "check_failed"
	ItemCheckPass              = "check_pass"
	ItemSuccess                = "success"
	ItemSkipped                = "skipped"

	ItemNoBackup        ItemStatus = "no_backup"
	ItemBackupSuccess              = "backup_success"
	ItemBackupFailed               = "backup_failed"
	ItemRollBackFailed             = "roll_backup_failed"
	ItemRollBackSuccess            = "roll_backup_success"
)

type TaskType = string

const (
	DML       TaskType = "DML"
	DDLCreate          = "CREATE"
	DDLUpdate          = "UPDATE"
)

func GetTypeName(taskType string) string {
	switch taskType {
	case DML:
		return "变更数据"
	case DDLCreate:
		return "建表"
	case DDLUpdate:
		return "改表"
	default:
		return "未知类型"
	}
}

type Action = string

const (
	EditItem Action = "editItem"
	DelItem         = "delItem"
	DoCancel        = "cancel"
	SkipAt          = "skipAt"
	BeginAt         = "beginAt"
	Progress        = "progress"
	DoReject        = "reject"
)

func HistoryStatus() []ItemStatus {
	return []ItemStatus{Reject, Cancel, ExecCancel, ExecFailed, ExecSuccess}
}

func ReviewerStatus() []ItemStatus {
	return []ItemStatus{CheckPass, CheckFailed}
}

func ExecStatus() []ItemStatus {
	return []ItemStatus{ReviewPass, DBAPass, ExecCancel, Executing, ExecWait}
}
