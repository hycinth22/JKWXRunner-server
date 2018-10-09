package view

import (
	"strconv"
)

func DistanceFormat(distance float64) string {
	return strconv.FormatFloat(distance, 'f', 3, 64)
}
