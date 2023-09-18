package modules

import (
	"encoding/json"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/anypb"
)

func buildFromMessage(pluginOptions *proto.PluginOptions, message pgs.Message, mo *proto.MessageOptions) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	schema.Type = "object"
	schema.Title = proto.GetTitleOrEmpty(mo)
	schema.Description = proto.GetDescriptionOrComment(message, mo)
	schema.Properties = jsonschema.NewOrderedSchemaMap()

	fillSchemaByObjectKeywords(schema, mo.GetObject())

	for _, field := range message.Fields() {
		// Skip OneOf Block field
		//if field.OneOf() != nil {
		//	continue
		//}

		propName := toPropertyName(field.Name())
		fieldSchema := &jsonschema.Schema{Ref: toRefId(field)}
		if !pluginOptions.GetMandatoryNullable() && (field.InRealOneOf() || field.HasOptionalKeyword()) || proto.GetFieldOptions(field).GetNullable() {
			fieldSchema = &jsonschema.Schema{OneOf: []*jsonschema.Schema{
				{Type: "null"},
				fieldSchema,
			}}
		}

		if !field.InRealOneOf() && !field.HasOptionalKeyword() && !field.Type().IsRepeated() && !field.Type().IsMap() {
			// If field is not a member of oneOf
			schema.Required = append(schema.Required, propName)
		}

		schema.Properties.Set(propName, fieldSchema)
	}

	// Convert Protobuf OneOfs to JSONSchema keywords
	for _, oneOf := range message.OneOfs() {
		propertyNames := lo.Map[pgs.Field, string](oneOf.Fields(), func(item pgs.Field, _ int) string {
			return toPropertyName(item.Name())
		})
		oneOfSchemas := lo.Map[string, *jsonschema.Schema](propertyNames, func(item string, _ int) *jsonschema.Schema {
			return &jsonschema.Schema{Required: []string{item}}
		})

		negativeSchema := &jsonschema.Schema{Not: &jsonschema.Schema{AnyOf: make([]*jsonschema.Schema, len(oneOfSchemas))}}
		copy(negativeSchema.Not.AnyOf, oneOfSchemas)

		combinedSchemas := append(oneOfSchemas, negativeSchema)
		schema.AllOf = append(schema.AllOf, &jsonschema.Schema{OneOf: combinedSchemas})
	}
	return schema
}

func buildFromMessageField(field pgs.Field, fo *proto.FieldOptions) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	schema.Title = proto.GetTitleOrEmpty(fo)
	schema.Description = proto.GetDescriptionOrComment(field, fo)

	if field.Type().IsRepeated() {
		schema.Ref = toRefId(field.Type().Element().Embed())
	} else {
		schema.Ref = toRefId(field.Type().Embed())
	}

	if field.Type().IsRepeated() {
		return wrapSchemaInArray(schema, field, fo)
	} else {
		return schema
	}
}

func buildFromMapField(pluginOptions *proto.PluginOptions, field pgs.Field, fo *proto.FieldOptions) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	schema.Title = proto.GetTitleOrEmpty(fo)
	schema.Description = proto.GetDescriptionOrComment(field, fo)
	schema.Type = "object"

	valueSchema := &jsonschema.Schema{}
	value := field.Type().Element()
	protoType := value.ProtoType()
	if protoType.IsInt() {
		if pluginOptions.GetInt64AsString() && isInt64(protoType) {
			schema.Type = "string"
			schema.Format = "int64"
		} else {
			schema.Type = "integer"
			fillSchemaByNumericKeywords(schema, fo.GetNumeric())
		}
	} else if protoType.IsNumeric() {
		valueSchema.Type = "number"
	} else if protoType == pgs.MessageT {
		valueSchema.Ref = toRefId(value.Embed())
	} else if protoType == pgs.BoolT {
		valueSchema.Type = "boolean"
	} else if protoType == pgs.EnumT {
		valueSchema.Ref = toRefId(value.Enum())
	} else if protoType == pgs.StringT {
		valueSchema.Type = "string"
	} else if protoType == pgs.BytesT {
		valueSchema.Type = "string"
		// Base64 Regex Expression
		valueSchema.Pattern = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$"
	}
	schema.AdditionalProperties = valueSchema
	return schema
}

func buildFromScalaField(pluginOptions *proto.PluginOptions, field pgs.Field, fo *proto.FieldOptions) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	schema.Title = proto.GetTitleOrEmpty(fo)
	schema.Description = proto.GetDescriptionOrComment(field, fo)

	protoType := field.Type().ProtoType()
	if protoType.IsInt() {
		if pluginOptions.GetInt64AsString() && isInt64(protoType) {
			schema.Type = "string"
			schema.Format = "int64"
		} else {
			schema.Type = "integer"
			fillSchemaByNumericKeywords(schema, fo.GetNumeric())
		}
	} else if protoType == pgs.DoubleT || protoType == pgs.FloatT {
		schema.Type = "number"
		fillSchemaByNumericKeywords(schema, fo.GetNumeric())
	} else if protoType == pgs.BoolT {
		schema.Type = "boolean"
	} else if protoType == pgs.EnumT {
		if field.Type().IsRepeated() {
			schema.Ref = toRefId(field.Type().Element().Enum())
		} else {
			schema.Ref = toRefId(field.Type().Enum())
		}
	} else if protoType == pgs.StringT {
		schema.Type = "string"
		fillSchemaByStringKeywords(schema, fo.GetString_())
	} else if protoType == pgs.BytesT {
		schema.Type = "string"
		// Base64 Regex Expression
		schema.Pattern = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$"
		fillSchemaByStringKeywords(schema, fo.GetString_())
	}

	if field.Type().IsRepeated() {
		return wrapSchemaInArray(schema, field, fo)
	} else {
		return schema
	}
}

