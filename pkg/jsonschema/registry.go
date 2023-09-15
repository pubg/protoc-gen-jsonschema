package jsonschema

type Registry struct {
	schemasByName SchemaMap
}

func NewRegistry() *Registry {
	return &Registry{schemasByName: NewOrderedSchemaMap()}
}

func (r *Registry) AddSchema(id string, schema *Schema) {
	r.schemasByName.Set(id, schema)
}

func (r *Registry) HasSchema(id string) bool {
	_, found := r.schemasByName.Get(id)
	return found
}

func (r *Registry) GetSchema(id string) *Schema {
	schema, _ := r.schemasByName.Get(id)
	return schema
}

func (r *Registry) GetKeys() []string {
	return r.schemasByName.Keys()
}

func (r *Registry) DeleteSchema(id string) {
	r.schemasByName.Delete(id)
}

func DeepCopyRegistry(registry *Registry) *Registry {
	newRegistry := NewRegistry()
	newRegistry.schemasByName = DeepCopyMap(registry.schemasByName)
	return newRegistry
}
