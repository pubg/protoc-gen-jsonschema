# protoc-gen-jsonschema

`protoc-gen-jsonschema` is a plugin that converts `protocol buffer` formats to `json-schema`. It doesn't aim to support the entire json-schema specification, but rather focuses solely on converting the protobuf specification. Support json-schema `draft-04, draft-06, draft-07, draft-2019-09, draft-2020-12` versions.

# Installation

If you have go runtime, you can `go install` it.
```
go install github.com/pubg/protoc-gen-jsonschema
```

Alternatively, you can download pre-built binary from [Github Release](https://github.com/pubg/protoc-gen-jsonschema/releases).

# Usage

Refer to the [Plugin Options](#plguin-options) section below for various options available for this plugin.

#### I'm not sure which options to use
This plugin provides default options that are ready to use. For testing or generating a basic json-schema file, the following command is sufficient without extra options.
```
protoc --jsonschema_out=. *.proto
```

#### Generate with yaml format
```
protoc --jsonschema_out=. --jsonschema_opt=output_file_suffix=.yaml *.proto
```

#### Shrink bytes for transfer over network.
```
protoc --jsonschema_out=. --jsonschema_opt=pretty_json_output=false
```

#### I'd like to comply with the protobuf JSON mapping standard
By default, this plugin does not comply to the Protobuf standard because most plugins and other JSON libraries do not address integers larger than a 53-bit value. To ensure greater compatibility with other libraries, this plugin converts int64 values to integers instead of strings. However, to comply with the Protobuf standard, int64 values should be converted to strings. If you wish to comply to the Protocol Buffer standard, the option below will assist you.
```
protoc --jsonschema_out=. --jsonschema_opt=int64_as_string=true
```

# Options

### Plguin Options

### entrypoint_message
```
entrypoint_message is used which message should be entrypoint object of schema.

default: null or empty
example:
    - --jsonschema_opt=entrypoint_message=MyMessage
```

### output_file_suffix
```
output_file_suffix is used to determine output file name suffix.
Values should end with '.json' or '.yaml' or '.yml'.

default: .schema.json
example:
    - --jsonschema_opt=output_file_suffix=.schema.json
    - --jsonschema_opt=output_file_suffix=.schema.yaml
```

### pretty_json_output
```
pretty_json_output is used to determine output json should be pretty printed.
This option is only used when output_file_suffix is '.json'.

default: true
example:
    - --jsonschema_opt=pretty_json_output=true
    - --jsonschema_opt=pretty_json_output=false
```

### draft
```
draft is used to determine which draft version should be used.
The value should be one of Draft04, Draft05, Draft06, Draft07, Draft201909, Draft202012.

default: Draft202012
example:
    - --jsonschema_opt=draft=Draft202012
```

### mandatory_nullable
```
mandatory_nullable determines whether this plugin should treat optional field as nullable.
Many programming languages do not differentiate between undefined and null.
However, scripting languages like JavaScript and TypeScript can distinguish between them.
By default, optional field is treated as nullable and undefined.

default: false
example:
    - --jsonschema_opt=mandatory_nullable=true
    - --jsonschema_opt=mandatory_nullable=false
```

### int64_as_string
```
int64_as_string determines whether int64 field treat as string.
Depends on Javascript specification, The JS stores integer to only 53bits.
So, if you want to use int64 field in JS, you should use string type.
References:

default: false
example:
    - --jsonschema_opt=int64_as_string=true
    - --jsonschema_opt=int64_as_string=false
```

### Protobuf Options

Below tables are not auto-generated features.
Check out the [Options](./options.md) file to see auto-generated options.

| protobuf label | jsonschema                           |
|----------------|--------------------------------------|
| required       | $.required.append(field)             |
| optional       | oneof {type: null, $original_schema} |
| repeated       | type: array, items: $original_schema |

| WellKnown Types                                 |
|-------------------------------------------------|
| k8s.io.apimachinery.pkg.util.intstr.IntOrString |
| k8s.io.api.core.v1.Volume                       |
| k8s.io.api.core.v1.SecretProjection             |
| k8s.io.api.core.v1.ConfigMapVolumeSource        |
| k8s.io.api.core.v1.ConfigMapProjection          |
| k8s.io.api.core.v1.ConfigMapKeySelector         |
| k8s.io.api.core.v1.SecretKeySelector            |
| k8s.io.api.core.v1.ConfigMapEnvSource           |
| k8s.io.api.core.v1.SecretEnvSource              |
| k8s.io.api.core.v1.Probe                        |
| k8s.io.api.core.v1.EphemeralContainer           |
| google.protobuf.Timestamp                       |
| google.protobuf.Duration                        |
| google.protobuf.Any                             |
| google.protobuf.NullValue                       |

If you'd like contribute well-known types, Please check [generator](./pkg/modules/1_middleend_generator.go) file.

# Output Examples
You can find basic example cases in the `./examples` directory and more complex cases in `./testdata/cases.`
