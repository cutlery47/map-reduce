package internal

import "unicode"

func Strip(str string) string {
	res := ""

	for _, el := range str {
		if unicode.IsDigit(el) || unicode.IsLetter(el) || el == ' ' {
			res += string(el)
		}
	}

	return res
}

func Sum(vals []int) int {
	res := 0

	for _, el := range vals {
		res += el
	}

	return res
}
