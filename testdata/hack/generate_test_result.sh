#!/bin/bash

set -eux

proto_dir=$1

cd $(dirname $0)

cd $(git rev-parse --show-toplevel)

go build -o ./testdata/hack/protoc-gen-jsonschema main.go

function generateProtoFile() {
  proto_dir=$1

  command="protoc \
  --plugin=protoc-gen-jsonschema=./testdata/hack/protoc-gen-jsonschema \
  --jsonschema_out=./testdata/cases/$proto_dir "  # 초기 명령어 추가

  # test.yaml 파일의 존재 여부 확인
  if [[ -f "./testdata/cases/$proto_dir/test.yaml" ]]; then
    # test.yaml에서 inputParameters 부분을 읽고 JSON 형식으로 파싱
    parameters=$(yq eval '.inputParameters' ./testdata/cases/$proto_dir/test.yaml)

    # inputParameters가 없거나 비어 있지 않은지 확인
    if [[ "$parameters" != "null" && "$parameters" != "{}" ]]; then
      # key-value 쌍을 반복하면서 --jsonschema_opt=<key>=<value> 형태로 문자열 결합
      while IFS=":" read -r key value; do
          key=$(echo $key | tr -d ' "')
          value=$(echo $value | tr -d ' "')
          command="${command} --jsonschema_opt=${key}=${value}"
      done <<< "$parameters"
    fi
  fi

  command="$command \
  -I ./ \
  -I testdata/cases/$proto_dir \
  testdata/cases/$proto_dir/*.proto
  "

  $command
}

generateProtoFile $proto_dir

rm ./testdata/hack/protoc-gen-jsonschema
