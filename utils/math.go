package utils

import "regexp"

func IsNumeric(input string) bool {
	// 定义正则表达式，匹配整数或浮点数
	re := regexp.MustCompile(`^[+-]?(\d+(\.\d*)?|\.\d+)$`)
	return re.MatchString(input)
}
