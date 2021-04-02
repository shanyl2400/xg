package utils

import "strconv"

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

func StringsToInts(s []string) ([]int, error) {
	ret := make([]int, len(s))
	var err error
	for i := range s {
		ret[i], err = strconv.Atoi(s[i])
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func IntsToStrings(d []int) []string {
	ret := make([]string, len(d))
	for i := range d {
		ret[i] = strconv.Itoa(d[i])
	}
	return ret
}
