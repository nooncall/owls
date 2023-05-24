package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func ErrsJoin(str string, err []error) string {
	if len(err) < 1 {
		return ""
	}
	result := ""
	for i, v := range err {
		if v == nil {
			continue
		}

		if i == 0 {
			result += v.Error()
			continue
		}
		result += v.Error() + str
	}
	return result
}

func StringifyByteDirectly(b []byte) string {
	var str []string
	for _, v := range b {
		str = append(str, fmt.Sprintf("%d", v))
	}
	return strings.Join(str, ",")
}

func ParseStringedByte(str string) []byte {
	splitStr := strings.Split(str, ",")
	var b []byte
	for _, v := range splitStr {
		intV, _ := strconv.Atoi(v)
		b = append(b, byte(intV))
	}
	return b
}

//删除首尾空格， 连续空格换成单个空格
func DelUselessSpace(str string) string {
	str = strings.TrimSpace(str)
	if strings.Contains(str, "  ") {
		return DelUselessSpace(strings.Replace(str, "  ", " ", -1))
	}
	if strings.Contains(str, "	") {
		return DelUselessSpace(strings.Replace(str, "	", " ", -1))
	}

	return str
}

func ErrsToString(err []error) string {
	var str string
	for _, v := range err {
		if v != nil {
			str += v.Error()
		}
	}
	return str
}

func ErrsToErr(err []error) error {
	if len(err) < 1 {
		return nil
	}

	var str string
	for _, v := range err {
		if v != nil {
			str += ";" + v.Error()
		}
	}
	return errors.New(str)
}

func StringArrayDeduplication(data []string) (result []string) {
	mapData := make(map[string]struct{})
	for _, v := range data {
		mapData[v] = struct{}{}
	}

	for k, _ := range mapData {
		result = append(result, k)
	}
	return
}
