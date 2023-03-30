// Package generator implements the code generator
package generator

import (
	"fmt"
	"path"
	"runtime/debug"

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
	tg := &Target{
		src:  pf,
		logs: g.logs.Named(fmt.Sprintf("target[%s]", *pf.Proto.Name)),
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		panic("failed to read build info: binary not build with modules support")
	}

	// tg idents provides various identifiers
	tg.idents.ddb = path.Join(bi.Path, "ddb")
	tg.idents.ddbpath = path.Join(bi.Path, "ddb", "ddbpath")
	tg.idents.ddbv1 = path.Join(bi.Path, "proto", "ddb", "v1")
	return tg
}
