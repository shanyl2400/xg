package service

import (
	"testing"
	"time"
)

func TestTimeDuration(t *testing.T) {
	t1 := time.Date(2020, 11, 20, 12, 15, 12, 10, time.UTC)
	t2 := time.Date(2020, 12, 20, 12, 15, 12, 10, time.UTC)
	dur := t1.Sub(t2)
	t.Log(dur > 0)
	t.Log(dur.Hours() / 24)

	t1 = t1.AddDate(0, 0, 20)
	t.Log(t1)
}

func TestForLoop(t *testing.T) {
	for i := 100; i <= 100; i++ {
		t.Log(i)
	}
}