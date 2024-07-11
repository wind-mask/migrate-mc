package api

import "strings"

// 将首字母转换为小写
func toLowerFirst(s string) string {
	if s == "" {
		return ""
	}
	r := []rune(s)
	r[0] = rune(strings.ToLower(string(r[0]))[0])
	return string(r)
}
