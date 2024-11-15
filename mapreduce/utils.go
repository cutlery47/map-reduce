package main

import "unicode"

func strip(str string) string {
	res := ""

	for _, el := range str {
		if unicode.IsDigit(el) || unicode.IsLetter(el) || el == ' ' {
			res += string(el)
		}
	}

	return res
}

func sum(vals []int) int {
	res := 0

	for _, el := range vals {
		res += el
	}

	return res
}
