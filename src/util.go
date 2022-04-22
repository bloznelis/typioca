package main

import "math"

func averageStringLen(strings []string) int {
	var totalLen int = 0
	var cnt int = 0

	for _, str := range strings {
		currentLen := len(dropAnsiCodes(str))
		totalLen += currentLen
		cnt += 1
	}

	return totalLen / cnt
}

func floor(value int) int32 {
	return int32(math.Max(0, float64(value)))
}

func dropLastString(strings []string) []string {
	le := len(strings)
	if le != 0 {
		return strings[:le-1]
	} else {
		return strings
	}
}

func dropLastRune(runes []rune) []rune {
	le := len(runes)
	if le != 0 {
		return runes[:le-1]
	} else {
		return runes
	}
}

func toKeysSlice(mp map[int]bool) []int {
	acc := []int{}
	for key := range mp {
		acc = append(acc, key)
	}
	return acc
}

func reverse(runes []rune) []rune {
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return runes
}
