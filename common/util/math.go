package util

import "math"

// MathRound 四舍五入一个浮点数，返回整数
func MathRound(x float64) int {
	mid := math.Floor(x) + 0.5
	if x >= mid {
		return int(math.Ceil(x))
	}
	return int(math.Floor(x))
}
