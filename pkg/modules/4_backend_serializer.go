package modules

import (
	"encoding/json"
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
	"sigs.k8s.io/yaml"
)

type BackendSerializer interface {
	Serialize(schema any, file pgs.File) ([]byte, error)
	ToFileName(file pgs.File) string
}

var _ BackendSerializer = (*SerializerImpl)(nil)

type SerializerImpl struct {
	pluginOptions *proto.PluginOptions
}

func NewSerializerImpl(pluginOptions *proto.PluginOptions) *SerializerImpl {
	return &SerializerImpl{pluginOptions: pluginOptions}
}

func (s *SerializerImpl) Serialize(schema any, file pgs.File) ([]byte, error) {
	outputFileSuffix := s.ToFileName(file)

	if strings.HasSuffix(outputFileSuffix, ".json") {
		if s.pluginOptions.GetPrettyJsonOutput() {
			return json.MarshalIndent(schema, "", "  ")
		} else {
			return json.Marshal(schema)
		}
	} else if strings.HasSuffix(outputFileSuffix, ".yaml") || strings.HasSuffix(outputFileSuffix, ".yml") {
		return yaml.Marshal(schema)
	} else {
		return nil, fmt.Errorf("unsupported output file suffix: `%s`, suffix should be endsWith `.json`, `.yaml`, `.yml`", outputFileSuffix)
	}
}

func (s *SerializerImpl) ToFileName(file pgs.File) string {
	outputFileSuffix := s.pluginOptions.OutputFileSuffix
	return file.InputPath().SetExt(outputFileSuffix).String()
}