func buildFromEnum(enum pgs.Enum) (*jsonschema.Schema, error) {
	eo := proto.GetEnumOptions(enum)

	schema := &jsonschema.Schema{}
	switch eo.GetMappingType() {
	case proto.EnumOptions_MapToString:
		schema.Type = "string"
	case proto.EnumOptions_MapToNumber:
		schema.Type = "number"
	case proto.EnumOptions_MapToCustom:
		schema.Type = "string"
	}
	schema.Title = proto.GetTitleOrEmpty(eo)
	schema.Description = proto.GetDescriptionOrComment(enum, eo)

	for _, enumValue := range enum.Values() {
		switch eo.GetMappingType() {
		case proto.EnumOptions_MapToString:
			schema.Enum = append(schema.Enum, enumValue.Name().String())
		case proto.EnumOptions_MapToNumber:
			schema.Enum = append(schema.Enum, enumValue.Value())
		case proto.EnumOptions_MapToCustom:
			evo := proto.GetEnumValueOptions(enumValue)

			customValue, err := parseScalaValueFromAny(evo.GetCustomValue())
			if err != nil {
				return nil, err
			}

			if customValue == nil {
				schema.Enum = append(schema.Enum, enumValue.Name().String())
			} else {
				schema.Enum = append(schema.Enum, customValue)
			}
		}
	}
	return schema, nil
}

func wrapSchemaInArray(schema *jsonschema.Schema, field pgs.Field, fo *proto.FieldOptions) *jsonschema.Schema {
	repeatedSchema := &jsonschema.Schema{}
	repeatedSchema.Title = schema.Title
	repeatedSchema.Description = schema.Description
	repeatedSchema.Type = "array"
	repeatedSchema.Items = schema

	fillSchemaByArrayKeywords(repeatedSchema, fo.GetArray())
	return repeatedSchema
}

func buildFromWellKnownField(field pgs.Field, fo *proto.FieldOptions) *jsonschema.Schema {
	schema := &jsonschema.Schema{}
	schema.Title = proto.GetTitleOrEmpty(fo)
	schema.Description = proto.GetDescriptionOrComment(field, fo)

	wellKnownType := getWellKnownType(field)
	if wellKnownType == WellKnownTypeNone {
		panic("not well known type")
	}

	switch wellKnownType {
	case WellKnownTypeTimestamp:
		schema.Type = "string"
		schema.Format = "date-time"
		fillSchemaByStringKeywords(schema, fo.GetString_())
	case WellKnownTypeDuration:
		schema.Type = "string"
		schema.Format = "duration"
		fillSchemaByStringKeywords(schema, fo.GetString_())
	case WellKnownTypeAny:
		schema.Type = "object"
	case WellKnownTypeNullValue:
		schema.Type = "null"
	case WellKnownTypeK8sIntOrString:
		schema.OneOf = []*jsonschema.Schema{
			{Type: "string"},
			{Type: "integer"},
		}
	}

	if field.Type().IsRepeated() {
		return wrapSchemaInArray(schema, field, fo)
	} else {
		return schema
	}
	return schema
}

func isWellKnownType(field pgs.Field) bool {
	return getWellKnownType(field) != WellKnownTypeNone
}

type WellKnownType int

const (
	WellKnownTypeNone WellKnownType = iota
	WellKnownTypeTimestamp
	WellKnownTypeDuration
	WellKnownTypeAny
	WellKnownTypeNullValue
	WellKnownTypeK8sIntOrString
)

func getWellKnownType(field pgs.Field) WellKnownType {
	if field.Type().IsMap() || field.Type().ProtoType() != pgs.MessageT {
		return WellKnownTypeNone
	}
	var message pgs.Message
	if field.Type().IsRepeated() {
		message = field.Type().Element().Embed()
	} else {
		message = field.Type().Embed()
	}

	switch message.FullyQualifiedName() {
	case ".google.protobuf.Timestamp":
		return WellKnownTypeTimestamp
	case ".google.protobuf.Duration":
		return WellKnownTypeDuration
	case ".google.protobuf.Any":
		return WellKnownTypeAny
	case ".google.protobuf.NullValue":
		return WellKnownTypeNullValue
	case ".k8s.io.apimachinery.pkg.util.intstr.IntOrString":
		return WellKnownTypeK8sIntOrString
	default:
		return WellKnownTypeNone
	}
}

func parseScalaValueFromAny(anyValue *anypb.Any) (any, error) {
	if anyValue == nil || anyValue.Value == nil {
		return nil, nil
	}

	var value any
	if err := json.Unmarshal(anyValue.Value, &value); err != nil {
		return nil, err
	}
	return value, nil
}

func isScalarType(field pgs.Field) bool {
	protoType := field.Type().ProtoType()
	if protoType.IsNumeric() {
		return true
	}
	if protoType == pgs.BoolT {
		return true
	}
	if protoType == pgs.EnumT {
		return true
	}
	if protoType == pgs.StringT || protoType == pgs.BytesT {
		return true
	}
	return false
}

func toPropertyName(name pgs.Name) string {
	return name.String()
}

type FqdnResolver interface {
	FullyQualifiedName() string
}

func toRefId(resolver FqdnResolver) jsonschema.RefId {
	return jsonschema.RefId(resolver.FullyQualifiedName())
}

func isInt64(protoType pgs.ProtoType) bool {
	switch protoType {
	case pgs.Int64T, pgs.UInt64T, pgs.SFixed64, pgs.SInt64, pgs.Fixed64T:
		return true
	}
	return false
}
