package modules

import (
	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/proto"
)

func getEntrypointFromFile(file pgs.File, pluginOptions *proto.PluginOptions) pgs.Message {
	entryPointMessage := proto.GetEntrypointMessage(pluginOptions, proto.GetFileOptions(file))
	if entryPointMessage == "" {
		return nil
	}

	for _, message := range file.Messages() {
		if message.Name().String() == entryPointMessage {
			return message
		}
	}
	return nil
}
