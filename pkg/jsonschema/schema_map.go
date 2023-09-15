package jsonschema

import "github.com/iancoleman/orderedmap"

type SchemaMap interface {
	Get(key string) (*Schema, bool)
	Set(key string, value *Schema)
	Keys() []string
	Delete(key string)
}

type SimpleSchemaMap map[string]*Schema

func (m SimpleSchemaMap) Get(key string) (*Schema, bool) {
	v, found := m[key]
	return v, found
}

func (m SimpleSchemaMap) Set(key string, value *Schema) {
	m[key] = value
}

func (m SimpleSchemaMap) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func (m SimpleSchemaMap) Delete(key string) {
	delete(m, key)
}

type orderedSchemaMap struct {
	memory *orderedmap.OrderedMap
}

func NewOrderedSchemaMap() SchemaMap {
	return &orderedSchemaMap{memory: orderedmap.New()}
}

func (m *orderedSchemaMap) Get(key string) (*Schema, bool) {
	if m == nil {
		return nil, false
	}

	v, found := m.memory.Get(key)
	if !found {
		return nil, false
	}
	return v.(*Schema), true
}

func (m *orderedSchemaMap) Set(key string, value *Schema) {
	m.memory.Set(key, value)
}

func (m *orderedSchemaMap) Keys() []string {
	if m == nil {
		return nil
	}
	return m.memory.Keys()
}

func (m *orderedSchemaMap) Delete(key string) {
	m.memory.Delete(key)
}
