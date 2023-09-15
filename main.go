package main

import (
	"fmt"
	"os"

	pgs "github.com/lyft/protoc-gen-star/v2"
	"github.com/pubg/protoc-gen-jsonschema/pkg/modules"
	"google.golang.org/protobuf/types/pluginpb"
)

var version string = "develop"

func main() {
	if len(os.Args) == 2 {
		if os.Args[1] == "--version" {
			fmt.Println(version)
		} else if os.Args[1] == "--help" {
			fmt.Println("USAGE:")
			fmt.Println("  protoc-gen-jsonschema --version  : print version")
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
