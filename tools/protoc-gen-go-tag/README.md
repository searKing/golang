[![GoDoc](https://pkg.go.dev/badge/github.com/searKing/golang/tools/protoc-gen-go-tag)](https://pkg.go.dev/github.com/searKing/golang/tools/protoc-gen-go-tag)
[![Report card](https://goreportcard.com/badge/github.com/searKing/golang/tools/protoc-gen-go-tag)](https://goreportcard.com/report/github.com/searKing/golang/tools/protoc-gen-go-tag)

# protoc-gen-go-tag

protoc-gen-go-tag generates `*.pb.go` as `protoc-gen-go` does, with custom struct tags of fields set by the
field option `field_tag`, or by `@gotags` comments, which require nothing to import.

For example, given [examples/example.proto](examples/example.proto),

```proto
syntax = "proto3";

package pb;

import "searking/tag/v1/tag.proto";
option go_package = "github.com/searKing/golang/tools/protoc-gen-go-tag/examples;pb";

message Http{
  string protocol = 1[(searking.tag.v1.field_tag) = {struct_tag: 'validate:"oneof=http https"'}];
  // version_default is the same as version_update
  string version_default = 2[json_name = "VersionDefault", (searking.tag.v1.field_tag) = {struct_tag: 'validate:"gte=0,lte=130" json:"version_with_default"'}];
  string version_update = 3[json_name = "VersionUpdate", (searking.tag.v1.field_tag) = {struct_tag: 'validate:"gte=0,lte=130" json:"version_with_update"', update_strategy: UPDATE_STRATEGY_UPDATE}];
  string version_replace = 4[json_name = "VersionReplace", (searking.tag.v1.field_tag) = {struct_tag: 'validate:"gte=0,lte=130" json:"version_with_replace"', update_strategy: UPDATE_STRATEGY_REPLACE}];

  Url url = 5[json_name = "Url", (searking.tag.v1.field_tag) = {struct_tag: "json:\"url_tag,omitempty\""}];
  message Url {
    string scheme = 1[json_name = "Scheme", (searking.tag.v1.field_tag) = {struct_tag: 'validate:"required" json:"schema_tag,omitempty"'}];
  }

  string test_order_tag_form_v8n = 10[(searking.tag.v1.field_tag) = {struct_tag: 'form:"currency_uso,omitempty" v8n:"required"'}];
  string test_order_tag_v8n_form = 11[(searking.tag.v1.field_tag) = {struct_tag: 'v8n:"required" form:"currency_uso,omitempty"'}];

  oneof test_oneof {
    string name = 6[json_name = "Name", (searking.tag.v1.field_tag) = {struct_tag: 'validate:"required" json:"name,omitempty"'}];
    string name_no_json_name = 7[(searking.tag.v1.field_tag) = {struct_tag: 'json:"name_no_json_name,omitempty"'}];
  }
}```

running this command, with [searking/tag/v1/tag.proto](proto/searking/tag/v1/tag.proto) on the `protoc` include path,
see [Get tag.proto](#get-tagproto),

```bash
protoc -I . -I third_party --go-tag_out=paths=source_relative:. example.proto
```

generates [examples/example.pb.go](examples/example.pb.go), containing

```go
type Http struct {
	Protocol string `protobuf:"bytes,1,opt,name=protocol,proto3" json:"protocol,omitempty" validate:"oneof=http https"`
	// version_default is the same as version_update
	VersionDefault      string    `protobuf:"bytes,2,opt,name=version_default,json=VersionDefault,proto3" json:"version_with_default,omitempty" validate:"gte=0,lte=130"`
	VersionUpdate       string    `protobuf:"bytes,3,opt,name=version_update,json=VersionUpdate,proto3" json:"version_with_update,omitempty" validate:"gte=0,lte=130"`
	VersionReplace      string    `protobuf:"bytes,4,opt,name=version_replace,json=VersionReplace,proto3" json:"version_with_replace" validate:"gte=0,lte=130"`
	Url                 *Http_Url `protobuf:"bytes,5,opt,name=url,json=Url,proto3" json:"url_tag,omitempty"`
	TestOrderTagFormV8N string    `protobuf:"bytes,10,opt,name=test_order_tag_form_v8n,json=testOrderTagFormV8n,proto3" json:"test_order_tag_form_v8n,omitempty" form:"currency_uso,omitempty" v8n:"required"`
	TestOrderTagV8NForm string    `protobuf:"bytes,11,opt,name=test_order_tag_v8n_form,json=testOrderTagV8nForm,proto3" json:"test_order_tag_v8n_form,omitempty" v8n:"required" form:"currency_uso,omitempty"`
	TestOneof           isHttp_TestOneof `protobuf_oneof:"test_oneof"`
	// ...
}

type Http_Url struct {
	Scheme string `protobuf:"bytes,1,opt,name=scheme,json=Scheme,proto3" json:"schema_tag,omitempty" validate:"required"`
	// ...
}

type Http_Name struct {
	Name string `protobuf:"bytes,6,opt,name=name,json=Name,proto3,oneof" json:"name,omitempty" validate:"required"`
}

type Http_NameNoJsonName struct {
	NameNoJsonName string `protobuf:"bytes,7,opt,name=name_no_json_name,json=nameNoJsonName,proto3,oneof" json:"name_no_json_name,omitempty"`
}
```

## Comments

Struct tags can be set by comment directives in leading or trailing comments of fields, without importing
`tag.proto`, as [protoc-go-inject-tag](https://github.com/favadi/protoc-go-inject-tag) does:

```proto
message Comment {
  // @gotags: validate:"required" json:"name_tag"
  string name = 1;
  string email = 2; // @gotags: validate:"email"
  // @gotags-replace: json:"age"
  int32 age = 3;
}
```

+ `// @gotags: ...` updates struct tags as `update_strategy: UPDATE_STRATEGY_UPDATE` does.
+ `// @gotags-replace: ...` replaces struct tags as `update_strategy: UPDATE_STRATEGY_REPLACE` does.
+ A directive takes a whole line of comments, and is kept in comments of the generated Go code.
+ A field may have many directives, and the option `field_tag` too; they are applied in order of the option
  `field_tag`, directives in leading comments, then directives in trailing comments.

## Get tag.proto

`tag.proto` is imported by `searking/tag/v1/tag.proto`, and declares the field option `searking.tag.v1.field_tag` of
the extension number 1391, registered in the
[Protobuf Global Extension Registry](https://github.com/protocolbuffers/protobuf/blob/main/docs/options.md).

### Export by protoc-gen-go-tag

`protoc-gen-go-tag` embeds `tag.proto` of its version, and exports it by the import path:

```sh
protoc-gen-go-tag -export_proto third_party
# third_party/searking/tag/v1/tag.proto
# third_party/github.com/searKing/golang/tools/protoc-gen-go-tag/tag/tag.proto, deprecated
protoc -I . -I third_party --go-tag_out=paths=source_relative:. example.proto
```

### Depend on BSR by buf

`tag.proto` is pushed to the [Buf Schema Registry](https://buf.build/searking-proto/protoc-gen-go-tag) from [proto](proto):

```yaml
# buf.yaml
version: v2
deps:
  - buf.build/searking-proto/protoc-gen-go-tag
```

```yaml
# buf.gen.yaml
version: v2
plugins:
  - local: protoc-gen-go-tag
    out: .
    opt: paths=source_relative
```

```sh
buf dep update
buf generate
```

### Declare field_tag in your proto

Declare `field_tag` of a message wire compatible with `searking.tag.v1.FieldTag` in your proto, instead of importing
`tag.proto`, and set its extension number by the parameter `field_tag_number`:

```proto
syntax = "proto3";

package pb;

import "google/protobuf/descriptor.proto";

extend google.protobuf.FieldOptions {
  FieldTag field_tag = 60001;
}

message FieldTag {
  enum UpdateStrategy {
    UPDATE = 0;
    REPLACE = 1;
  }
  string struct_tag = 1;
  UpdateStrategy update_strategy = 2;
}

message Http {
  string version = 1 [(field_tag) = {struct_tag: 'validate:"gte=0,lte=130"'}];
}
```

```sh
protoc -I . --go-tag_out=paths=source_relative,field_tag_number=60001:. http.proto
```

`field_tag_number` may be repeated, and is 65000 of the deprecated `google.protobuf.field_tag` by default, so set it
to another number if 65000 conflicts with other extensions. `searking.tag.v1.field_tag` is read always.

### Migrate from google.protobuf.field_tag

`google.protobuf.field_tag` of `github.com/searKing/golang/tools/protoc-gen-go-tag/tag/tag.proto` is deprecated, for
it is declared in the package `google.protobuf` of protobuf itself, with the extension number 65000 in the range for
internal use of organizations, which may conflict with others. It is still read by default, migrate by:

+ Replace `import "github.com/searKing/golang/tools/protoc-gen-go-tag/tag/tag.proto";` with
  `import "searking/tag/v1/tag.proto";`.
+ Replace `(google.protobuf.field_tag)` with `(searking.tag.v1.field_tag)`.
+ Replace `update_strategy: update` and `update_strategy: replace` with `update_strategy: UPDATE_STRATEGY_UPDATE` and
  `update_strategy: UPDATE_STRATEGY_REPLACE`.

## Struct tags

+ `update_strategy: UPDATE_STRATEGY_UPDATE`, the default, updates struct tags generated by `protoc-gen-go` with `struct_tag`:
  names are replaced, and options are added, such as `json:"version_with_update,omitempty"`.
+ `update_strategy: UPDATE_STRATEGY_REPLACE` replaces struct tags generated by `protoc-gen-go` with `struct_tag`,
  such as `json:"version_with_replace"`.
+ Struct tags are ordered as `protobuf`, `json`, then other tags in the order written in `struct_tag`.
+ `struct_tag` must not set `protobuf` or `protobuf_oneof`, which the protobuf runtime relies on.
+ Struct tags are set on fields of Go structs named by `protoc-gen-go`, including renamed ones such as `Reset_` of
  the field `reset`, wrapper structs of oneof fields, and fields of proto3 `optional`.
+ `field_tag` options are applied in order of `searking.tag.v1.field_tag`, then numbers of `field_tag_number`.

## Download/Install

```sh
go install github.com/searKing/golang/tools/protoc-gen-go-tag@latest
```

The Go code is generated by `protoc-gen-go` of the version of `google.golang.org/protobuf` that `protoc-gen-go-tag`
requires.

### protoc-gen-go-tag: program not found or is not executable

The `protoc` compiler expects to find plugins named `proto-gen-xxx` on the execution `$PATH`. So first:

```sh
export PATH=${PATH}:$(go env GOPATH)/bin
```

### xxx.proto: Tried to write the same file twice.

+ The `protoc-gen-go-tag` rewrites `*.pb.go` by `proto-gen-go` already. So use `protoc-gen-go-tag` instead
  of `protoc-gen-go`
+ Do not run `--go_out` with `--go-tag_out` together, for `protoc-gen-go-tag` generates `*.pb.go` as `protoc-gen-go` does
+ Basically the magical incantation (apart from includes) is the `--go-tag_out`. That triggers the
  `protoc-gen-go-tag` plugin to generate `*.pb.go`. That's it :)

### If you use `protoc-gen-go` of `github.com/golang/protobuf`

+ [github.com/golang/protobuf](https://pkg.go.dev/github.com/golang/protobuf) has been superseded by
  the [google.golang.org/protobuf](https://pkg.go.dev/google.golang.org/protobuf) module
+ install `protoc-gen-go-tag`
  by `go install github.com/searKing/golang/tools/protoc-gen-go-tag@e27209892e58d3df53c587f1feefb6290af17d85`
+ `protoc -I . --go_out=paths=source_relative:. --go-tag_out=paths=source_relative:. *.proto`
+ `*.pb.go` generated by `protoc-gen-go` will be rewritten by `protoc-gen-go-tag`

## Development

Requires `protoc`, and `google/protobuf/descriptor.proto` in the include directory installed with `protoc`, or in
`PROTOC_INCLUDE`.

```sh
# Regenerate examples/example.pb.go, and testdata/protos.binpb, the descriptor set of protos in testdata/proto.
go generate
# Update golden files of testdata/golden, generated from testdata/protos.binpb.
go test -update
```

Push `searking/tag/v1/tag.proto` to the BSR module `buf.build/searking-proto/protoc-gen-go-tag`, after `buf registry login`,
if it is changed:

```sh
buf build proto
# main is the default label, which deps without labels resolve to.
buf push proto --label main --label v1.2.147
```

## Inspiring projects

* [protobuf](https://pkg.go.dev/google.golang.org/protobuf)
