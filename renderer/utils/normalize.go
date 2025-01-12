package utils

// Normalize x, which is in the range from N to M, to the range from -1 to 1
func Normalize(x, n, m float32) (normalized float32) {
	normalized = (x - n) / (m - n)
	normalized = normalized*2 - 1
	return
}
