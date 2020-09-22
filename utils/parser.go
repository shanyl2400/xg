package utils

import (
	"strconv"
	"strings"
	"xg/log"
)

func ParseInt(str string) int {
	x, err := strconv.Atoi(str)
	if err != nil{
		log.Warning.Println("Parse in failed, error:", err)
		return 0
	}
	return x
}

func ParseInts(str string) []int {
	strList := strings.Split(str, ",")
	ret := make([]int, 0)
	for i := range strList {
		id, err := strconv.Atoi(strList[i])
		if err == nil {
			ret = append(ret, id)
		}
	}
	if len(ret) < 1 {
		return nil
	}
	return ret
}


func ParseFloats(str string) []float64 {
	strList := strings.Split(str, ",")
	ret := make([]float64, 0)
	for i := range strList {
		id, err := strconv.ParseFloat(strList[i], 64)
		if err == nil {
			ret = append(ret, id)
		}
	}
	if len(ret) < 1 {
		return nil
	}
	return ret
}