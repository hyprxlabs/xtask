package types

import (
	"iter"

	om "github.com/wk8/go-ordered-map/v2"
)

type PrimitiveMap struct {
	values *om.OrderedMap[string, PrimitiveValue]
}

type PrimitiveValue struct {
	Value interface{}
	Type  string
}

func NewPrimitiveMap() *PrimitiveMap {
	return &PrimitiveMap{
		values: &om.OrderedMap[string, PrimitiveValue]{},
	}
}

type ErrUnsupportedPrimitiveType struct {
	Key   string
	Value interface{}
}

func (e *ErrUnsupportedPrimitiveType) Error() string {
	return "unsupported primitive type for key " + e.Key
}

func (pm *PrimitiveMap) Has(key string) bool {
	_, ok := pm.values.Get(key)
	return ok
}

func (pm *PrimitiveMap) Len() int {
	return pm.values.Len()
}

func (pm *PrimitiveMap) Set(key string, value PrimitiveValue) {
	pm.values.Set(key, value)
}

func (pm *PrimitiveMap) Get(key string) (PrimitiveValue, bool) {
	return pm.values.Get(key)
}

func (pm *PrimitiveMap) Delete(key string) {
	pm.values.Delete(key)
}

func (pm *PrimitiveMap) AsMap() om.OrderedMap[string, PrimitiveValue] {
	return *pm.values
}

func (pm *PrimitiveMap) Keys() []string {
	keys := make([]string, 0, pm.values.Len())
	for el := pm.values.Oldest(); el != nil; el = el.Next() {
		keys = append(keys, el.Key)
	}
	return keys
}

func (pm *PrimitiveMap) Values() []PrimitiveValue {
	values := make([]PrimitiveValue, 0, pm.values.Len())
	for el := pm.values.Oldest(); el != nil; el = el.Next() {
		values = append(values, el.Value)
	}
	return values
}

func (pm *PrimitiveMap) Iter() iter.Seq2[string, PrimitiveValue] {
	return func(yield func(string, PrimitiveValue) bool) {
		for el := pm.values.Oldest(); el != nil; el = el.Next() {
			if !yield(el.Key, el.Value) {
				return
			}
		}
	}
}

func (pm *PrimitiveMap) SetString(key, value string) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "string"})
}

func (pm *PrimitiveMap) SetStrings(key string, value []string) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "[]string"})
}

func (pm *PrimitiveMap) SetI64(key string, value int64) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "int64"})
}

func (pm *PrimitiveMap) SetU64(key string, value uint64) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "int64"})
}

func (pm *PrimitiveMap) SetU32(key string, value uint32) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "uint32"})
}

func (pm *PrimitiveMap) SetBool(key string, value bool) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "bool"})
}

func (pm *PrimitiveMap) SetF64(key string, value float64) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "float64"})
}

func (pm *PrimitiveMap) SetF32(key string, value float32) {
	pm.values.Set(key, PrimitiveValue{Value: value, Type: "float32"})
}

func (pm *PrimitiveMap) GetStrings(key string) ([]string, bool) {
	v, ok := pm.values.Get(key)
	if !ok || v.Type != "[]string" {
		return nil, false
	}
	return v.Value.([]string), true
}
