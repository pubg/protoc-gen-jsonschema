package pjsonschema

import (
	"github.com/iancoleman/orderedmap"
	"github.com/invopop/jsonschema"
)

type Registry struct {
	schemasByName *orderedmap.OrderedMap
}

func NewRegistry() *Registry {
	return &Registry{schemasByName: orderedmap.New()}
}

func (r *Registry) AddSchema(id string, schema *jsonschema.Schema) {
	r.schemasByName.Set(id, schema)
}

func (r *Registry) HasSchema(id string) bool {
	_, found := r.schemasByName.Get(id)
	return found
}

func (r *Registry) GetSchema(id string) *jsonschema.Schema {
	schema, _ := r.schemasByName.Get(id)
	return schema.(*jsonschema.Schema)
}

func (r *Registry) GetKeys() []string {
	return r.schemasByName.Keys()
}

func Walk(registry *Registry, entryPoint string) {

}
