package stats

func IsPowerOfTwo(n uint8) bool {
	return (n & (n - 1)) == 0
}
