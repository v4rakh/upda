package str

import (
	"regexp"
	"strings"
)

// ValuesString concatenate all values of a map split by comma
func ValuesString(m map[string]string) string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return strings.Join(values, ", ")
}

// ExtractBetween extracts all occurrences of a string within left and right delimiter, first inner array item is with delimiters, second one without
func ExtractBetween(str string, leftDelimiter string, rightDelimiter string) [][]string {
	if str == "" || leftDelimiter == "" || rightDelimiter == "" {
		return make([][]string, 0)
	}

	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(leftDelimiter) + `(.*?)` + regexp.QuoteMeta(rightDelimiter))
	return rx.FindAllStringSubmatch(str, -1)
}

// AllContained determines if all a slice values are in the b slice
func AllContained(a, b []string) bool {
	lookup := make(map[string]struct{}, len(b))
	for _, v := range b {
		lookup[v] = struct{}{}
	}
	for _, item := range a {
		if _, ok := lookup[item]; !ok {
			return false
		}
	}
	return true
}
