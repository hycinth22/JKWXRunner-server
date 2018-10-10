package view

import (
	"strconv"
	"time"
)

func DistanceFormat(distance float64) string {
	return strconv.FormatFloat(distance, 'f', 3, 64)
}
func TimeFormat(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}
