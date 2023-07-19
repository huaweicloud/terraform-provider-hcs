package utils

import (
	"crypto/rand"
	"math/big"
)

var reader = rand.Reader

func Max32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func Min32(x, y int32) int32 {
	if x > y {
		return y
	}
	return x
}

func RandInt32(min, max int32) int32 {
	if min >= max || min == 0 || max == 0 {
		return max
	}

	b, err := rand.Int(reader, big.NewInt(int64(max)))
	if err != nil {
		return max
	}
	return int32(b.Int64()) + min
}

func Pow32(x, y int32) int32 {
	var ans int32 = 1
	for y != 0 {
		ans *= x
		y--
	}
	return ans
}
