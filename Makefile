options:
	protoc \
	--go_out=./pkg/proto \
	--go_opt=paths=source_relative \
	--doc_out=. \
	--doc_opt=markdown,options.md \
	./*.proto

test-examples:
	go build -o protoc-gen-jsonschema main.go
	protoc \
	--plugin=protoc-gen-jsonschema=./protoc-gen-jsonschema \
	--jsonschema_out=./ \
	--jsonschema_opt=additional_properties=true \
	--jsonschema_opt=draft=Draft04 \
	--jsonschema_opt=output_file_suffix=.schema.json \
	-I ./ \
	-I ./examples \
	./examples/example.proto
	rm protoc-gen-jsonschema


generate-examples:
	go build -o protoc-gen-jsonschema main.go
	protoc \
	--plugin=protoc-gen-jsonschema=./protoc-gen-jsonschema \
	--jsonschema_out=./ \
	--jsonschema_opt=draft=Draft04 \
	--jsonschema_opt=output_file_suffix=.draft-04.json \
	--plugin=protoc-gen-jsonschema06=./protoc-gen-jsonschema \
	--jsonschema06_out=./ \
	--jsonschema06_opt=draft=Draft06 \
	--jsonschema06_opt=output_file_suffix=.draft-06.json \
	--plugin=protoc-gen-jsonschema07=./protoc-gen-jsonschema \
	--jsonschema07_out=./ \
	--jsonschema07_opt=draft=Draft07 \
	--jsonschema07_opt=output_file_suffix=.draft-07.json \
	--plugin=protoc-gen-jsonschema19=./protoc-gen-jsonschema \
	--jsonschema19_out=./ \
	--jsonschema19_opt=draft=Draft201909 \
	--jsonschema19_opt=output_file_suffix=.draft-2019-09.json \
	--plugin=protoc-gen-jsonschema20=./protoc-gen-jsonschema \
	--jsonschema20_out=./ \
	--jsonschema20_opt=draft=Draft202012 \
	--jsonschema20_opt=output_file_suffix=.draft-2020-12.json \
	--plugin=protoc-gen-jsonschemayaml=./protoc-gen-jsonschema \
	--jsonschemayaml_out=./ \
	--jsonschemayaml_opt=output_file_suffix=.yaml \
	-I ./ \
	-I ./examples \
	./examples/example.proto
	rm protoc-gen-jsonschema

dump-testdata:
	protoc \
    --debug_out=./ \
    --debug_opt=dump_binary=true \
    --debug_opt=file_binary=request.pb.bin \
    -I ./ \
    -I ./testdata/cases/jsonschema-object \
    ./testdata/cases/jsonschema-object/test.proto
