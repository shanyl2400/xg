package tests

import (
	"math/rand"
	"time"
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	numCharset = "0123456789"
)
var(
	numHeaders = []string{"138", "175", "189"}
	subjects = []string{
		"外语-英语",
		"外语-日语",
		"外语-韩语",
		"外语-俄罗斯语",
		"外语-葡萄牙语",
		"外语-柬埔寨语",
		"设计-Photoshop",
		"设计-Illustrator",
		"设计-Sketch",
		"设计-手绘",
		"设计-服装",
		"计算机-Java",
		"计算机-前端",
		"计算机-React",
		"计算机-Golang",
		"计算机-Python",
	}

	addresses = []string {
		"上海市闵行区",
		"上海市黄浦区",
		"上海市徐汇区",
		"上海市长宁区",
		"上海市静安区",
		"上海市普陀区",
		"上海市虹口区",
		"上海市杨浦区",
		"上海市宝山区",
		"上海市嘉定区",
		"上海市金山区",
		"上海市松江区",
		"上海市青浦区",
		"上海市奉贤区",
		"上海市崇明区",
		"上海市浦东新区",
	}

)

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandTelephone() string {
	index := RandIndex(len(numHeaders))
	numHead := numHeaders[index]
	return numHead + StringWithCharset(8, numCharset)
}

func RandEmail(length int) string {
	return StringWithCharset(length, charset) + "@163.com"
}

func RandString(length int) string {
	return StringWithCharset(length, charset)
}

func RandInt() int{
	return rand.Int() % 1000
}

func RandIndex(length int) int {
	return rand.Int() % length
}
func RandArray(arr []string) string {
	return arr[RandIndex(len(arr))]
}

func RandArrayList(arr []string, length int) []string {
	ret := make([]string, length)
	for i := 0; i < length; i ++ {
		ret[i] = RandArray(arr)
	}
	return ret
}