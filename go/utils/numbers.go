package utils

import "strconv"

func SubTwoStrNum(numA, numB string) (int64, error) {
	a, err := strconv.ParseInt(numA, 10, 64)
	if err != nil {
		return 0, err
	}

	b, err := strconv.ParseInt(numB, 10, 64)
	if err != nil {
		return 0, err
	}

	return a - b, nil
}

func Int64ArrayDeduplication(data []int64) (result []int64) {
	mapData := make(map[int64]struct{})
	for _, v := range data {
		mapData[v] = struct{}{}
	}

	for k, _ := range mapData {
		result = append(result, k)
	}
	return
}
