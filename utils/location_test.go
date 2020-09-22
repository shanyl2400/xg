package utils

import (
	"testing"
	"xg/conf"
)

func TestGetAddressLocation(t *testing.T) {
	conf.Set(&conf.Config{
		AMapKey:   "5ffb23f18540eea567aa563445bff6ed",
	})
	res, err := GetAddressLocation("上海徐汇区")
	t.Log(err)
	t.Log(res)
}