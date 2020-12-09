package pgsql

import (
	"strconv"
	"strings"
)

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

// Unique ...
func Unique(input []int) []int {
	result := make([]int, 0, len(input))
	values := make(map[int]bool)

	for _, val := range input {
		if _, ok := values[val]; !ok {
			values[val] = true
			result = append(result, val)
		}
	}

	return result
}

// Contains ...
func Contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}
