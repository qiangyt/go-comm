package qcoll

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_OrderedMap_happy(t *testing.T) {
	a := require.New(t)

	m := NewOrderedMap("")
	a.Equal(0, m.Len())
	a.Equal(0, m.Keys().Size())
	a.Len(m.Values(), 0)

	m.Put("k1", "v1")
	m.Put("k2", "v2")

	a.Equal(2, m.Len())
	a.Equal(2, m.Keys().Size())
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

	m := NewOrderedMap("")
	m.Put("k1", "v1")

	json, err := m.MarshalJSON()
	a.Nil(err)
	a.Equal(`{"k1"
:"v1"
}`, (string(json)))

	m2 := NewOrderedMap("")
	err = m2.UnmarshalJSON(json)
	a.Nil(err)
	a.Equal(m, m2)
}

func Test_OrderedMap_putAll(t *testing.T) {
	a := require.New(t)

	type Element struct {
		Name  string
		Value string
	}

	m := NewOrderedMap[*Element](nil)
	m.PutAll(func(elt *Element) string {
		return elt.Name
	}, []*Element{
		{Name: "n1", Value: "v1"},
		{Name: "n2", Value: "v2"},
	})

	a.Equal(2, m.Len())
	a.Equal(2, m.Keys().Size())
	a.Len(m.Values(), 2)

	a.Equal("v1", m.Get("n1").Value)
	a.Equal("v2", m.Get("n2").Value)
}

func Test_OrderedMap_putIfAbsent(t *testing.T) {
	a := require.New(t)

	m := NewOrderedMap("")
	a.True(m.PutIfAbsent("n", "v"))
	a.False(m.PutIfAbsent("n", "v-changed"))

	a.Equal("v", m.Get("n"))
}

func Test_OrderedMap_sort(t *testing.T) {
	a := require.New(t)

	m := NewOrderedMap("")
	m.Put("k1", "v1")
	m.Put("k2", "v2")

	m.SortByKey(false)
	a.Equal(2, m.Len())
	a.True(m.Keys().Contains("k1"))
	a.True(m.Keys().Contains("k2"))
	a.Equal("v1", m.Values()[0])
	a.Equal("v2", m.Values()[1])

	m.SortByKey(true)
	a.Equal(2, m.Len())
	a.True(m.Keys().Contains("k1"))
	a.True(m.Keys().Contains("k2"))
	a.Equal("v2", m.Values()[0])
	a.Equal("v1", m.Values()[1])
}

func TestNewOrderedMap(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	a.NotNil(om)
	a.Equal(0, om.Len())
}

func TestOrderedMap_Put_Get(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("key1", "value1")
	om.Put("key2", "value2")

	a.Equal("value1", om.Get("key1"))
	a.Equal("value2", om.Get("key2"))
	a.Equal(2, om.Len())
}

func TestOrderedMap_Get_NotFound(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("default")
	result := om.Get("nonexistent")
	a.Equal("default", result)
}

func TestOrderedMap_Find(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[int](0)
	om.Put("key1", 100)

	value, exists := om.Find("key1")
	a.True(exists)
	a.Equal(100, value)

	value, exists = om.Find("key2")
	a.False(exists)
	a.Equal(0, value)
}

func TestOrderedMap_Has(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("key1", "value1")

	a.True(om.Has("key1"))
	a.False(om.Has("key2"))
}

func TestOrderedMap_PutIfAbsent(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")

	// First put should succeed
	added := om.PutIfAbsent("key1", "value1")
	a.True(added)
	a.Equal("value1", om.Get("key1"))

	// Second put with same key should fail
	added = om.PutIfAbsent("key1", "value2")
	a.False(added)
	a.Equal("value1", om.Get("key1"))
}

func TestOrderedMap_PutAll(t *testing.T) {
	a := require.New(t)

	type TestItem struct {
		Name  string
		Value int
	}

	om := NewOrderedMap[TestItem](TestItem{})

	items := []TestItem{
		{Name: "item1", Value: 1},
		{Name: "item2", Value: 2},
	}

	om.PutAll(func(v TestItem) string { return v.Name }, items)

	a.Equal(2, om.Len())
	a.Equal(1, om.Get("item1").Value)
	a.Equal(2, om.Get("item2").Value)
}

func TestOrderedMap_Delete(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("key1", "value1")
	om.Put("key2", "value2")

	a.Equal(2, om.Len())

	om.Delete("key1")
	a.Equal(1, om.Len())
	a.False(om.Has("key1"))
	a.True(om.Has("key2"))
}

func TestOrderedMap_Keys(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("key1", "value1")
	om.Put("key2", "value2")

	keys := om.Keys()
	a.Equal(2, keys.Size())
	a.True(keys.Contains("key1"))
	a.True(keys.Contains("key2"))
}

func TestOrderedMap_Values(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[int](0)
	om.Put("key1", 100)
	om.Put("key2", 200)

	values := om.Values()
	a.Equal(2, len(values))
	a.Contains(values, 100)
	a.Contains(values, 200)
}

func TestOrderedMap_Entries(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("key1", "value1")
	om.Put("key2", "value2")

	entries := om.Entries()
	a.Equal(2, len(entries))

	// Check first entry
	a.Equal("key1", entries[0].Key)
	a.Equal("value1", entries[0].Value)

	// Check second entry
	a.Equal("key2", entries[1].Key)
	a.Equal("value2", entries[1].Value)
}

func TestOrderedMap_MarshalJSON(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("key1", "value1")
	om.Put("key2", "value2")

	data, err := json.Marshal(om)
	a.NoError(err)
	a.Contains(string(data), "key1")
	a.Contains(string(data), "value1")
}

func TestOrderedMap_UnmarshalJSON(t *testing.T) {
	a := require.New(t)

	jsonData := `{"key1":"value1","key2":"value2"}`

	om := NewOrderedMap[any](nil)
	err := json.Unmarshal([]byte(jsonData), om)
	a.NoError(err)

	a.Equal(2, om.Len())
	a.True(om.Has("key1"))
	a.True(om.Has("key2"))
}

func TestOrderedMap_SortByKey(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")
	om.Put("z", "value_z")
	om.Put("a", "value_a")
	om.Put("m", "value_m")

	// Sort ascending
	om.SortByKey(false)
	entries := om.Entries()
	a.Equal("a", entries[0].Key)
	a.Equal("m", entries[1].Key)
	a.Equal("z", entries[2].Key)

	// Sort descending
	om.SortByKey(true)
	entries = om.Entries()
	a.Equal("z", entries[0].Key)
	a.Equal("m", entries[1].Key)
	a.Equal("a", entries[2].Key)
}

func TestOrderedMap_ToMap(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[int](0)
	om.Put("key1", 100)
	om.Put("key2", 200)

	m := om.ToMap()
	a.Equal(2, len(m))
	a.Equal(100, m["key1"])
	a.Equal(200, m["key2"])
}

func TestOrderedMap_Empty(t *testing.T) {
	a := require.New(t)

	om := NewOrderedMap[string]("")

	a.Equal(0, om.Len())
	a.Equal(0, om.Keys().Size())
	a.Empty(om.Values())
	a.Empty(om.Entries())
	a.Empty(om.ToMap())
}
