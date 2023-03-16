// Package generator implements the code generator
package generator

import (
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/protobuf/compiler/protogen"
)

// Config for configuring the generator
type Config struct{}

// Generator generates DynamoDB helper functions
type Generator struct {
	cfg  Config
	logs *zap.Logger
}

// NewGenerator inits the generator
func NewGenerator(logs *zap.Logger, opts Config) (g *Generator, err error) {
	g = &Generator{
		logs: logs.Named("generator"),
		cfg:  opts,
	}

	return g, nil
}

// CreateTarget inits a target for a generator
func (g Generator) CreateTarget(pf *protogen.File) *Target {
	return &Target{
		src:  pf,
		logs: g.logs.Named(fmt.Sprintf("target[%s]", *pf.Proto.Name)),
	}
}
