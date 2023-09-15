package modules

import (
	"fmt"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/jsonschema"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
	"google.golang.org/protobuf/encoding/protojson"
)

type Module struct {
	*pgs.ModuleBase
	pluginOptions *proto.PluginOptions
}

func NewModule() *Module {
	return &Module{ModuleBase: &pgs.ModuleBase{}}
}

func (m *Module) Name() string {
	return "JsonSchemaModule"
}

func (m *Module) InitContext(c pgs.BuildContext) {
	m.ModuleBase.InitContext(c)
	m.pluginOptions = proto.GetPluginOptions(c.Parameters())

	m.Debugf("pluginOptions: %v", protojson.MarshalOptions{EmitUnpopulated: true}.Format(m.pluginOptions))
}

func (m *Module) Execute(targets map[string]pgs.File, packages map[string]pgs.Package) []pgs.Artifact {
	// Phase: IntermediateSchemaGenerate
	visitor := NewVisitor(m)
	for _, pkg := range packages {
		m.CheckErr(pgs.Walk(visitor, pkg), fmt.Sprintf("failed to walk package %s", pkg.ProtoName().String()))
	}
	m.Debugf("# of IntermediateSchemas: %d", len(visitor.registry.GetKeys()))

	// Phase: TargetSchemaGenerate
	var optimizer Optimizer = NewOptimizerImpl(m.pluginOptions)
	var generator TargetSchemaGenerator = NewMultiDraftGenerator(m.pluginOptions)
	var serializer Serializer = NewSerializerImpl(m.pluginOptions)
	for _, file := range targets {
		m.Push("TargetSchemaGenerate")
		m.Push(file.Name().String())
		m.Debugf("FileOptions: %v", protojson.MarshalOptions{EmitUnpopulated: true}.Format(proto.GetFileOptions(file)))

		copiedRegistry := jsonschema.DeepCopyRegistry(visitor.registry)

		optimizer.Optimize(copiedRegistry, file)
		m.Debugf("# of Schemas After Optimized : %d", len(copiedRegistry.GetKeys()))

		rootSchema := generator.Generate(copiedRegistry, file)
		if rootSchema == nil {
			m.Logf("Cannot find matched entrypointMessage in %s", file.Name().String())
			continue
		}

		content, err := serializer.Serialize(rootSchema, file)
		m.CheckErr(err, fmt.Sprintf("failed to serialize file %s", file.Name().String()))
		fileName := serializer.ToFileName(file)
		m.Debugf("GeneratedFileName: %s", fileName)

		m.AddGeneratorFile(fileName, string(content))
	}

	return m.Artifacts()
}
