package modules

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/invopop/jsonschema"
	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
	"gopkg.in/yaml.v3"
)

type Serializer interface {
	Serialize(schema *jsonschema.Schema, pluginOptions *proto.PluginOptions, file pgs.File) (string, []byte, error)
}

var _ Serializer = (*SerializerImpl)(nil)

type SerializerImpl struct {
}

func NewSerializerImpl() *SerializerImpl {
	return &SerializerImpl{}
}

func (s *SerializerImpl) Serialize(schema *jsonschema.Schema, pluginOptions *proto.PluginOptions, file pgs.File) (string, []byte, error) {
	fileOptions := proto.GetFileOptions(file)
	content, err := s.serializeToBytes(schema, pluginOptions, fileOptions)
	if err != nil {
		return "", nil, err
	}

	outputFileSuffix := proto.GetOutputFileSuffix(pluginOptions, fileOptions)
	fileName := file.InputPath().SetExt(outputFileSuffix).String()
	return fileName, content, nil
}

func (s *SerializerImpl) serializeToBytes(schema *jsonschema.Schema, pluginOptions *proto.PluginOptions, fileOptions *proto.FileOptions) ([]byte, error) {
	outputFileSuffix := proto.GetOutputFileSuffix(pluginOptions, fileOptions)
	if strings.HasSuffix(outputFileSuffix, ".json") {
		if proto.GetPrettyJsonOutput(pluginOptions, fileOptions) {
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
