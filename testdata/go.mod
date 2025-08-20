module github.com/pubg/protoc-gen-jsonschema/testdata

go 1.22

toolchain go1.24.3

require (
	github.com/lyft/protoc-gen-star/v2 v2.0.4
	github.com/pubg/protoc-gen-jsonschema v0.0.0
	github.com/samber/lo v1.51.0
	github.com/stretchr/testify v1.8.2
	k8s.io/apimachinery v0.28.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/afero v1.3.3 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/text v0.22.0 // indirect
	google.golang.org/protobuf v1.36.7 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	sigs.k8s.io/json v0.0.0-20221116044647-bc3834ca7abd // indirect
	sigs.k8s.io/yaml v1.6.0 // indirect
)

replace github.com/pubg/protoc-gen-jsonschema => ../
