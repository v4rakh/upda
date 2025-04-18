//go:build !integration

package str

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindInSlice(t *testing.T) {
	a := assert.New(t)
	a.True(FindInSlice([]string{""}, ""))
	a.True(FindInSlice([]string{"test", "abc"}, "test"))
	a.False(FindInSlice([]string{"abc"}, "test"))
	a.False(FindInSlice([]string{""}, "test"))
}

func TestToSnakeCase(t *testing.T) {
	a := assert.New(t)
	a.Equal("test_this", ToSnakeCase("TestThis"))
	a.Equal("this", ToSnakeCase("This"))
	a.Equal("", ToSnakeCase(""))
	a.NotEqual("test_this", ToSnakeCase("Test This"))
}

func TestRandomString(t *testing.T) {
	a := assert.New(t)
	a.Equal("", RandomString(-1))
	a.Equal(4, len(RandomString(4)))
}

func TestExtractValuesFromString(t *testing.T) {
	a := assert.New(t)
	a.Equal("", ValuesString(nil))
	a.Equal("val1", ValuesString(map[string]string{"key1": "val1"}))
	valuesString := ValuesString(map[string]string{"key1": "val1", "key2": "val2"})
	a.Contains(valuesString, "val1")
	a.Contains(valuesString, "val2")
}

func TestExtractBetweenEmpty(t *testing.T) {
	a := assert.New(t)
	a.Equal(0, len(ExtractBetween("", "", "")))
	a.Equal(0, len(ExtractBetween("test", "", "")))
	a.Equal(0, len(ExtractBetween("test", "test", "")))
}

func TestExtractBetweenVars(t *testing.T) {
	a := assert.New(t)

	left := "<VAR>"
	right := "</VAR>"

	//str := "A new update arrived on <VAR>HOST</VAR> for <VAR>APPLICATION</VAR>. Its version is <VAR>VERSION</VAR>"
	str := "<VAR>MY_VAR</VAR> A new update arrived on <VA>HOST</VAR> for for <VAR>APPLICATION</VAR> (<VAR>MY_SECOND_VAR</VAR>) ... some random other tag <SECRET>MY_SECRET</SECRET>"
	matches := ExtractBetween(str, left, right)

	a.Equal(true, len(matches) > 0)
	a.Equal("MY_VAR", matches[0][1])
	a.Equal("APPLICATION", matches[1][1])
	a.Equal("MY_SECOND_VAR", matches[2][1])
}
