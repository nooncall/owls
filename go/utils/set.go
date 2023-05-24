package utils

func MergeStringMaps(m1, m2 map[string]string) map[string]string {
	for k, v := range m1 {
		m2[k] = v
	}

	return m2
}
