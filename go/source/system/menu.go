package system

import (
	"github.com/pkg/errors"
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/model/system"
	"gorm.io/gorm"
)

var BaseMenu = new(menu)

type menu struct{}

func (m *menu) TableName() string {
	return "sys_base_menus"
}

func (m *menu) Initialize() error {
	entities := []system.SysBaseMenu{
		{GVA_MODEL: global.GVA_MODEL{ID: 1}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "dashboard", Name: "dashboard", Component: "view/dashboard/index.vue", Sort: 1, Meta: system.Meta{Title: "仪表盘", Icon: "odometer"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 2}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "about", Name: "about", Component: "view/about/index.vue", Sort: 7, Meta: system.Meta{Title: "关于我们", Icon: "info-filled"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 3}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "admin", Name: "superAdmin", Component: "view/superAdmin/index.vue", Sort: 3, Meta: system.Meta{Title: "超级管理员", Icon: "user"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 4}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "authority", Name: "authority", Component: "view/superAdmin/authority/authority.vue", Sort: 1, Meta: system.Meta{Title: "角色管理", Icon: "avatar"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 5}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "menu", Name: "menu", Component: "view/superAdmin/menu/menu.vue", Sort: 2, Meta: system.Meta{Title: "菜单管理", Icon: "tickets", KeepAlive: true}},
		{GVA_MODEL: global.GVA_MODEL{ID: 6}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "api", Name: "api", Component: "view/superAdmin/api/api.vue", Sort: 3, Meta: system.Meta{Title: "api管理", Icon: "platform", KeepAlive: true}},
		{GVA_MODEL: global.GVA_MODEL{ID: 7}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "user", Name: "user", Component: "view/superAdmin/user/user.vue", Sort: 4, Meta: system.Meta{Title: "用户管理", Icon: "coordinate"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 8}, MenuLevel: 0, Hidden: true, ParentId: "0", Path: "person", Name: "person", Component: "view/person/person.vue", Sort: 4, Meta: system.Meta{Title: "个人信息", Icon: "message"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 9}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "example", Name: "example", Component: "view/example/index.vue", Sort: 6, Meta: system.Meta{Title: "示例文件", Icon: "management"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 10}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "excel", Name: "excel", Component: "view/example/excel/excel.vue", Sort: 4, Meta: system.Meta{Title: "excel导入导出", Icon: "takeaway-box"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 11}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "upload", Name: "upload", Component: "view/example/upload/upload.vue", Sort: 5, Meta: system.Meta{Title: "媒体库（上传下载）", Icon: "upload"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 12}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "breakpoint", Name: "breakpoint", Component: "view/example/breakpoint/breakpoint.vue", Sort: 6, Meta: system.Meta{Title: "断点续传", Icon: "upload-filled"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 13}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "customer", Name: "customer", Component: "view/example/customer/customer.vue", Sort: 7, Meta: system.Meta{Title: "客户列表（资源示例）", Icon: "avatar"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 14}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "systemTools", Name: "systemTools", Component: "view/systemTools/index.vue", Sort: 5, Meta: system.Meta{Title: "系统工具", Icon: "tools"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 15}, MenuLevel: 0, Hidden: false, ParentId: "14", Path: "autoCode", Name: "autoCode", Component: "view/systemTools/autoCode/index.vue", Sort: 1, Meta: system.Meta{Title: "代码生成器", Icon: "cpu", KeepAlive: true}},
		{GVA_MODEL: global.GVA_MODEL{ID: 16}, MenuLevel: 0, Hidden: false, ParentId: "14", Path: "formCreate", Name: "formCreate", Component: "view/systemTools/formCreate/index.vue", Sort: 2, Meta: system.Meta{Title: "表单生成器", Icon: "magic-stick", KeepAlive: true}},
		{GVA_MODEL: global.GVA_MODEL{ID: 17}, MenuLevel: 0, Hidden: false, ParentId: "14", Path: "system", Name: "system", Component: "view/systemTools/system/system.vue", Sort: 3, Meta: system.Meta{Title: "系统配置", Icon: "operation"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 18}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "dictionary", Name: "dictionary", Component: "view/superAdmin/dictionary/sysDictionary.vue", Sort: 5, Meta: system.Meta{Title: "字典管理", Icon: "notebook"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 19}, MenuLevel: 0, Hidden: true, ParentId: "3", Path: "dictionaryDetail/:id", Name: "dictionaryDetail", Component: "view/superAdmin/dictionary/sysDictionaryDetail.vue", Sort: 1, Meta: system.Meta{Title: "字典详情", Icon: "order"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 20}, MenuLevel: 0, Hidden: false, ParentId: "3", Path: "operation", Name: "operation", Component: "view/superAdmin/operation/sysOperationRecord.vue", Sort: 6, Meta: system.Meta{Title: "操作历史", Icon: "pie-chart"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 21}, MenuLevel: 0, Hidden: false, ParentId: "9", Path: "simpleUploader", Name: "simpleUploader", Component: "view/example/simpleUploader/simpleUploader", Sort: 6, Meta: system.Meta{Title: "断点续传（插件版）", Icon: "upload"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 22}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "https://www.gin-vue-admin.com", Name: "https://www.gin-vue-admin.com", Component: "/", Sort: 0, Meta: system.Meta{Title: "官方网站", Icon: "home-filled"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 23}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "state", Name: "state", Component: "view/system/state.vue", Sort: 6, Meta: system.Meta{Title: "服务器状态", Icon: "cloudy"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 24}, MenuLevel: 0, Hidden: false, ParentId: "14", Path: "autoCodeAdmin", Name: "autoCodeAdmin", Component: "view/systemTools/autoCodeAdmin/index.vue", Sort: 1, Meta: system.Meta{Title: "自动化代码管理", Icon: "magic-stick"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 25}, MenuLevel: 0, Hidden: true, ParentId: "14", Path: "autoCodeEdit/:id", Name: "autoCodeEdit", Component: "view/systemTools/autoCode/index.vue", Sort: 0, Meta: system.Meta{Title: "自动化代码（复用）", Icon: "magic-stick"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 26}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "TidbOrMysql", Name: "TidbOrMysql", Component: "view/tidbOrMysql/index.vue", Sort: 3, Meta: system.Meta{Title: "TiDB(Mysql)", Icon: "coin"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 27}, MenuLevel: 0, Hidden: false, ParentId: "26", Path: "ClusterManager", Name: "ClusterManager", Component: "view/tidbOrMysql/cluster/cluster.vue", Sort: 6, Meta: system.Meta{Title: "集群管理", Icon: "aim"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 28}, MenuLevel: 0, Hidden: false, ParentId: "26", Path: "SubmitSql", Name: "SubmitSql", Component: "view/tidbOrMysql/submit/submit.vue", Sort: 1, Meta: system.Meta{Title: "提交SQL", Icon: "chat-line-square"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 29}, MenuLevel: 0, Hidden: false, ParentId: "26", Path: "Review", Name: "Review", Component: "view/tidbOrMysql/review/review.vue", Sort: 2, Meta: system.Meta{Title: "审核与执行", Icon: "operation"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 30}, MenuLevel: 0, Hidden: false, ParentId: "26", Path: "Rule", Name: "Rule", Component: "view/tidbOrMysql/rule/rule.vue", Sort: 3, Meta: system.Meta{Title: "审核规则", Icon: "question-filled"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 31}, MenuLevel: 0, Hidden: true, ParentId: "26", Path: "exec/:id", Name: "exec", Component: "view/tidbOrMysql/exec/exec.vue", Sort: 4, Meta: system.Meta{Title: "任务执行", Icon: "pointer"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 32}, MenuLevel: 0, Hidden: false, ParentId: "26", Path: "TaskHistory", Name: "TaskHistory", Component: "view/tidbOrMysql/history/history.vue", Sort: 5, Meta: system.Meta{Title: "历史任务", Icon: "message-box"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 33}, MenuLevel: 0, Hidden: false, ParentId: "26", Path: "Read", Name: "Read", Component: "view/tidbOrMysql/read.vue", Sort: 0, Meta: system.Meta{Title: "数据读取", Icon: "search"}},

		{GVA_MODEL: global.GVA_MODEL{ID: 34}, MenuLevel: 0, Hidden: false, ParentId: "0", Path: "auth", Name: "auth", Component: "view/auth/apply.vue", Sort: 2, Meta: system.Meta{Title: "权限", Icon: "view"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 35}, MenuLevel: 0, Hidden: false, ParentId: "34", Path: "apply", Name: "apply", Component: "view/auth/apply.vue", Sort: 1, Meta: system.Meta{Title: "申请", Icon: "zoom-in"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 36}, MenuLevel: 0, Hidden: false, ParentId: "34", Path: "approval", Name: "approval", Component: "view/auth/approval.vue", Sort: 2, Meta: system.Meta{Title: "审批", Icon: "ship"}},
		{GVA_MODEL: global.GVA_MODEL{ID: 37}, MenuLevel: 0, Hidden: false, ParentId: "34", Path: "auths", Name: "auths", Component: "view/auth/auths.vue", Sort: 3, Meta: system.Meta{Title: "权限", Icon: "view"}},
	}
	if err := global.GVA_DB.Create(&entities).Error; err != nil { // 创建 model.User 初始化数据
		return errors.Wrap(err, m.TableName()+"表数据初始化失败!")
	}
	return nil
}

func (m *menu) CheckDataExist() bool {
	if errors.Is(global.GVA_DB.Where("path = ?", "autoCodeEdit/:id").First(&system.SysBaseMenu{}).Error, gorm.ErrRecordNotFound) { // 判断是否存在数据
		return false
	}
	return true
}
