package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randRangeInt63(min, max int64) int64 {
	return min + rand.Int63()%(max-min+1)
}

func randSleep(min time.Duration, max time.Duration) {
	randNum := randRangeInt63(min.Nanoseconds(), max.Nanoseconds())
	time.Sleep(time.Duration(randNum) * time.Nanosecond)
}
