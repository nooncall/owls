package checker

import (
	"fmt"
	"testing"

	"github.com/qingfeng777/owls/server/service/tidb_or_mysql/service"
)

func TestIndexMatchCondition(t *testing.T) {
	service.InitConfigLog()

	keys := &[]KeysInfo{
		{
			KeyName:    "PRIMARY",
			SeqInIndex: 1,
			ColumnName: "id",
		},
		{
			KeyName:    "idx_record_video",
			SeqInIndex: 1,
			ColumnName: "roomid",
		},
	}
	resp := indexMatchConditionOrdinal(keys, []string{"id"})
	t.Log("resp is ", resp)
}

func TestOperateDisableIndex(t *testing.T) {
	service.InitConfigLog()

	type data struct {
		origin string
		target bool
	}
	datas := []data{
		{
			origin: "FirstName='Bush' AND LastName='hello' and age >17 and name like '%hello%' and high between 175 and 180",
			target: true,
		},
		{
			origin: "FirstName='Bush' AND LastName='hello' and age >17 and name like 'hello%' and high between 175 and 180",
			target: true,
		},
		{
			origin: "FirstName='Bush' AND LastName='hello' and age >17 and high between 175 and 180 and name like 'hello%' ",
			target: true,
		},
		{
			origin: "FirstName='Bush' AND LastName='hello' and age >17 and high = 180 ",
			target: true,
		},
		{
			origin: "FirstName='Bush' AND LastName='hello' and age =17 and high between 175 and 180 ",
			target: false,
		},
		{
			origin: "",
			target: false,
		},
	}

	for i, v := range datas {
		tar := operateDisableIndex(v.origin)
		if tar != v.target {
			t.Log(fmt.Sprintf("expert: %v, \n got : %v, sql index :%d ", tar, v.target, i))
			t.FailNow()
		}
	}
}

func TestDmlSqlToCount(t *testing.T) {
	service.InitConfigLog()

	type data struct {
		origin string
		target string
		isErr  bool
		opType opType
	}
	datas := []data{
		{
			"update test_table set name=zhangsan where age=19 ",
			"select count(*) from test_table where age=19",
			false,
			Update,
		},
		{
			"",
			"",
			true,
			Update,
		},
		{
			"delete from test_table  where age in (select age from test_table where naem=zhangsan) ",
			"select count(*) from test_table where age in (select age from test_table where naem=zhangsan)",
			false,
			Delete,
		},
	}
	for _, v := range datas {
		target, err := dmlSqlToCount(v.origin)
		if err != nil {
			if v.isErr {
				continue
			}
			t.Log("err not expected: ", err.Error())
			t.FailNow()
		}
		if v.target != target {
			t.Log(fmt.Sprintf("translate err , expect : %s, got : %s", v.target, target))
			t.FailNow()
		}
	}
}
