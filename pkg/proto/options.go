package proto

import (
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
)

func GetPluginOptions(params pgs.Parameters) *PluginOptions {
	options := &PluginOptions{}

	visibilityLevel, _ := params.UintDefault("visibility_level", 0)
	options.VisibilityLevel = uint32(visibilityLevel)

	options.EntrypointMessage = params.StrDefault("entrypoint_message", "")
	options.OutputFileSuffix = params.StrDefault("output_file_suffix", ".schema.json")
	options.PrettyJsonOutput, _ = params.BoolDefault("pretty_json_output", true)
	options.MandatoryNullable, _ = params.BoolDefault("mandatory_nullable", false)
	options.Int64AsString, _ = params.BoolDefault("int64_as_string", false)
	options.PreserveProtoFieldNames, _ = params.BoolDefault("preserve_proto_field_names", false)

	if _, ok := params["additional_properties"]; ok {
		val, _ := params.Bool("additional_properties")
		options.AdditionalProperties = &val
	}

	// Default Value
	options.Draft = Draft_Draft202012
	for draft, index := range Draft_value {
		if params.Str("draft") == draft {
			options.Draft = Draft(index)
		}
	}
	return options
}

func GetFileOptions(file pgs.File) *FileOptions {
	options := &FileOptions{}
	if ok, _ := file.Extension(E_File, options); ok {
		return options
	}
	return nil
}

func GetMessageOptions(message pgs.Message) *MessageOptions {
	options := &MessageOptions{}
	if ok, _ := message.Extension(E_Message, options); ok {
		return options
	}
	return nil
}

func GetFieldOptions(field pgs.Field) *FieldOptions {
	options := &FieldOptions{}
	if ok, _ := field.Extension(E_Field, options); ok {
		return options
	}
	return nil
}

func GetEnumOptions(enum pgs.Enum) *EnumOptions {
	options := &EnumOptions{}
	if ok, _ := enum.Extension(E_Enum, options); ok {
		return options
	}
	return nil
}

func GetEnumValueOptions(enumValue pgs.EnumValue) *EnumValueOptions {
	options := &EnumValueOptions{}
	if ok, _ := enumValue.Extension(E_EnumValue, options); ok {
		return options
	}
	return nil
}

type TitleResolver interface {
	GetTitle() string
}

func GetTitleOrEmpty(title TitleResolver) string {
	if title != nil && title.GetTitle() != "" {
		return title.GetTitle()
	}
	return ""
}

type DescriptionResolver interface {
	GetDescription() string
}

type NameWithSourceResolver interface {
	SourceCodeInfo() pgs.SourceCodeInfo
	Name() pgs.Name
}

func GetDescriptionOrEmpty(description DescriptionResolver) string {
	if description != nil && description.GetDescription() != "" {
		return description.GetDescription()
	}
	return ""
}

func GetDescriptionOrComment(name NameWithSourceResolver, description DescriptionResolver) string {
	if description != nil && description.GetDescription() != "" {
		return description.GetDescription()
	}

	srcInfo := name.SourceCodeInfo()
	if srcInfo == nil {
		return ""
	}

	builder := &strings.Builder{}
	for _, comment := range srcInfo.LeadingDetachedComments() {
		builder.WriteString(comment)
	}
	builder.WriteString(srcInfo.LeadingComments())
	builder.WriteString(srcInfo.TrailingComments())
	comment := strings.TrimSpace(strings.Trim(builder.String(), "\n"))

	if comment != "" {
		return comment
	}
	return ""
}

func GetEntrypointMessage(pluginOptions *PluginOptions, fileOptions *FileOptions) string {
	if fileOptions != nil && fileOptions.GetEntrypointMessage() != "" {
		return fileOptions.GetEntrypointMessage()
	} else {
		return pluginOptions.GetEntrypointMessage()
	}
}
