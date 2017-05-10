package jsonpp

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadValue(t *testing.T) {
	json := `{"key1":"value1","key2":"value2"}`

	value := new(string)
	err := Read(strings.NewReader(json), []string{"key1"}, value)

	assert.Equal(t, "value1", *value)
	assert.NoError(t, err)
}

func TestReadArray(t *testing.T) {
	json := `{"key1":"value1","key2":{"key3":[1, 2, 3]},"key4":"value4"}`

	value := new([]int)
	err := Read(strings.NewReader(json), []string{"key2", "key3"}, value)

	assert.Equal(t, []int{1, 2, 3}, *value)
	assert.NoError(t, err)
}

func TestIdentity(t *testing.T) {
	json := `{"key":"value1","key2":"value2"}`

	modified, err := Rewrite(strings.NewReader(json), nil, nil)

	assert.Equal(t, json, modified)
	assert.NoError(t, err)
}

func TestRewriteValue(t *testing.T) {
	json := `{"key1":"value1","key2":"value2"}`
	expected := `{"key1":"newValue","key2":"value2"}`

	modified, err := Rewrite(strings.NewReader(json), []string{"key1"}, "newValue")

	assert.Equal(t, expected, modified)
	assert.NoError(t, err)
}

func TestRewriteValueDeep(t *testing.T) {
	json := `{"key1":"value1","key2":{"key3":"value3"},"key4":"value4"}`
	expected := `{"key1":"value1","key2":{"key3":"newValue"},"key4":"value4"}`

	modified, err := Rewrite(strings.NewReader(json), []string{"key2", "key3"}, "newValue")

	assert.Equal(t, expected, modified)
	assert.NoError(t, err)
}

func TestRewriteIntValue(t *testing.T) {
	json := `{"key1":"value1","key2":{"key3":1},"key4":"value4"}`
	expected := `{"key1":"value1","key2":{"key3":42},"key4":"value4"}`

	modified, err := Rewrite(strings.NewReader(json), []string{"key2", "key3"}, 42)

	assert.Equal(t, expected, modified)
	assert.NoError(t, err)
}

func TestRewriteArrayValue(t *testing.T) {
	json := `{"key1":"value1","key2":{"key3":[1,42,3]},"key4":"value4"}`
	expected := `{"key1":"value1","key2":{"key3":[1,2,3]},"key4":"value4"}`

	modified, err := Rewrite(strings.NewReader(json), []string{"key2", "key3", "[1]"}, 2)

	assert.Equal(t, expected, modified)
	assert.NoError(t, err)
}

func TestRewriteArray(t *testing.T) {
	json := `{"key1":"value1","key2":{"key3":[1, 42, 3]},"key4":"value4"}`
	expected := `{"key1":"value1","key2":{"key3":["42"]},"key4":"value4"}`

	modified, err := Rewrite(strings.NewReader(json), []string{"key2", "key3"}, []string{"42"})

	assert.Equal(t, expected, modified)
	assert.NoError(t, err)
}
