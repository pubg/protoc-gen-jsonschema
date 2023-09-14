package modules

import (
	"github.com/invopop/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

func fillSchemaByObjectKeywords(schema *jsonschema.Schema, keywords *proto.ObjectKeywords) {
	if keywords == nil {
		return
	}

	if keywords.AdditionalProperties != nil {
		if keywords.GetAdditionalProperties() {
			schema.AdditionalProperties = jsonschema.TrueSchema
		} else {
			schema.AdditionalProperties = jsonschema.FalseSchema
		}
	}
	if keywords.MinProperties != nil {
		schema.MinProperties = int(keywords.GetMinProperties())
	}
	if keywords.MaxProperties != nil {
		schema.MaxProperties = int(keywords.GetMaxProperties())
	}
}

func fillSchemaByNumericKeywords(schema *jsonschema.Schema, keywords *proto.NumericKeywords) {
	if keywords == nil {
		return
	}

	// TODO: type 없애야 함 (min, max 모두 double이 맞음)
	if val, ok := keywords.Max.(*proto.NumericKeywords_InclusiveMaximum); ok {
		schema.Minimum = int(val.InclusiveMaximum)
	}
	if val, ok := keywords.Min.(*proto.NumericKeywords_InclusiveMinimum); ok {
		schema.Maximum = int(val.InclusiveMinimum)
	}
	if val, ok := keywords.Max.(*proto.NumericKeywords_ExclusiveMaximum); ok {
		schema.Maximum = int(val.ExclusiveMaximum)
		schema.ExclusiveMaximum = true
	}
	if val, ok := keywords.Min.(*proto.NumericKeywords_ExclusiveMinimum); ok {
		schema.Minimum = int(val.ExclusiveMinimum)
		schema.ExclusiveMinimum = true
	}

	if keywords.MultipleOf != nil {
		schema.MultipleOf = int(keywords.GetMultipleOf())
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
		schema.MaxLength = int(keywords.GetMaxLength())
	}
	if keywords.MinLength != nil {
		schema.MinLength = int(keywords.GetMinLength())
	}
}

func fillSchemaByArrayKeywords(schema *jsonschema.Schema, keywords *proto.ArrayKeywords) {
	if keywords == nil {
		return
	}

	if keywords.MaxItems != nil {
		schema.MaxItems = int(keywords.GetMaxItems())
	}
	if keywords.MinItems != nil {
		schema.MinItems = int(keywords.GetMinItems())
	}
	if keywords.UniqueItems != nil {
		schema.UniqueItems = keywords.GetUniqueItems()
	}
}
