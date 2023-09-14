#!/bin/bash

set -eux

cd $(dirname $0)

go build -o protoc-gen-jsonschema ../main.go

export DEBUG=true

protoc \
  --plugin=protoc-gen-jsonschema=./protoc-gen-jsonschema \
  --jsonschema_out=./ \
  -I ./ \
  -I ../proto \
  ./example.proto

rm protoc-gen-jsonschema
