package rainbow

import "strings"

// MessageCharsに含まれる文字だけを使った文字列のうち、
// strの次に辞書順最小の文字列を返す
func nextPermutation(str string) string {
	revStr := reverseString(str)
	nextRevStr := next(revStr, 0)
	nextStr := reverseString(nextRevStr)
	return nextStr
}

func reverseString(str string) string {
	revStr := make([]rune, len(str))
	for i, v := range str {
		revStr[len(revStr)+^i] = v
	}
	return string(revStr)
}

func next(str string, index int) string {
	if index >= MessageLength {
		return strings.Repeat(MessageChars[0:1], MessageLength)
	}

	i := strings.Index(MessageChars, str[index:index+1])
	if len(MessageChars)-1 > i {
		// 現在の桁を1増やす
		return str[:index] + MessageChars[i+1:i+2] + str[index+1:]
	}

	// 現在の桁を0に戻す
	str = str[:index] + MessageChars[:1] + str[index+1:]
	return next(str, index+1)
}
