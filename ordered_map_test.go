package comm

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

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
