package modules

import (
	"github.com/invopop/jsonschema"
)

type Optimizer interface {
	Optimize(schema *jsonschema.Schema) *jsonschema.Schema
}

var _ Optimizer = (*OptimizerImpl)(nil)

type OptimizerImpl struct {
	schemaByRef map[string]*jsonschema.Schema
}

func NewOptimizerImpl() *OptimizerImpl {
	return &OptimizerImpl{}
}

const refCountKey = "refCount"

func (o *OptimizerImpl) Optimize(rootSchema *jsonschema.Schema) *jsonschema.Schema {
	// TODO: Need Deepcopy all schema
	o.buildSchemaByRef(rootSchema)

	o.increaseSchemaRefCount(rootSchema.Ref)
	o.visitSchema(o.schemaByRef[rootSchema.Ref])

	o.optimizeDefinitions(rootSchema)

	o.cleanupSchemas()
	return rootSchema
}

func (o *OptimizerImpl) cleanupSchemas() {
	for _, schema := range o.schemaByRef {
		if schema.Extras != nil {
			delete(schema.Extras, refCountKey)
		}
	}
}

func (o *OptimizerImpl) optimizeDefinitions(rootSchema *jsonschema.Schema) {
	newDefinitions := jsonschema.Definitions{}
	for key, schema := range rootSchema.Definitions {
		if schema.Extras == nil {
			continue
		}
		refCount, found := schema.Extras[refCountKey]
		if !found {
			continue
		}
		if refCount.(int) == 0 {
			continue
		}
		newDefinitions[key] = schema
	}
	rootSchema.Definitions = newDefinitions
}

func (o *OptimizerImpl) buildSchemaByRef(rootSchema *jsonschema.Schema) {
	o.schemaByRef = map[string]*jsonschema.Schema{}
	for ref, schema := range rootSchema.Definitions {
		o.schemaByRef[jsonschema.EmptyID.Def(ref).String()] = schema
	}
}

func (o *OptimizerImpl) increaseSchemaRefCount(ref string) {
	schema := o.schemaByRef[ref]
	if schema.Extras == nil {
		schema.Extras = map[string]interface{}{}
		schema.Extras[refCountKey] = int(1)
	} else if _, exist := schema.Extras[refCountKey]; exist {
		schema.Extras[refCountKey] = schema.Extras[refCountKey].(int) + 1
	} else {
		schema.Extras[refCountKey] = int(1)
	}
}

func (o *OptimizerImpl) visitSchema(schema *jsonschema.Schema) {
	if schema == nil {
		return
	}

	if schema.Ref != "" {
		o.increaseSchemaRefCount(schema.Ref)
		o.visitSchema(o.schemaByRef[schema.Ref])
	}
	o.visitSchemaMap(schema.Definitions)

	o.visitSchemaArray(schema.AllOf)
	o.visitSchemaArray(schema.AnyOf)
	o.visitSchemaArray(schema.OneOf)
	o.visitSchema(schema.Not)

	o.visitSchema(schema.If)
	o.visitSchema(schema.Then)
	o.visitSchema(schema.Else)
	o.visitSchemaMap(schema.DependentSchemas)

	o.visitSchemaArray(schema.PrefixItems)
	o.visitSchema(schema.Items)
	o.visitSchema(schema.Contains)

	if schema.Properties != nil {
		for _, key := range schema.Properties.Keys() {
			schema, _ := schema.Properties.Get(key)
			o.visitSchema(schema.(*jsonschema.Schema))
		}
	}
	o.visitSchemaMap(schema.PatternProperties)
	o.visitSchema(schema.AdditionalProperties)
	o.visitSchema(schema.PropertyNames)

	o.visitSchema(schema.ContentSchema)
}

func (o *OptimizerImpl) visitSchemaArray(schemas []*jsonschema.Schema) {
	for _, schema := range schemas {
		o.visitSchema(schema)
	}
}

func (o *OptimizerImpl) visitSchemaMap(schemas map[string]*jsonschema.Schema) {
	for _, schema := range schemas {
		o.visitSchema(schema)
	}
}
