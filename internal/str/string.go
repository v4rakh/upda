package str

import (
	"math/rand"
	"regexp"
	"strings"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
)

var (
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// FindInSlice finds value in a slice
func FindInSlice(slice []string, val string) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

// ValuesString concatenate all values of a map split by comma
func ValuesString(m map[string]string) string {
	values := make([]string, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return strings.Join(values, ", ")
}

// ToSnakeCase converts string to snake case
func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	if n <= 0 {
		return ""
	}
	b := make([]byte, n)
	for i := 0; i < n; {
		if idx := int(rand.Int63() & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i++
		}
	}
	return string(b)
}

// ExtractBetween extracts all occurrences of a string within left and right delimiter, first inner array item is with delimiters, second one without
func ExtractBetween(str string, leftDelimiter string, rightDelimiter string) [][]string {
	if str == "" || leftDelimiter == "" || rightDelimiter == "" {
		return make([][]string, 0)
	}

	rx := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(leftDelimiter) + `(.*?)` + regexp.QuoteMeta(rightDelimiter))
	return rx.FindAllStringSubmatch(str, -1)
}
