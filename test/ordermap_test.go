package test

import (
	"testing"

	"github.com/fastgh/go-comm/v2"
	"github.com/stretchr/testify/require"
)

func Test_OrderedMap_happy(t *testing.T) {
	a := require.New(t)

	m := comm.NewOrderedMap("")
	a.Equal(0, m.Len())
	a.Len(m.Keys(), 0)
	a.Len(m.Values(), 0)

	m.Set("k1", "v1")
	m.Set("k2", "v2")

	a.Equal(2, m.Len())
	a.Len(m.Keys(), 2)
	a.Len(m.Values(), 2)

	a.Equal("v1", m.Get("k1"))
	a.True(m.Has("k1"))
	v1, b1 := m.Find("k1")
	a.True(b1)
	a.Equal("v1", v1)

	a.Equal("v2", m.Get("k2"))
	a.True(m.Has("k2"))
	v2, b2 := m.Find("k2")
	a.True(b2)
	a.Equal("v2", v2)

	a.Equal("", m.Get("k3"))
	a.False(m.Has("k3"))

	m.Delete("k2")
	a.Equal(1, m.Len())
	a.Equal("", m.Get("k2"))
	a.False(m.Has("k2"))

	entries := m.Entries()
	a.Len(entries, 1)
	a.Equal("k1", entries[0].Key)
	a.Equal("v1", entries[0].Value)
}

func Test_OrderedMap_json(t *testing.T) {
	a := require.New(t)

	m := comm.NewOrderedMap("")
	m.Set("k1", "v1")

	json, err := m.MarshalJSON()
	a.Nil(err)
	a.Equal(`{"k1"
:"v1"
}`, (string(json)))

	m2 := comm.NewOrderedMap("")
	err = m2.UnmarshalJSON(json)
	a.Nil(err)
	a.Equal(m, m2)
}
