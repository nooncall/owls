package utils

import (
	"fmt"
	"testing"
)

func TestDelUselessSpace(t *testing.T) {
	type Data struct {
		Str    string
		Target string
	}
	data := []Data{
		{
			Str:    "  hgetall      hello",
			Target: "hgetall hello",
		},
		{
			Str: "	hgetall hello",
			Target: "hgetall hello",
		},
		{
			Str: "	he	  llo	",
			Target: "he llo",
		},
	}

	for _, v := range data {
		resp := DelUselessSpace(v.Str)
		if v.Target != resp {
			t.Log(fmt.Sprintf("test DelUselessSpace failed, expect: %s, got:%s", v.Target, resp))
			t.FailNow()
		}
	}
}

func TestSubTwoStrNum(t *testing.T) {
	type Data struct {
		a      string
		b      string
		target int64
	}
	data := []Data{
		{
			a:      "10",
			b:      "8",
			target: 2,
		},
		{
			a:      "7",
			b:      "8",
			target: -1,
		},
	}

	for _, v := range data {
		resp, err := SubTwoStrNum(v.a, v.b)
		if err != nil {
			t.Log("test sub str num err: ", err.Error())
			t.FailNow()
		}

		if v.target != resp {
			t.Log(fmt.Sprintf("test DelUselessSpace failed, expect: %d, got:%d", v.target, resp))
			t.FailNow()
		}
	}
}
