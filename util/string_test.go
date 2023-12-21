package util

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
