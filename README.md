# protoc-gen-dynamodb

Use Protobuf to define DynamoDB item encoding using Go (golang).

## features

- Uses sdk v2
- Unit and e2e testing
- Generate table definitions
- Use protobuf field numbers
- use official 'attributevalue' with ability to customize its behaviour

## Ideas

- Define the sk/pk on the message, and set the pk/sk member of the resulting item, error when sk/pk is not set?
- Consider using the field position numbers instead of names (since they are supposed to be stable)
- Generate table definitions for use in AWS Cloudformation/CDK
- E2E testing with LocalDynamodb docker container
- Fuzz testing with complicated protobuf message
- Generate methods for just generating "Key" attribute maps
- Generate methods for creating put/get items for transactions
- Generate methods for handling dynamodb stream Lambda events, use: https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue#FromDynamoDBStreamsMap
