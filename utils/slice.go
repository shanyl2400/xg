package utils
func SliceDeduplication(s []string) []string {
	temp := make(map[string]bool)
	for i := range s {
		temp[s[i]] = true
	}

	result := make([]string, 0, len(temp))
	for k, v := range temp {
		if v {
			result = append(result, k)
		}
	}

	return result
}

func SliceDeduplicationInt(s []int) []int {
	temp := make(map[int]bool)
	for i := range s {
		temp[s[i]] = true
	}

	result := make([]int, 0, len(temp))
	for k, v := range temp {
		if v {
			result = append(result, k)
		}
	}

	return result
}
