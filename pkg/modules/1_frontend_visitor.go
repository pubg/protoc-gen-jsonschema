package modules

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

// MiddleendVisitor generate intermediate jsonschema from protobuf
type MiddleendVisitor struct {
	pgs.Visitor

	debugger pgs.DebuggerCommon

	registry      *jsonschema.Registry
	pluginOptions *proto.PluginOptions
}

var _ pgs.Visitor = (*MiddleendVisitor)(nil)

func NewVisitor(debugger pgs.DebuggerCommon, pluginOptions *proto.PluginOptions) *MiddleendVisitor {
	v := &MiddleendVisitor{
		debugger:      debugger,
		registry:      jsonschema.NewRegistry(),
		pluginOptions: pluginOptions,
	}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

func (v *MiddleendVisitor) VisitMessage(message pgs.Message) (pgs.Visitor, error) {
	mo := proto.GetMessageOptions(message)
	if mo.GetVisibilityLevel() < v.pluginOptions.GetVisibilityLevel() {
		return nil, nil
	}

	var schema *jsonschema.Schema
	// if message is well-known type
	if isWellKnownMessage(message) {
		schema = buildFromWellKnownMessage(v.pluginOptions, message, mo)
	} else {
		schema = buildFromMessage(v.pluginOptions, message, mo)
	}
	v.registry.AddSchema(message.FullyQualifiedName(), schema)
	return v, nil
}

func (v *MiddleendVisitor) VisitField(field pgs.Field) (pgs.Visitor, error) {
	fo := proto.GetFieldOptions(field)
	if fo.GetVisibilityLevel() < v.pluginOptions.GetVisibilityLevel() {
		return nil, nil
	}

	// if field is well-known type
	if isWellKnownField(field) {
		schema := buildFromWellKnownField(field, fo)
		v.registry.AddSchema(field.FullyQualifiedName(), schema)
		return v, nil
	}

	// if field is message or map type
	fieldType := field.Type()
	if fieldType.IsMap() {
		schema := buildFromMapField(v.pluginOptions, field, fo)
		v.registry.AddSchema(field.FullyQualifiedName(), schema)
		return v, nil
	} else if fieldType.ProtoType() == pgs.MessageT {
		schema := buildFromMessageField(field, fo)
		v.registry.AddSchema(field.FullyQualifiedName(), schema)
		return v, nil
	}

	// if field is scala type
	// scala = boolean, string, number
	if isScalarType(field) {
		schema := buildFromScalaField(v.pluginOptions, field, fo)
		v.registry.AddSchema(field.FullyQualifiedName(), schema)
		return v, nil
	}

	panic("not supported field type")
	return v, nil
}

func (v *MiddleendVisitor) VisitEnum(enum pgs.Enum) (pgs.Visitor, error) {
	schema, err := buildFromEnum(enum)
	if err != nil {
		return nil, err
	}
	v.registry.AddSchema(enum.FullyQualifiedName(), schema)
	return v, nil
}
