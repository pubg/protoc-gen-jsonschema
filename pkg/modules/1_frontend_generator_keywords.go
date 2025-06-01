package modules

import (
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
	"github.com/pubg/protoc-gen-jsonschema/pkg/utils"
)

func fillSchemaByObjectKeywords(pluginOptions *proto.PluginOptions, schema *jsonschema.Schema, keywords *proto.ObjectKeywords) {
	if additionalProperties := proto.GetAdditionalProperties(pluginOptions, keywords); additionalProperties != nil {
		schema.AdditionalProperties = jsonschema.NewBooleanSchema(*additionalProperties)
	}

	if keywords == nil {
		return
	}
	if keywords.MinProperties != nil {
		schema.MinProperties = utils.UInt32(keywords.GetMinProperties())
	}
	if keywords.MaxProperties != nil {
		schema.MaxProperties = utils.UInt32(keywords.GetMaxProperties())
	}
}

func fillSchemaByNumericKeywords(schema *jsonschema.Schema, keywords *proto.NumericKeywords) {
	if keywords == nil {
		return
	}

	if val, ok := keywords.Min.(*proto.NumericKeywords_InclusiveMinimum); ok {
		schema.Minimum = &val.InclusiveMinimum
	}
	if val, ok := keywords.Max.(*proto.NumericKeywords_InclusiveMaximum); ok {
		schema.Maximum = &val.InclusiveMaximum
	}
	if val, ok := keywords.Max.(*proto.NumericKeywords_ExclusiveMaximum); ok {
		schema.ExclusiveMaximum = &val.ExclusiveMaximum
	}
	if val, ok := keywords.Min.(*proto.NumericKeywords_ExclusiveMinimum); ok {
		schema.ExclusiveMinimum = &val.ExclusiveMinimum
	}

	if keywords.MultipleOf != nil {
		schema.MultipleOf = utils.Int32(keywords.GetMultipleOf())
	}
}

func fillSchemaByStringKeywords(schema *jsonschema.Schema, keywords *proto.StringKeywords) {
	if keywords == nil {
		return
	}

	if keywords.Pattern != "" {
		schema.Pattern = keywords.GetPattern()
	}
	if keywords.Format != "" {
		schema.Format = keywords.GetFormat()
	}

	if keywords.MaxLength != nil {
		schema.MaxLength = utils.UInt32(keywords.GetMaxLength())
	}
	if keywords.MinLength != nil {
		schema.MinLength = utils.UInt32(keywords.GetMinLength())
	}
}

func fillSchemaByArrayKeywords(schema *jsonschema.Schema, keywords *proto.ArrayKeywords) {
	if keywords == nil {
		return
	}

	if keywords.MaxItems != nil {
		schema.MaxItems = utils.UInt32(keywords.GetMaxItems())
	}
	if keywords.MinItems != nil {
		schema.MinItems = utils.UInt32(keywords.GetMinItems())
	}
	if keywords.UniqueItems != nil {
		schema.UniqueItems = keywords.UniqueItems
	}
}
