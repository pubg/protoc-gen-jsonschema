# TestCases

This directory contains test cases for the protoc-gen-jsonschema plugin.

# Make a new test case

1. Create a new directory in the `testdata/cases` directory.
2. Create `test.proto` file and `test.yaml` file in the new directory.
3. Fill the `test.proto` file with proper content.
4. Fill the `test.yaml` file with expected test result.

### test.yaml format

```yaml
name: "proto-json-name-ignore" # test case name, should be unique
description: "test for ignore proto json_name option" # description of the test case

inputFiles: ["test.proto"] # input proto file
inputParameters: # protoc-gen-jsonschema plugin options. Map<string, any> type
  preserve_proto_field_names: true

expectResultFiles: ["test.schema.json"] # expected output schema file
expectResultIsNull: false # if true, expectResultFile is empty test will pass
expectFail: false # if true, generate schema fail test will pass

expectedBehaviorDescription: "" # not used
```

### Generate all test case results

```sh
cd hack
./generate_all_test_result.sh
```

### Generate single test case result

```sh
cd hack
./generate_test_result.sh <test case name>
```

# Test all cases

```sh
go test do_test.go
```

# Test single case

Add ad-hoc filter code to the test code
```go
46 for _, dir := range dirs {
47     if dir.Name() != "proto-json-name-ignore" {
48        continue
49    }
    ...
}
```

```sh
go test do_test.go
```

# Description of each test case

proto-oneof-single oneof가 잘 전환되는지 테스트

proto-oneof-multiple oneof가 잘 전환되는지 테스트

proto-number int32, int64등이 number로 잘 전환되는지 테스트

proto-string string가 string으로 잘 전환되는지 테스트, bytes가 base64 string으로 잘 전환되는지 테스트

proto-repeated-scala repeated가 array로 잘 전환되는지 테스트

proto-repeated-wellknown

proto-repeated-message-and-map

proto-optional optional이 있다면 required에서 빠지는지 테스트

proto-message message가 object로 잘 전환되는지 테스트

proto-circular-dependent 순환 의존성이 있는 메시지를 잘 전환하는지 테스트

proto-enum enum이 enum으로 잘 전환되는지 테스트

plugin-preserve-field-name-false json_name이 잘 적용되는지 테스트

plugin-preserve-field-name-true json_name이 무시되는지 테스트

proto-map map이 object로 잘 전환되는지 테스트, array 속성이 들어가면 안됨

proto-nested-message nested message가 잘 전환되는지 테스트

jsonschema-nullable nullable 속성이 있으면 oneof로 잘 전환되는지 테스트

jsonschema-number number validation keyword 테스트

jsonschema-string string validation keyword 테스트

jsonschema-array array validation keyword 테스트

jsonschema-object object validation keyword 테스트

[//]: # (jsonschema-enum enum validation keyword 테스트)
[//]: # (jsonschema-const const validation keyword 테스트)

plugin-entrypoint-options plugin entrypoint가 잘 작동하는지 테스트

plugin-output-options output 관련 옵션이 잘 작동하는지 테스트

plugin-draft-04 draft-04 스키마가 잘 작동하는지 테스트

plugin-draft-06 draft-06 스키마가 잘 작동하는지 테스트

plugin-draft-07 draft-07 스키마가 잘 작동하는지 테스트

plugin-draft-201909 draft-2019-09 스키마가 잘 작동하는지 테스트

plugin-draft-202012 draft-2020-12 스키마가 잘 작동하는지 테스트

wellknown-any google.protobuf.Any가 잘 전환되는지 테스트

wellknown-timestamp google.protobuf.Timestamp가 잘 전환되는지 테스트

wellknown-duration google.protobuf.Duration가 잘 전환되는지 테스트

wellknown-k8s-intorstring k8s.io.apimachinery.pkg.util.intstr.IntOrString가 잘 전환되는지 테스트

