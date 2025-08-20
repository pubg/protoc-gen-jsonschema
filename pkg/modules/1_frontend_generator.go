package modules

import (
	"encoding/json"
	"slices"

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

	fillSchemaByObjectKeywords(pluginOptions, schema, mo.GetObject())

	for _, field := range message.Fields() {
		// Skip OneOf Block field
		//if field.OneOf() != nil {
		//	continue
		//}

		propName := toPropertyName(field, pluginOptions.GetPreserveProtoFieldNames())
		fieldSchema := &jsonschema.Schema{Ref: toRefId(field)}
		if !pluginOptions.GetMandatoryNullable() && (field.InRealOneOf() || field.HasOptionalKeyword()) || proto.GetFieldOptions(field).GetNullable() {
			fieldSchema = &jsonschema.Schema{OneOf: []*jsonschema.Schema{
				{Type: "null"},
				fieldSchema,
			}}
		}

		if isRequiredField(pluginOptions, field) {
			// If field is not a member of oneOf
			schema.Required = append(schema.Required, propName)
		}

		schema.Properties.Set(propName, fieldSchema)
	}

	// Convert Protobuf OneOfs to JSONSchema keywords
	for _, oneOf := range message.OneOfs() {
		propertyNames := lo.Map[pgs.Field, string](oneOf.Fields(), func(item pgs.Field, _ int) string {
			return toPropertyName(item, pluginOptions.GetPreserveProtoFieldNames())
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

func isRequiredField(pluginOptions *proto.PluginOptions, field pgs.Field) bool {
	if pluginOptions.GetRespectProtojsonPresence() {
		// To see the below link for more details
		// https://github.com/protocolbuffers/protobuf/blob/main/docs/implementing_proto3_presence.md#to-test-whether-a-field-should-have-presence
		return !field.InRealOneOf() && field.HasPresence()
	} else {
		return !field.InRealOneOf() && !field.HasOptionalKeyword() && !field.Type().IsRepeated() && !field.Type().IsMap()
	}
}

func buildFromWellKnownMessage(pluginOptions *proto.PluginOptions, message pgs.Message, mo *proto.MessageOptions) *jsonschema.Schema {
	baseSchema := buildFromMessage(pluginOptions, message, mo)

	wellKnownType := getWellKnownMessageType(message)
	if wellKnownType == WellKnownMessageTypeNone {
		panic("not well known type")
	}
	switch wellKnownType {
	case WellKnownMessageTypeK8sIntOrString:
		schema := &jsonschema.Schema{}
		schema.Title = baseSchema.Title
		schema.Description = baseSchema.Description
		schema.OneOf = []*jsonschema.Schema{
			{Type: "string"},
			{Type: "integer"},
		}
		return schema
	case WellKnownMessageTypeK8sVolume:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.VolumeSource"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "volumeSource")
		baseSchema.Properties.Delete("volumeSource")
	case WellKnownMessageTypeK8sSecretProjection:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sConfigMapVolumeSource:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sConfigMapProjection:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sConfigMapKeySelector:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sSecretKeySelector:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sConfigMapEnvSource:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sSecretEnvSource:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.LocalObjectReference"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "localObjectReference")
		baseSchema.Properties.Delete("localObjectReference")
	case WellKnownMessageTypeK8sProbe:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.ProbeHandler"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "handler")
		baseSchema.Properties.Delete("handler")
	case WellKnownMessageTypeK8sEphemeralContainer:
		baseSchema.OneOf = []*jsonschema.Schema{{Ref: ".k8s.io.api.core.v1.EphemeralContainerCommon"}}
		baseSchema.Required = deletePropertyInRequired(baseSchema.Required, "ephemeralContainerCommon")
		baseSchema.Properties.Delete("ephemeralContainerCommon")
	}
	return baseSchema
}

type WellKnownMessageType int

const (
	WellKnownMessageTypeNone WellKnownMessageType = iota
	WellKnownMessageTypeK8sIntOrString
	WellKnownMessageTypeK8sVolume
	WellKnownMessageTypeK8sSecretProjection
	WellKnownMessageTypeK8sConfigMapVolumeSource
	WellKnownMessageTypeK8sConfigMapProjection
	WellKnownMessageTypeK8sConfigMapKeySelector
	WellKnownMessageTypeK8sSecretKeySelector
	WellKnownMessageTypeK8sConfigMapEnvSource
	WellKnownMessageTypeK8sSecretEnvSource
	WellKnownMessageTypeK8sProbe
	WellKnownMessageTypeK8sEphemeralContainer
)

func isWellKnownMessage(message pgs.Message) bool {
	return getWellKnownMessageType(message) != WellKnownMessageTypeNone
}

func getWellKnownMessageType(message pgs.Message) WellKnownMessageType {
	switch message.FullyQualifiedName() {
	case ".k8s.io.apimachinery.pkg.util.intstr.IntOrString":
		return WellKnownMessageTypeK8sIntOrString
	case ".k8s.io.api.core.v1.Volume":
		return WellKnownMessageTypeK8sVolume
	case ".k8s.io.api.core.v1.SecretProjection":
		return WellKnownMessageTypeK8sSecretProjection
	case ".k8s.io.api.core.v1.ConfigMapVolumeSource":
		return WellKnownMessageTypeK8sConfigMapVolumeSource
	case ".k8s.io.api.core.v1.ConfigMapProjection":
		return WellKnownMessageTypeK8sConfigMapProjection
	case ".k8s.io.api.core.v1.ConfigMapKeySelector":
		return WellKnownMessageTypeK8sConfigMapKeySelector
	case ".k8s.io.api.core.v1.SecretKeySelector":
		return WellKnownMessageTypeK8sSecretKeySelector
	case ".k8s.io.api.core.v1.ConfigMapEnvSource":
		return WellKnownMessageTypeK8sConfigMapEnvSource
	case ".k8s.io.api.core.v1.SecretEnvSource":
		return WellKnownMessageTypeK8sSecretEnvSource
	case ".k8s.io.api.core.v1.Probe":
		return WellKnownMessageTypeK8sProbe
	case ".k8s.io.api.core.v1.EphemeralContainer":
		return WellKnownMessageTypeK8sEphemeralContainer
	}

	return WellKnownMessageTypeNone
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
		if pluginOptions.GetRespectProtojsonInt64() && isInt64(protoType) {
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
		if pluginOptions.GetRespectProtojsonInt64() && isInt64(protoType) {
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

	wellKnownType := getWellKnownFieldType(field)
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
	}

	if field.Type().IsRepeated() {
		return wrapSchemaInArray(schema, field, fo)
	} else {
		return schema
	}
	return schema
}

func isWellKnownField(field pgs.Field) bool {
	return getWellKnownFieldType(field) != WellKnownTypeNone
}

type WellKnownFieldType int

const (
	WellKnownTypeNone WellKnownFieldType = iota
	WellKnownTypeTimestamp
	WellKnownTypeDuration
	WellKnownTypeAny
	WellKnownTypeNullValue
)

func getWellKnownFieldType(field pgs.Field) WellKnownFieldType {
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

func toPropertyName(field pgs.Field, preserveProto bool) string {
	if !preserveProto && field.Descriptor().JsonName != nil {
		return field.Descriptor().GetJsonName()
	}
	return field.Name().String()
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

func deletePropertyInRequired(required []string, item string) []string {
	return slices.DeleteFunc(required, func(val string) bool {
		return val == item
	})
}
