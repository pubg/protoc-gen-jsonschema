# TestCases
proto-oneof-single oneof가 잘 전환되는지 테스트
proto-oneof-multiple oneof가 잘 전환되는지 테스트
proto-number int32, int64등이 number로 잘 전환되는지 테스트
proto-string string가 string으로 잘 전환되는지 테스트, bytes가 base64 string으로 잘 전환되는지 테스트
proto-repeated-scala repeated가 array로 잘 전환되는지 테스트
proto-repeated-wellknown 
proto-repeated-message-and-map
proto-optional optional이 있다면 required에서 빠지는지 테스트 
proto-message message가 object로 잘 전환되는지 테스트
proto-enum enum이 enum으로 잘 전환되는지 테스트
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

