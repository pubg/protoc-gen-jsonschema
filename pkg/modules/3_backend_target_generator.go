package modules

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema/draft_04"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema/draft_06"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema/draft_07"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema/draft_201909"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema/draft_202012"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

type BackendTargetGenerator interface {
	Generate(registry *jsonschema.Registry, entrypointMessage pgs.Message, fileOptions *proto.FileOptions) any
}

var _ BackendTargetGenerator = (*MultiDraftGenerator)(nil)

type MultiDraftGenerator struct {
	pluginOptions *proto.PluginOptions
}

func NewMultiDraftGenerator(pluginOptions *proto.PluginOptions) *MultiDraftGenerator {
	return &MultiDraftGenerator{pluginOptions: pluginOptions}
}

const draft04Version = "http://json-schema.org/draft-04/schema#"
const draft06Version = "http://json-schema.org/draft-06/schema#"
const draft07Version = "http://json-schema.org/draft-07/schema#"
const draft201909Version = "https://json-schema.org/draft/2019-09/schema"
const draft202012Version = "https://json-schema.org/draft/2020-12/schema"

func (g *MultiDraftGenerator) Generate(registry *jsonschema.Registry, entrypointMessage pgs.Message, fileOptions *proto.FileOptions) any {
	rootSchema := &jsonschema.Schema{}
	rootSchema.Title = proto.GetTitleOrEmpty(fileOptions)
	rootSchema.Description = proto.GetDescriptionOrEmpty(fileOptions)
	rootSchema.Type = "object"

	rootSchema.Ref = toRefId(entrypointMessage)

	rootSchema.Definitions = jsonschema.NewOrderedSchemaMap()
	for _, key := range registry.GetKeys() {
		rootSchema.Definitions.Set(key, registry.GetSchema(key))
	}

	switch g.pluginOptions.GetDraft() {
	case proto.Draft_Draft04:
		rootSchema.Version = draft04Version
		return draft_04.New(rootSchema)
	case proto.Draft_Draft06:
		rootSchema.Version = draft06Version
		return draft_06.New(rootSchema)
	case proto.Draft_Draft07:
		rootSchema.Version = draft07Version
		return draft_07.New(rootSchema)
	case proto.Draft_Draft201909:
		rootSchema.Version = draft201909Version
		return draft_201909.New(rootSchema)
	case proto.Draft_Draft202012:
		rootSchema.Version = draft202012Version
		return draft_202012.New(rootSchema)
	}
	return nil
}

func (g *MultiDraftGenerator) getEntrypointMessage(messages []pgs.Message, fileOptions *proto.FileOptions) pgs.Message {
	entryPointMessage := proto.GetEntrypointMessage(g.pluginOptions, fileOptions)
	if entryPointMessage == "" {
		// TODO: print debug message
		return nil
	}

	for _, message := range messages {
		if message.Name().String() == entryPointMessage {
			return message
		}
	}
	return nil
}
