# protoc-gen-jsonschema

목표:

1. protobuf 파일을 jsonschema로 변환한다
2. 전체 json schema 스펙을 지원하지 않는다. (protobuf의 스펙만 json-schema로 변환한다

# 기능 대응

## protobuf field

| protobuf label | jsonschema    keyword |
|----------------|-----------------------|
| required       | required 포함           |
| optional       | required 미포함          |
| repeated       | array instance로 전환    |

| protobuf options       | jsonschema keyword        |
|------------------------|---------------------------|
| bool nullable          | oneof {type: null}        |
| stringOptions string   | -                         |
| numericOptions numeric | -                         |
| arrayOptions array     | -, only if has repeated label |
| rawOptions message     | -                         |



## protobuf message

| protobuf options          | jsonschema keyword   |
|---------------------------|----------------------|
| bool additionalProperties | additionalProperties |

## protobuf message, field
| protobuf options | jsonschema keyword |
|------------------|--------------------|
| string description | description        |
| string $comment   | $comment           |
| any default      | default            |

