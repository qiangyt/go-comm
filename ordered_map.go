package comm

import (
	"github.com/iancoleman/orderedmap"
)

type OrderedMap[K any] struct {
	backend  *orderedmap.OrderedMap
	nilValue K
}

type KeyValue[K any] struct {
	Key   string
	Value K
}

func NewOrderedMap[K any](nilValue K) *OrderedMap[K] {
	backend := orderedmap.New()
	backend.SetEscapeHTML(false)

	return &OrderedMap[K]{
		backend:  backend,
		nilValue: nilValue,
	}
}

func (me *OrderedMap[K]) Find(key string) (K, bool) {
	r, exists := me.backend.Get(key)
	if !exists {
		return me.nilValue, false
	}
	return r.(K), true
}

func (me *OrderedMap[K]) Get(key string) K {
	r, exists := me.backend.Get(key)
	if !exists {
		return me.nilValue
	}
	return r.(K)
}

func (me *OrderedMap[K]) Has(key string) bool {
	_, exists := me.backend.Get(key)
	return exists
}

func (me *OrderedMap[K]) Len() int {
	return len(me.backend.Keys())
}

func (me *OrderedMap[K]) Set(key string, value K) {
	me.backend.Set(key, value)
}

func (me *OrderedMap[K]) Delete(key string) {
	me.backend.Delete(key)
}

func (me *OrderedMap[K]) Keys() []string {
	return me.backend.Keys()
}

func (me *OrderedMap[K]) Values() []K {
	r := make([]K, 0, me.Len())

	for _, k := range me.backend.Keys() {
		v, _ := me.Find(k)
		r = append(r, v)
	}

	return r
}

func (me *OrderedMap[K]) Entries() []*KeyValue[K] {
	r := make([]*KeyValue[K], 0, me.Len())

	for _, k := range me.backend.Keys() {
		v, exists := me.Find(k)
		if exists {
			kv := &KeyValue[K]{
				Key:   k,
				Value: v,
			}
			r = append(r, kv)
		}
	}

	return r
}

func (me *OrderedMap[K]) UnmarshalJSON(bytes []byte) error {
	return me.backend.UnmarshalJSON(bytes)
}

func (me *OrderedMap[K]) MarshalJSON() ([]byte, error) {
	return me.backend.MarshalJSON()
}
