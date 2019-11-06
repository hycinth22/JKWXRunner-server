package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func randRangeInt63(min, max int64) int64 {
	return min + rand.Int63()%(max-min+1)
}

func randSleepDuration(min time.Duration, max time.Duration) time.Duration {
	randNum := randRangeInt63(min.Nanoseconds(), max.Nanoseconds())
	return time.Duration(randNum) * time.Nanosecond
}

// can't guarantee must be not timeout
func sleepPartOfTotalTime(totalCount int64, totalTime time.Duration) time.Duration {
	totalTime = time.Duration(0.8 * float64(totalTime)) // 20% for delay & other
	single := totalTime.Nanoseconds() / totalCount
	d := randSleepDuration(time.Duration(0.8*float64(single)), time.Duration(1.2*float64(single)))
	log.Println("Sleep ", d.String())
	time.Sleep(d)
	return d
}

func sleepUtil(t time.Time) {
	time.Sleep(time.Until(t))
}

// %v the value in a default format, adds field names
func dumpStructValue(data interface{}) string {
	return fmt.Sprintf("%+v", data)
}

// %#v	a Go-syntax representation of the value
func dumpStruct(data interface{}) string {
	return fmt.Sprintf("%#v", data)
}
