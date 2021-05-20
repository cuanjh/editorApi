package utils

import "math/rand"

// 函　数：生成随机数
// 概　要：
// 参　数：
// min: 最小值
// max: 最大值
// 返回值：
//      int64: 生成的随机数
func RandInt64(min, max int64) int64 {
	if min >= max || min == 0 || max == 0 {
		return max
	}
	return rand.Int63n(max-min) + min
}

func GetBaseRand(num int) int64 {
	var rand_num int64
	switch num {
	case 1:
		rand_num = RandInt64(700, 1000)
	case 2:
		rand_num = RandInt64(500, 700)
	default:
		rand_num = RandInt64(300, 500)
	}
	return rand_num
}
