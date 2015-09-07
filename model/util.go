package model

import (
	"unicode"
	"unicode/utf8"
)

func toSnake(src string) string {
	buf := make([]byte, 0)
	runeBuf := make([]byte, 1)
	for i := 0; i < len(src); i++ {
		val := rune(src[i])
		if unicode.IsUpper(val) {
			if i > 0 {
				buf = append(buf, '_')
			}
			utf8.EncodeRune(runeBuf, unicode.ToLower(val))
			buf = append(buf, runeBuf...)
		} else {
			utf8.EncodeRune(runeBuf, val)
			buf = append(buf, runeBuf...)
		}
	}
	return string(buf)
}
