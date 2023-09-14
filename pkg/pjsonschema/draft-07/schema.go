package draft_07

type Schema struct {
	// for root Schema
	Schema string `json:"$schema,omitempty"`

	// base properties
	Id          string `json:"$id,omitempty"`
	Ref         string `json:"$ref,omitempty"`
	Comment     string `json:"$comment,omitempty"`
	Description string `json:"description,omitempty"`

	Defs map[string]*Schema `json:"$defs,omitempty"`
}

type AnyInstance struct {
	Type  string    `json:"type,omitempty"`
	Enum  []*Schema `json:"enum,omitempty"`
	Const *Schema   `json:"const,omitempty"`
}

type NumericKeywords struct {
	/*
		The value of "multipleOf" MUST be a number, strictly greater than 0.

		A numeric instance is valid only if division by this keyword's value
		results in an integer.
	*/
	MultipleOf *float64 `json:"multipleOf,omitempty"`
	/*
		The value of "maximum" MUST be a number, representing an inclusive
		upper limit for a numeric instance.

		If the instance is a number, then this keyword validates only if the
		instance is less than or exactly equal to "maximum".
	*/
	Maximum *float64 `json:"maximum,omitempty"`
	/*
		The value of "exclusiveMaximum" MUST be number, representing an
		exclusive upper limit for a numeric instance.

		If the instance is a number, then the instance is valid only if it
		has a value strictly less than (not equal to) "exclusiveMaximum".
	*/
	ExclusiveMaximum *float64 `json:"exclusiveMaximum,omitempty"`
	/*
		The value of "minimum" MUST be a number, representing an inclusive
		lower limit for a numeric instance.

		If the instance is a number, then this keyword validates only if the
		instance is greater than or exactly equal to "minimum".
	*/
	Minimum *float64 `json:"minimum,omitempty"`
	/*
		The value of "exclusiveMinimum" MUST be number, representing an
		exclusive lower limit for a numeric instance.

		If the instance is a number, then the instance is valid only if it
		has a value strictly greater than (not equal to) "exclusiveMinimum".
	*/
	ExclusiveMinimum *float64 `json:"exclusiveMinimum,omitempty"`
}

type StringKeywords struct {
	MaxLength *int   `json:"maxLength,omitempty"`
	MinLength *int   `json:"minLength,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
}

type ArrayKeywords struct {
	MaxProperties        *int                   `json:"maxProperties,omitempty"`
	MinProperties        *int                   `json:"minProperties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	Properties           map[string]*Schema     `json:"properties,omitempty"`
	PatternProperties    map[string]*Schema     `json:"patternProperties,omitempty"`
	AdditionalProperties *Schema                `json:"additionalProperties,omitempty"`
	Dependencies         map[string]interface{} `json:"dependencies,omitempty"`
	PropertyNames        *Schema                `json:"propertyNames,omitempty"`
}
