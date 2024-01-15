package utils

import "fmt"

// バージョン情報を文字列に変換する
func ParseVersion(v uint32) string {
	// メジャーバージョン1byte, マイナーバージョン1byte, パッチバージョン2byte
	// 上位8ビット
	major := v >> 24
	// 次の8ビット
	minor := (v >> 16) & 0xff
	// 下位16ビット
	patch := v & 0xffff
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
