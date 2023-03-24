# protoc-gen-dynamodb

[![Checks](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/checks.yaml/badge.svg)](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/checks.yaml)
[![Test](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/test.yaml/badge.svg)](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/test.yaml)

Use Protobuf to define DynamoDB item encoding using Go (golang).

## features

- Uses sdk v2
- Unit and e2e testing
- Type-safe expression path building
- Generate table definitions
- use official 'attributevalue'
- Wide(r) range of types support: everything in the canonical json table
  - Including maps with all basic types, including bool as keys
- Allow messages external to the package to be usable as field messages without problem
- Uses field position numbers by default instead of names (since they are supposed to be stable)
- Support well-knowns, but are generated to maps with strings for their fields, instead of field numbers
  - Document "Any" format in particular: "Value" stored always stored as binary
  - Document "FieldMask" format: "StringSet"
  - Structpb.Value is formatted in dynamodb
- Does no logic to support formatting pk/sk, instead supports the use code to do this
- Support of embedding fields as json
