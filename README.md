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

- [ ] SHOULD Add test that errors when unsupported map type is used
- [ ] SHOULD Make sure generated error handling prints the field name and a more descriptive error
- [ ] SHOULD allow customizing the encoder/decoder options
- [ ] MUST add file header that states that the file is generated
- [ ] SHOULD make it configurable on how to handle nil/empty fields like stdlib json package
- [ ] SHOULD test that messages from external packages that DO implement the MarshalDynamoItem can be used in fields without problem
- [ ] SHOULD turn panics into errors (or add catch mechanism)
- [ ] SHOULD add support of StringSets, NumberSets, ByteSets etc
- [ ] SHOULD allow skipping certain fields for all dynamodb marshalling/unmarshalling
- [x] SHOULD error when marshalling map of messages and the key is an empty string (not allowed)
