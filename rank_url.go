package mycelium

import (
	"math"
)

func RankLength(url string) float64 {
	maxLength := 100.0

	return math.Min(float64(len(url))/maxLength, 1.0)
}
