package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randBetween0to1() float64 {
	return rand.Float64()
}
