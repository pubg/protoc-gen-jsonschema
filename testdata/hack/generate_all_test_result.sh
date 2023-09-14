#!/bin/bash

set -eux

cd $(dirname $0)

go build -o protoc-gen-jsonschema ../../main.go

proto_dirs=$(ls ../cases)

for proto_dir in $proto_dirs; do
  protoc \
    --plugin=protoc-gen-jsonschema=./protoc-gen-jsonschema \
    --jsonschema_out=../cases/$proto_dir \
    --jsonschema_opt=pretty_json_output=true \
    -I ../../proto \
    -I ../cases/$proto_dir \
    ../cases/$proto_dir/*.proto
done

rm protoc-gen-jsonschema
