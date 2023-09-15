package jsonschema

import "github.com/pubg/protoc-gen-jsonschema/pkg/utils"

type Schema struct {
	Version       string
	ID            string
	Anchor        string
	DynamicAnchor string
	Ref           RefId
	DynamicRef    string
	Definitions   SchemaMap
	Comments      string

	AllOf []*Schema
	AnyOf []*Schema
	OneOf []*Schema
	Not   *Schema

	If               *Schema
	Then             *Schema
	Else             *Schema
	DependentSchemas SchemaMap

	PrefixItems []*Schema
	Items       *Schema
	Contains    *Schema

	Properties           SchemaMap
	PatternProperties    SchemaMap
	AdditionalProperties *Schema
	PropertyNames        *Schema

	Type              string
	Enum              []any
	Const             *any
	MultipleOf        *int
	Maximum           *float64
	ExclusiveMaximum  *float64
	Minimum           *float64
	ExclusiveMinimum  *float64
	MaxLength         *int
	MinLength         *int
	Pattern           string
	MaxItems          *int
	MinItems          *int
	UniqueItems       *bool
	MaxContains       *int
	MinContains       *int
	MaxProperties     *int
	MinProperties     *int
	Required          []string
	DependentRequired map[string][]string

	Format string

	ContentEncoding  string
	ContentMediaType string
	ContentSchema    *Schema

	Title       string
	Description string
	Default     *any
	Deprecated  *bool
	ReadOnly    *bool
	WriteOnly   *bool
	Examples    []any

	Extras map[string]any

	// Indicated boolean schema or generic schema
	boolean *bool
}

func (s *Schema) SetExtrasItem(key string, value any) {
	if s.Extras == nil {
		s.Extras = map[string]any{}
	}
	s.Extras[key] = value
}

func (s *Schema) ClearExtras() {
	s.Extras = nil
}

func (s *Schema) GetExtrasItem(key string) any {
	if s.Extras == nil {
		return nil
	}
	return s.Extras[key]
}

func DeepCopy(origin *Schema) *Schema {
	if origin == nil {
		return nil
	}

	dst := &Schema{}
	dst.Version = origin.Version
	dst.ID = origin.ID
	dst.Anchor = origin.Anchor
	dst.DynamicAnchor = origin.DynamicAnchor
	dst.Ref = origin.Ref
	dst.DynamicRef = origin.DynamicRef
	dst.Definitions = DeepCopyMap(origin.Definitions)
	dst.Comments = origin.Comments

	dst.AllOf = DeepCopyArray(origin.AllOf)
	dst.AnyOf = DeepCopyArray(origin.AnyOf)
	dst.OneOf = DeepCopyArray(origin.OneOf)
	dst.Not = DeepCopy(origin.Not)

	dst.If = DeepCopy(origin.If)
	dst.Then = DeepCopy(origin.Then)
	dst.Else = DeepCopy(origin.Else)
	dst.DependentSchemas = DeepCopyMap(origin.DependentSchemas)

	dst.PrefixItems = DeepCopyArray(origin.PrefixItems)
	dst.Items = DeepCopy(origin.Items)
	dst.Contains = DeepCopy(origin.Contains)

	dst.Properties = DeepCopyMap(origin.Properties)
	dst.PatternProperties = DeepCopyMap(origin.PatternProperties)
	dst.AdditionalProperties = DeepCopy(origin.AdditionalProperties)
	dst.PropertyNames = DeepCopy(origin.PropertyNames)

	dst.Type = origin.Type
	dst.Enum = utils.CopyAnyArray(origin.Enum)
	dst.Const = utils.CopyAnyP(origin.Const)
	dst.MultipleOf = utils.CopyIntP(origin.MultipleOf)
	dst.Maximum = utils.CopyFloat64P(origin.Maximum)
	dst.ExclusiveMaximum = utils.CopyFloat64P(origin.ExclusiveMaximum)
	dst.Minimum = utils.CopyFloat64P(origin.Minimum)
	dst.ExclusiveMinimum = utils.CopyFloat64P(origin.ExclusiveMinimum)
	dst.MaxLength = utils.CopyIntP(origin.MaxLength)
	dst.MinLength = utils.CopyIntP(origin.MinLength)
	dst.Pattern = origin.Pattern
	dst.MaxItems = utils.CopyIntP(origin.MaxItems)
	dst.MinItems = utils.CopyIntP(origin.MinItems)
	dst.UniqueItems = utils.CopyBoolP(origin.UniqueItems)
	dst.MaxContains = utils.CopyIntP(origin.MaxContains)
	dst.MinContains = utils.CopyIntP(origin.MinContains)
	dst.MaxProperties = utils.CopyIntP(origin.MaxProperties)
	dst.MinProperties = utils.CopyIntP(origin.MinProperties)
	dst.Required = utils.CopyStringArray(origin.Required)
	dst.DependentRequired = utils.CopyMapStringArray(origin.DependentRequired)

	dst.Format = origin.Format

	dst.ContentEncoding = origin.ContentEncoding
	dst.ContentMediaType = origin.ContentMediaType
	dst.ContentSchema = DeepCopy(origin.ContentSchema)

	dst.Title = origin.Title
	dst.Description = origin.Description
	dst.Default = utils.CopyAnyP(origin.Default)
	dst.Deprecated = utils.CopyBoolP(origin.Deprecated)
	dst.ReadOnly = utils.CopyBoolP(origin.ReadOnly)
	dst.WriteOnly = utils.CopyBoolP(origin.WriteOnly)

	dst.Extras = utils.CopyMapAny(origin.Extras)
	return dst
}

func DeepCopyArray(schemas []*Schema) []*Schema {
	if schemas == nil {
		return nil
	}

	newSchemas := make([]*Schema, len(schemas))
	for i, schema := range schemas {
		newSchemas[i] = DeepCopy(schema)
	}
	return newSchemas
}

func DeepCopyMap(schemas SchemaMap) SchemaMap {
	if schemas == nil {
		return nil
	}

	newSchemas := NewOrderedSchemaMap()
	for _, key := range schemas.Keys() {
		schema, _ := schemas.Get(key)
		newSchemas.Set(key, DeepCopy(schema))
	}
	return newSchemas
}

// Boolean Schema is special type of schema
var (
	// TrueSchema equals to `{}` or `true`
	TrueSchema = *NewBooleanSchema(true)
	// FalseSchema equals to `{"not": {}}` or `false`
	FalseSchema = *NewBooleanSchema(false)
)

func NewBooleanSchema(b bool) *Schema {
	return &Schema{boolean: &b}
}
