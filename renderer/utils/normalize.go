package utils

// N〜Mの範囲にあるxを-1〜1に正規化する
func Normalize(x, n, m float32) (normalized float32) {
	normalized = (x - n) / (m - n)
	normalized = normalized*2 - 1
	return
}
