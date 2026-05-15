package utils

import (
	"encoding/binary"
	"strings"
	"unicode"
)

func ClearStr(str string) string {
	str = strings.ToLower(strings.TrimSpace(str))
	var b strings.Builder
	b.Grow(len(str))
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func ClearStrSpace(str string) string {
	str = strings.ToLower(strings.TrimSpace(str))
	var b strings.Builder
	b.Grow(len(str))
	var last rune = -1
	for _, r := range str {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			last = r
		} else if last != -1 && last != ' ' {
			b.WriteByte(' ')
			last = ' '
		}
	}
	return b.String()
}

func I2B(num int64) []byte {
	value := make([]byte, 8)
	binary.BigEndian.PutUint64(value, uint64(num))
	return value
}

func B2I(value []byte) int64 {
	return int64(binary.BigEndian.Uint64(value))
}
