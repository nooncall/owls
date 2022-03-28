package utils

import (
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
