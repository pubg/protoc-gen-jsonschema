package modules

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/pjsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

type JsonSchemaVisitor struct {
	pgs.Visitor

	debugger pgs.DebuggerCommon

	registry      *pjsonschema.Registry
	pluginOptions *proto.PluginOptions
}

var _ pgs.Visitor = (*JsonSchemaVisitor)(nil)

func NewJsonSchemaVisitor(debugger pgs.DebuggerCommon) *JsonSchemaVisitor {
	v := &JsonSchemaVisitor{
		debugger: debugger,
		registry: pjsonschema.NewRegistry(),
	}
	v.Visitor = pgs.PassThroughVisitor(v)
	return v
}

func (jv *JsonSchemaVisitor) VisitMessage(message pgs.Message) (v pgs.Visitor, err error) {
	mo := proto.GetMessageOptions(message)
	if mo.GetVisibilityLevel() < jv.pluginOptions.GetVisibilityLevel() {
		return nil, nil
	}

	schema := buildFromMessage(message, mo)
	jv.registry.AddSchema(message.FullyQualifiedName(), schema)
	return jv, nil
}

func (jv *JsonSchemaVisitor) VisitField(field pgs.Field) (v pgs.Visitor, err error) {
	fo := proto.GetFieldOptions(field)
	if fo.GetVisibilityLevel() < jv.pluginOptions.GetVisibilityLevel() {
		return nil, nil
	}

	// TODO: if field is well-known type

	// if field is message or map type
	fieldType := field.Type()
	if fieldType.IsMap() {
		schema := buildFromMapField(field, fo)
		jv.registry.AddSchema(field.FullyQualifiedName(), schema)
		return jv, nil
	} else if fieldType.ProtoType() == pgs.MessageT {
		schema := buildFromMessageField(field, fo)
		jv.registry.AddSchema(field.FullyQualifiedName(), schema)
		return jv, nil
	}

	// if field is scala type
	// scala = boolean, string, number
	if isScalarType(field) {
		schema := buildFromScalaField(field, fo)
		jv.registry.AddSchema(field.FullyQualifiedName(), schema)
		return jv, nil
	}

	panic("not supported field type")
	return jv, nil
}

func (jv *JsonSchemaVisitor) VisitEnum(enum pgs.Enum) (v pgs.Visitor, err error) {
	schema, err := buildFromEnum(enum)
	if err != nil {
		return nil, err
	}
	jv.registry.AddSchema(enum.FullyQualifiedName(), schema)
	return jv, nil
}
