package sql_util

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/nooncall/owls/go/core"
	"github.com/nooncall/owls/go/global"
	"github.com/nooncall/owls/go/utils/logger"
)

func initConfigLog() {
	global.GVA_VP = core.Viper() // 初始化Viper
	logger.InitLog(global.GVA_CONFIG.DBFilter.LogDir, "test.log", "debug")
}

func TestGetUpdateColumn(t *testing.T) {
	initConfigLog()

	type data struct {
		origin string
		target []string
		isErrr bool
	}
	datas := []data{
		{
			origin: "update test_table set name=zhangsan where age=19 ",
			target: []string{"name"},
			isErrr: false,
		},
		{
			origin: "update test_table set name=zhangsan",
			target: []string{"name"},
			isErrr: true,
		},
		{
			origin: "UPDATE test_table set name=zhangsan, age = 18 where age=19",
			target: []string{"name", "age"},
			isErrr: false,
		},
		{
			origin: "",
			target: nil,
			isErrr: true,
		},
	}

	for _, v := range datas {
		tar, err := GetSqlColumn(v.origin)
		if v.isErrr && err != nil {
			continue
		}
		if err != nil || !reflect.DeepEqual(tar, v.target) {
			t.Log(fmt.Sprintf("expert: %v, \n got : %v, err: %v", v.target, tar, err))
			t.FailNow()
		}
	}
}

func TestDeleteSpecifyCharAtHead(t *testing.T) {
	initConfigLog()

	type data struct {
		origin string
		target string
	}
	datas := []data{
		{
			origin: " update ",
			target: "update ",
		},
		{
			origin: "; ;  ;update",
			target: "update",
		},
		{
			origin: `
;
hello`,
			target: "hello",
		},
		{
			origin: "",
			target: "",
		},
	}

	for _, v := range datas {
		tar := deleteSpecifyCharAtHead(v.origin)
		if tar != v.target {
			t.Log(fmt.Sprintf("expert :%s, \n got :%s", v.target, tar))
			t.Fail()
		}
	}
}

func TestBuildDelRollBackSql(t *testing.T) {
	initConfigLog()

	tableName := "auth"

	column := []Column{
		{
			Field: "id",
		},
		{
			Field: "channels",
		},
		{
			Field: "permission",
		},
		{
			Field: "topic",
		},
		{
			Field: "backup_status",
		},
		{
			Field: "backup_id",
		},
	}
	data := [][]string{
		{"1", "hello", "2", "hello", "4", "5"},
		{"4", "hi", "3", "boy", "5", "6"},
	}
	target := `INSERT INTO auth (id, channels, permission, topic, backup_status, backup_id) VALUES ('1', 'hello', '2', 'hello', '4', '5'), ('4', 'hi', '3', 'boy', '5', '6');`

	sql, err := buildDelRollBackSql(column, data, tableName)
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}
	if sql != target {
		t.Fail()
	}
}

func TestReplaceSpecifyChar(t *testing.T) {
	initConfigLog()

	type data struct {
		origin string
		target string
	}
	datas := []data{
		{
			origin: " update  ",
			target: " update ",
		},
		{
			origin: "; ;  ;update  ",
			target: "; ; ;update ",
		},
		{
			origin: `
;
hello`,
			target: " ; hello",
		},
		{
			origin: "",
			target: "",
		},
		{
			origin: "delete from  picture",
			target: "delete from picture",
		},
	}

	for _, v := range datas {
		tar := replaceSpecifyChar(v.origin)
		if tar != v.target {
			t.Log(fmt.Sprintf("expert :%s, \n got :%s", v.target, tar))
			t.Fail()
		}
	}
}

func TestHandleKeyWordForCondition(t *testing.T) {
	initConfigLog()

	type data struct {
		origin string
		target string
	}
	datas := []data{
		{
			origin: "FirstName='Bush' AND LastName='hello' and age >17 and name like '%hello%' and high between 175 and 180",
			target: "FirstName='Bush' AND LastName='hello' and age >17 and name like '%hello%' and high between 175 and 180",
		},
		{
			origin: "",
			target: "",
		},
		{
			origin: "uid = 40792760 AND courseid = 295748327620620 AND `index` >1 and name not in ('hi','hei');",
			target: "uid = 40792760 AND courseid = 295748327620620 AND `index` >1 and name not in ('hi','hei');",
		},
		{
			origin: "uid = 40792760 AND courseid = 295748327620620 AND index >1 and user not in ('hi','hei');",
			target: "uid = 40792760 AND courseid = 295748327620620 AND `index` >1 and user not in ('hi','hei');",
		},
	}

	for _, v := range datas {
		resp := HandelKeyWorldForCondition(v.origin)
		if !reflect.DeepEqual(resp, v.target) {
			t.Log(fmt.Sprintf("expert: %v, \n got : %v", v.target, resp))
			t.FailNow()
		}
	}
}

func TestGetCondition(t *testing.T) {
	initConfigLog()

	type data struct {
		origin string
		target []string
	}
	datas := []data{
		{
			origin: "FirstName='Bush' AND LastName='hello' and age >17 and name like '%hello%' and high between 175 and 180",
			target: []string{"firstname", "lastname", "age", "name", "high"},
		},
		{
			origin: "",
			target: nil,
		},
		{
			origin: "uid = 40792760 AND courseid = 295748327620620 AND `index` >1 and name not in ('hi','hei');",
			target: []string{"uid", "courseid", "`index`", "name"},
		},
	}

	for _, v := range datas {
		tar := GetCondition(v.origin)
		if !reflect.DeepEqual(tar, v.target) {
			t.Log(fmt.Sprintf("expert: %v, \n got : %v", tar, v.target))
			t.FailNow()
		}
	}
}

func TestIsSubKey(t *testing.T) {
	initConfigLog()

	type Key struct {
		KeyS   string
		KeyL   string
		Result bool
	}

	keys := []Key{
		{
			"a",
			"ab+cd",
			false,
		},
		{
			"ab",
			"ab+cd",
			true,
		},
	}
	for _, v := range keys {
		if resp := IsSubKey(v.KeyL, v.KeyS); resp != v.Result {
			t.FailNow()
		}
		if resp := IsSubKey(v.KeyS, v.KeyL); resp != v.Result {
			t.FailNow()
		}
	}
}
