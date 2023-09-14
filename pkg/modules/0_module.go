package modules

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

type JsonSchemaModule struct {
	*pgs.ModuleBase
	pluginOptions *proto.PluginOptions
}

func NewJsonSchemaModule() *JsonSchemaModule {
	return &JsonSchemaModule{ModuleBase: &pgs.ModuleBase{}}
}

func (m *JsonSchemaModule) Name() string {
	return "JsonSchemaModule"
}

func (m *JsonSchemaModule) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.pluginOptions = proto.GetPluginOptions(c.Parameters())
}

func (m *JsonSchemaModule) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	visitor := NewJsonSchemaVisitor(m)
	for _, pkg := range packages {
		m.CheckErr(pgs.Walk(visitor, pkg), fmt.Sprintf("failed to walk package %s", pkg.ProtoName().String()))
	}

	transformer := NewDraft202012Transformer()
	optimizer := NewOptimizerImpl()
	serializer := NewSerializerImpl()
	for _, file := range targets {
		rootSchema := transformer.Transform(visitor.registry, file, m.pluginOptions)
		if rootSchema == nil {
			m.Logf("Cannot find matched entrypointMessage in %s", file.Name().String())
			continue
		}

		rootSchema = optimizer.Optimize(rootSchema)

		fileName, content, err := serializer.Serialize(rootSchema, m.pluginOptions, file)
		m.CheckErr(err, fmt.Sprintf("failed to serialize file %s", file.Name().String()))

		m.AddGeneratorFile(fileName, string(content))
	}

	return m.Artifacts()
}
