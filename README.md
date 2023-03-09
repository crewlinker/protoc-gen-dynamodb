# protoc-gen-dynamodb

Use Protobuf to define DynamoDB item encoding using Go (golang).

## features

- Uses sdk v2
- Unit and e2e testing
- Generate table definitions
- Use protobuf field numbers
- use official 'attributevalue' with ability to customize its behaviour
- Wide(r) range of types support: everything in the canonical json table
  - Including maps with all basic types, including bool as keys
- No external dependencies of the generated code except the aws SDK
- Allow messages external to the package to be usable as field messages without problem
- Support well-knowns, but are generated to maps with strings for their fields, instead of field numbers
  - Document "Any" format in particular: "Value" stored always stored as binary
  - Document "FieldMask" format: "StringSet"
  - Structpb.Value is formatted in dynamodb

## Ideas

- Define the sk/pk on the message, and set the pk/sk member of the resulting item, error when sk/pk is not set?
- Consider using the field position numbers instead of names (since they are supposed to be stable)
- Generate table definitions for use in AWS Cloudformation/CDK
- E2E testing with LocalDynamodb docker container
- Fuzz testing with complicated protobuf message
- Generate methods for just generating "Key" attribute maps
- Generate methods for creating put/get items for transactions
- Generate methods for handling dynamodb stream Lambda events, use: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#FromDynamoDBStreamsMap
- Allow nested messages to be stored as protojson/protobinary instead of nested maps
- Similar to: https://github.com/GoogleCloudPlatform/protoc-gen-bq-schema

## Backlog

- [ ] SHOULD get a good grip on "Nullability" of fields, the documentation for "WrapperTypes" suggests that: "Wrappers use the same representation in JSON as the wrapped primitive type, except that null is allowed and preserved during data conversion and transfer.". Which suggests that
      null values are normally not preserved. So fields have "nil"value when marschalling (or the zero) value of
      a message should not have any "null" attributes.

  - Dynamodb Null documentation
    - Is a difference between "presence" and setting it to null: https://stackoverflow.com/questions/37325862/dynamodb-what-difference-does-it-make-whether-i-set-an-attribute-to-null-tr
  - A zero value of a message should not marshal to any attributes
  - Never set an explicit null, unless: Map/List value, or Wrapper types (?)
  - Also influenced by the "optional" keyword
  - Related: https://protobuf.dev/programming-guides/field_presence/
    - Go example: https://protobuf.dev/programming-guides/field_presence/#go-example
  - Related: https://protobuf.dev/programming-guides/proto3/#default
  - Json also mentions how default for null may be set: https://protobuf.dev/programming-guides/proto3/#json

    ```
    When generating JSON-encoded output from a protocol buffer, if a protobuf field has the default value and if the field doesn’t support field presence, it will be omitted from the output by default. An implementation may provide options to include fields with default values in the output.

    A proto3 field that is defined with the optional keyword supports field presence. Fields that have a value set and that support field presence always include the field value in the JSON-encoded output, even if it is the default value.
    ```

```

- [x] SHOULD match output field presence to that of json encoding
- [ ] SHOULD fix staticcheck checks on generated code such as S1002
- [ ] SHOULD test with coverage test as described here: https://go.dev/blog/integration-test-coverage
- [ ] SHOULD test with OneOf field
- [ ] SHOULD Add test that errors when unsupported map type is used
- [x] SHOULD Make sure generated error handling prints the field name and a more descriptive error
- [ ] SHOULD allow customizing the encoder/decoder options
- [ ] MUST add file header that states that the file is generated
- [ ] SHOULD make it configurable on how to handle nil/empty fields like stdlib json package
- [ ] SHOULD test that messages from external packages that DO implement the MarshalDynamoItem can be used in fields without problem
- [ ] SHOULD turn panics into errors (or add catch mechanism)
- [ ] SHOULD add support of StringSets, NumberSets, ByteSets etc
- [x] SHOULD allow customizing the dynamodb attribute name (from default "field number")
- [ ] SHOULD allow skipping certain fields for all dynamodb marshalling/unmarshalling.
  - make sure that logic applies to both marshalling, and unmarshalling (new test case)
- [x] SHOULD error when marshalling map of messages and the key is an empty string (not allowed)
- [ ] COULD improve usability of FieldMask encoding, instead of slice of strings of the field names in proto definition, could/should be the dynamodb attribute names. But this probably means implement another version of the fieldmaskpb.New() function. But https://pkg.go.dev/google.golang.org/protobuf/types/known/fieldmaskpb#Intersect states that "field.number" paths are also valid
```
