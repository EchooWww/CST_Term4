package util

import (
	"math/rand"
	"time"
)

func Work(min, max int) {
	if min >= max {
		return
	}
	n := rand.Intn(max - min)
	time.Sleep(time.Duration(n) * time.Millisecond)
}