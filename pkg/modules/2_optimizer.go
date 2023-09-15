package modules

import (
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

type Optimizer interface {
	Optimize(registry *jsonschema.Registry, file pgs.File)
}

var _ Optimizer = (*OptimizerImpl)(nil)

type OptimizerImpl struct {
	schemaByRef   jsonschema.SchemaMap
	pluginOptions *proto.PluginOptions
}

func NewOptimizerImpl(pluginOptions *proto.PluginOptions) *OptimizerImpl {
	return &OptimizerImpl{schemaByRef: jsonschema.NewOrderedSchemaMap(), pluginOptions: pluginOptions}
}

const refCountKey = "refCount"

func (o *OptimizerImpl) Optimize(registry *jsonschema.Registry, file pgs.File) {
	entrypointMessage := o.getEntrypointMessage(file.Messages(), proto.GetFileOptions(file))
	entrypointSchemaRef := toRefId(entrypointMessage)
	entrypointSchema := registry.GetSchema(entrypointSchemaRef.String())

	o.increaseSchemaRefCount(registry, entrypointSchemaRef.String())
	o.visitSchema(registry, entrypointSchema)

	o.optimizeDefinitions(registry)
}

func (o *OptimizerImpl) optimizeDefinitions(registry *jsonschema.Registry) {
	var deleteKeys []string
	for _, key := range registry.GetKeys() {
		schema := registry.GetSchema(key)
		if schema.GetExtrasItem(refCountKey) == nil {
			deleteKeys = append(deleteKeys, key)
		} else if schema.GetExtrasItem(refCountKey).(int) == 0 {
			deleteKeys = append(deleteKeys, key)
		}
	}

	for _, key := range deleteKeys {
		registry.DeleteSchema(key)
	}
}

func (o *OptimizerImpl) idFromRef(ref string) string {
	return strings.Replace(ref, "#/definitions/", "", 1)
}

func (o *OptimizerImpl) getEntrypointMessage(messages []pgs.Message, fileOptions *proto.FileOptions) pgs.Message {
	entryPointMessage := proto.GetEntrypointMessage(o.pluginOptions, fileOptions)
	if entryPointMessage == "" {
		// TODO: print debug message
		return nil
	}

	for _, message := range messages {
		if message.Name().String() == entryPointMessage {
			return message
		}
	}
	return nil
}

func (o *OptimizerImpl) increaseSchemaRefCount(registry *jsonschema.Registry, ref string) {
	schema := registry.GetSchema(ref)

	rawValue := schema.GetExtrasItem(refCountKey)
	if rawValue == nil {
		schema.SetExtrasItem(refCountKey, int(1))
	} else {
		schema.SetExtrasItem(refCountKey, rawValue.(int)+1)
	}
}

func (o *OptimizerImpl) visitSchema(registry *jsonschema.Registry, schema *jsonschema.Schema) {
	if schema == nil {
		return
	}

	if !schema.Ref.IsEmpty() {
		o.increaseSchemaRefCount(registry, schema.Ref.String())
		o.visitSchema(registry, registry.GetSchema(schema.Ref.String()))
	}
	o.visitSchemaMap(registry, schema.Definitions)

	o.visitSchemaArray(registry, schema.AllOf)
	o.visitSchemaArray(registry, schema.AnyOf)
	o.visitSchemaArray(registry, schema.OneOf)
	o.visitSchema(registry, schema.Not)

	o.visitSchema(registry, schema.If)
	o.visitSchema(registry, schema.Then)
	o.visitSchema(registry, schema.Else)
	o.visitSchemaMap(registry, schema.DependentSchemas)

	o.visitSchemaArray(registry, schema.PrefixItems)
	o.visitSchema(registry, schema.Items)
	o.visitSchema(registry, schema.Contains)

	o.visitSchemaMap(registry, schema.Properties)
	o.visitSchemaMap(registry, schema.PatternProperties)
	o.visitSchema(registry, schema.AdditionalProperties)
	o.visitSchema(registry, schema.PropertyNames)

	o.visitSchema(registry, schema.ContentSchema)
}

func (o *OptimizerImpl) visitSchemaArray(registry *jsonschema.Registry, schemas []*jsonschema.Schema) {
	for _, schema := range schemas {
		o.visitSchema(registry, schema)
	}
}

func (o *OptimizerImpl) visitSchemaMap(registry *jsonschema.Registry, schemaMap jsonschema.SchemaMap) {
	if schemaMap == nil {
		return
	}
	for _, key := range schemaMap.Keys() {
		schema, _ := schemaMap.Get(key)
		o.visitSchema(registry, schema)
	}
}
