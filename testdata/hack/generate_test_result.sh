#!/bin/bash

set -eux

proto_dir=$1

cd $(dirname $0)

go build -o protoc-gen-venus ../../main.go

protoc \
  --plugin=protoc-gen-venus=./protoc-gen-venus \
  --venus_out=../cases/$proto_dir \
  -I ../../proto \
  -I ../cases/$proto_dir \
  ../cases/$proto_dir/*.proto

rm protoc-gen-venus
