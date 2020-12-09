package pgsql

import (
	"strconv"
	"strings"
)

// StrWithCommaToArrStr ...
func StrWithCommaToArrStr(values []string) []string {
	var tag string
	for _, v := range values {
		tag = v
	}
	return strings.Split(tag, ",")
}

// Find ...
func Find(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// ReplaceSQL replaces the instance occurrence of any string pattern with an increasing $n based sequence
func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

// Eq ...
func Eq(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}

	for k, v := range a {
		if w, ok := b[k]; !ok || v != w {
			return false
		}
	}

	return true
}
