package main

import (
	"fmt"
	"os"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/modules"
	"google.golang.org/protobuf/types/pluginpb"
)

var version string = "develop"

const helpMessage = `Usage protoc-gen-jsonschema:
protoc --jsonschema_out=. (OPTIONS) *.proto

This plugin converts 'protocol buffer' formats to 'json-schema'. Support json-schema 'draft-04, draft-06, draft-07, draft-2019-09, draft-2020-12' versions.

EXAMPLE: I'm not sure which options to use
protoc --jsonschema_out=. *.proto

EXAMPLE: Generate with yaml format
protoc --jsonschema_out=. --jsonschema_opt=output_file_suffix=.yaml *.proto

EXAMPLE: Shrink bytes for transfer over network
protoc --jsonschema_out=. --jsonschema_opt=pretty_json_output=false

EXAMPLE: I'd like to comply with the protobuf JSON mapping standard
protoc --jsonschema_out=. --jsonschema_opt=int64_as_string=true

EXAMPLE: I'm not satisfied with the plugin's options. I want to customize every fields
cp jsonschema.proto examples/jsonschema.proto
protoc --jsonschema_out=. --proto_path=examples examples/jsonschema.proto

Please check https://github.com/pubg/protoc-gen-jsonschema for more details.

FLAGS:
  --version  : print version
  --help     : print help
`

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "--version" {
			fmt.Println(version)
		} else if os.Args[1] == "--help" {
			fmt.Print(helpMessage)
		}
		return
	}

	feature := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.Init(
		pgs.DebugEnv("DEBUG"),
		pgs.SupportedFeatures(&feature)).
		RegisterModule(modules.NewModule()).
		Render()
}
