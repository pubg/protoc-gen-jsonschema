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

	builder := &strings.Builder{}
	for _, comment := range name.SourceCodeInfo().LeadingDetachedComments() {
		builder.WriteString(comment)
	}
	builder.WriteString(name.SourceCodeInfo().LeadingComments())
	builder.WriteString(name.SourceCodeInfo().TrailingComments())
	comment := strings.TrimSpace(strings.Trim(builder.String(), "\n"))

	if comment != "" {
		return comment
	}
	return ""
}

func GetOutputFileSuffix(pluginOptions *PluginOptions, fileOptions *FileOptions) string {
	if fileOptions != nil && fileOptions.GetOutputFileSuffix() != "" {
		return fileOptions.GetOutputFileSuffix()
	} else {
		return pluginOptions.GetOutputFileSuffix()
	}
}

func GetPrettyJsonOutput(pluginOptions *PluginOptions, fileOptions *FileOptions) bool {
	if fileOptions != nil && fileOptions.PrettyJsonOutput != nil {
		return fileOptions.GetPrettyJsonOutput()
	} else {
		return pluginOptions.GetPrettyJsonOutput()
	}
}

func GetEntrypointMessage(pluginOptions *PluginOptions, fileOptions *FileOptions) string {
	if fileOptions != nil && fileOptions.GetEntrypointMessage() != "" {
		return fileOptions.GetEntrypointMessage()
	} else {
		return pluginOptions.GetEntrypointMessage()
	}
}
