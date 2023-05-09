package util

import "strconv"

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func ParseLong(s string) int64 {
	num, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		num = 0
	}
	return num
}
