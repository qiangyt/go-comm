package comm

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestStruct struct {
	Name  string
	Value int
}

func TestStructToMap_happy(t *testing.T) {
	a := require.New(t)

	s := TestStruct{
		Name:  "test",
		Value: 42,
	}

	result := StructToMap(s)
	a.NotNil(result)
	a.Equal("test", result["Name"])
	a.Equal(42, result["Value"])
}

func TestMap2Struct_happy(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"Name":  "test",
		"Value": 100,
	}

	var s TestStruct
	err := Map2Struct(m, &s)
	a.NoError(err)
	a.Equal("test", s.Name)
	a.Equal(100, s.Value)
}

func TestMap2StructP_happy(t *testing.T) {
	a := require.New(t)

	m := map[string]any{
		"Name":  "hello",
		"Value": 200,
	}

	var s TestStruct
	Map2StructP(m, &s)
	a.Equal("hello", s.Name)
	a.Equal(200, s.Value)
}

func TestStructToMap_empty(t *testing.T) {
	a := require.New(t)

	s := TestStruct{}

	result := StructToMap(s)
	a.NotNil(result)
	a.Equal("", result["Name"])
	a.Equal(0, result["Value"])
}

func TestMap2Struct_empty(t *testing.T) {
	a := require.New(t)

	m := map[string]any{}

	var s TestStruct
	err := Map2Struct(m, &s)
	a.NoError(err)
	a.Equal("", s.Name)
	a.Equal(0, s.Value)
}
