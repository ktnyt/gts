package gts

// Unpack the integer pair to its elements.
func Unpack(p [2]int) (int, int) {
	return p[0], p[1]
}

const intSize = 32 << (^uint(0) >> 63)

// Abs returns the absolute value of the given integer.
func Abs(x int) int {
	y := x >> (intSize - 1)
	return (x ^ y) - y
}

// Compare the two integers and return the result.
func Compare(i, j int) int {
	switch {
	case i < j:
		return -1
	case j < i:
		return 1
	default:
		return 0
	}
}

// Min returns the smaller integer.
func Min(i, j int) int {
	if i < j {
		return i
	}
	return j
}

// Max returns the bigger integer.
func Max(i, j int) int {
	if j < i {
		return i
	}
	return j
}
