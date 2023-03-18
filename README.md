# protoc-gen-dynamodb

[![Checks](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/checks.yaml/badge.svg)](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/checks.yaml)
[![Test](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/test.yaml/badge.svg)](https://github.com/crewlinker/protoc-gen-dynamodb/actions/workflows/test.yaml)

Use Protobuf to define DynamoDB item encoding using Go (golang).

## features

- Uses sdk v2
- Unit and e2e testing
- Generate table definitions
- Use protobuf field numbers
- use official 'attributevalue' with ability to customize its behaviour
- Wide(r) range of types support: everything in the canonical json table
  - Including maps with all basic types, including bool as keys
- ~~No external dependencies of the generated code except the aws SDK~~
- Allow messages external to the package to be usable as field messages without problem
- Uses field position numbers by default instead of names (since they are supposed to be stable)
- Support well-knowns, but are generated to maps with strings for their fields, instead of field numbers
  - Document "Any" format in particular: "Value" stored always stored as binary
  - Document "FieldMask" format: "StringSet"
  - Structpb.Value is formatted in dynamodb
- Does no logic to support formatting pk/sk, instead supports the use code to do this

## Ideas

- Define the sk/pk on the message, and set the pk/sk member of the resulting item, error when sk/pk is not set?
- Define indexes (and maybe streams) so they can be used by aws CDK, indexes may generate extra methods
- Generate table definitions for use in AWS Cloudformation/CDK, and or generate/support create table for initializing tables on LocalDynamo
- E2E testing with LocalDynamodb docker container
- Fuzz testing with complicated protobuf message
- Generate methods for just generating "Key" attribute maps
- Generate methods for creating put/get items for transactions
- Generate methods for handling dynamodb stream Lambda events, use: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#FromDynamoDBStreamsMap
- Allow nested messages (oneof) to be stored as protojson/protobinary instead of nested maps
- Similar to: https://github.com/GoogleCloudPlatform/protoc-gen-bq-schema
- Don't generate central unmarshal/marshal method more than once per package. Simply check file existence?

## Expression generation utility

Provide functionality similar to a query builder to support building the differente

- ConditionExpression: PutItem,DeleteItem,UpdateItem (TransactWriteItem)
- UpdateExpression: UpdateItem
- ProjectionExpression: GetItem, BatchGetItem (TransactGetItems), Query
- FilterExpression: Query
- KeyConditionExpression: Query

Integrate with: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression

- Idea: Generate methods to return "NameBuilder" instances, include method method that returns "NameBuilder" for partition/sort key. So can return PartitionKeyName().AttributeExists()

## Key generation and utility

We would like to generate some way to generate Dynamo key helper methods. Carefull key construction in DynamoDB
is important because it directly determines how data can be queried effiently. It is unlikely that
we can capture this logic in Protobuf. Some ideas

- Allow certain fields to be involved in the key generation, then generate a method/function that takes
  a lambda and should output the key(s)
  - It will also error if the fields are not set before this is called
- PK/SK are often values stringed together for example with "#"
- Parts of the PK/SK can be constants, maybe allow enums to be used?
- Parts of the PK/SK can be constant given the message type
- Would be nice if these pk/sk constructions are flagged when changed in a backwards incompatible way
- Could generate functions that create key attribute maps

What should the helping do for the various methods

- GetItem, UpdateItem, DeleteItem: turn basic Go type (string, int), passed in a request parameters into a
  attribute map with just the key values
- PutItem: should produce a full item, including keys into a dynamodb attribute map
- QueryItem: should produce a exact partition key, and variety of expressions on the sort key

## Minimal Viable Backlog

- [ ] MUST generate methods that return PartitionKey (name/value), and SortKey (name/value)
- [ ] MUST deploy a buf module so users can easily include options

## Documentation backlog

- [ ] Write the rules for what kind of fields can be set as a PK, or SK
- [ ] Write what well-known types encode to what

## Feature Backlog

- [ ] SHOULD add field option to support (un)marshalling StringSets, NumberSets, ByteSets etc
- [ ] SHOULD allow skipping certain fields for all dynamodb marshalling/unmarshalling: ignore option
- [ ] SHOULD support encoding compex types (messages, maps, strucpb, oneof values as json AND/OR binary protobuf)
- [ ] COULD generate query building structure for FilterExpressions/KeyConditionExpressions etc
- [ ] COULD add errors to the "MarshalDynamoKey" method if range/sort key is empty string, or empty bytes (0 number is fine?)
- [ ] COULD generate var/consts that return the attribute name for a member, so it can be used in DynamoExpressions
      , so instead of `attribute_not_exists(pk)` it can be `attribute_not_exists("+modelv1.ProfileSortKey+")` - Should generates methods that returns the SortKey() and PartitionKey() so interfaces can be defined
      for helper methods
- [ ] COULD (if we have annotations on what are the keys), to implement a GetItem implementation that fetches
      the rest of the item from the database.
- [ ] COULD add option to "skip unsupported" instead of error (but why not just "ignore" the field?)
- [ ] COULD make it configurable on how to handle nil/empty fields like stdlib json package
- [ ] COULD support pk/sk options for fields that are message types, as long as the message has textencoding interface of some sort
- [ ] COULD allow customizing the encoder/decoder options
  - make sure that logic applies to both marshalling, and unmarshalling (new test case)
- [ ] COULD improve usability of FieldMask encoding, instead of slice of strings of the field names in
      proto definition, could/should be the dynamodb attribute names. But this probably means implement another version of the fieldmaskpb.New() function. But https://pkg.go.dev/google.golang.org/protobuf/types/known/fieldmaskpb#Intersect states that "field.number" paths are also valid

## Hardening Backlog

- [ ] SHOULD unit test the "ddb" shared package
- [ ] SHOULD test boolean key maps
- [ ] SHOULD fix go vet checks failure
- [ ] SHOULD Add test that errors when unsupported map type is used
- [ ] SHOULD test with coverage test as described here: https://go.dev/blog/integration-test-coverage
- [ ] SHOULD Fuzz the "FieldPresence" message as well, but this might require revamping the fuzzing setup because
      of interface types in the generated types
- [ ] SHOULD test that messages from external packages that DO implement the MarshalDynamoItem can be used
      in fields without problem
- [ ] SHOULD turn panics into errors (or add catch mechanism)
- [ ] COULD fuzz using official go toolchain fuzzing

## Done Backlog

- [x] SHOULD add method that marshals just the keys (if any keys are configured), fail if more keys, fail if no data in keys
- [x] SHOULD add code generation that adds methods to return the PartitionKey and SortKey from a message
- [x] SHOULD run test with parralel and -race enabled
- [x] MUST add file header that states that the file is generated
- [x] SHOULD test support of wrapper types (what about optional field with wrapper types?)
- [x] SHOULD add a test that passes in empty (or half empty) attribute maps into unmarshal and unmarshalled struct
      to match what would be marshalled from from json.
- [x] SHOULD match output field presence to that of json encoding
- [x] SHOULD test with OneOf field
- [x] SHOULD Make sure generated error handling prints the field name and a more descriptive error
- [x] SHOULD allow customizing the dynamodb attribute name (from default "field number")
- [x] SHOULD error when marshalling map of messages and the key is an empty string (not allowed)
- [x] SHOULD get a good grip on "Nullability" of fields, the documentation for "WrapperTypes" suggests that: "Wrappers use the same representation in JSON as the wrapped primitive type, except that null is allowed and preserved during data conversion and transfer.". Which suggests that
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
