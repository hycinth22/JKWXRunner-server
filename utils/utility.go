package utils

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var TimeZoneBeijing = time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))

func RandRangeInt63(min, max int64) int64 {
	return min + rand.Int63()%(max-min+1)
}

func RandSleepDuration(min time.Duration, max time.Duration) time.Duration {
	randNum := RandRangeInt63(min.Nanoseconds(), max.Nanoseconds())
	return time.Duration(randNum) * time.Nanosecond
}

// can't guarantee must be not timeout
func SleepPartOfTotalTime(totalCount int, totalTime time.Duration) {
	totalTime = time.Duration(0.8 * float64(totalTime)) // 20% for delay & other
	single := totalTime.Nanoseconds() / int64(totalCount)

	var d time.Duration
	if time.Duration(single) > 5*time.Minute {
		d = RandSleepDuration(15*time.Second, 5*time.Minute)
	} else {
		d = RandSleepDuration(time.Duration(0.8*float64(single)), time.Duration(1.2*float64(single)))
	}

	log.Println("Sleep ", d.String())
	time.Sleep(d)
}

func SleepUtil(t time.Time) {
	time.Sleep(time.Until(t))
}

// %v the value in a default format, adds field names
func DumpStructValue(data interface{}) string {
	return fmt.Sprintf("%+v", data)
}

// %#v	a Go-syntax representation of the value
func DumpStruct(data interface{}) string {
	return fmt.Sprintf("%#v", data)
}
