package masks

import "strings"

func All[T ~string](s T) T {
	return "[MASKED]"
}

func Left[T ~string](s T, visibleLetterCount int) T {
	length := len(s)
	if visibleLetterCount >= length {
		return s
	}
	return T(strings.Repeat("*", length-visibleLetterCount)) + s[length-visibleLetterCount:]
}
