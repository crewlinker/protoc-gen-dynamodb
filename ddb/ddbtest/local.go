// Package ddbtest provides helper code for testing DynamoDB code
package ddbtest

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// NewLocalClient returns a DynamoDB client that is configured to connect to a locally hosted DynamoDB. An
// optional url to the locally hosted dynamodb client can be provided. By default it is set to localhost:8000
func NewLocalClient(epurl ...string) (*dynamodb.Client, error) {
	ep := aws.Endpoint{URL: "http://localhost:8000", SigningRegion: "localhost"}
	if len(epurl) > 0 {
		ep.URL = epurl[0]
	}

	// for local testing we won't be doing io, so context will probably not be used but we
	// configure it just to be sure
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, opts ...any) (aws.Endpoint, error) { return ep, nil })),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider("local-test-key-id", "local-test-key-secret", "")),
		config.WithRegion("localhost"))
	if err != nil {
		return nil, fmt.Errorf("failed to load default config: %w", err)
	}

	return dynamodb.NewFromConfig(cfg), nil
}
