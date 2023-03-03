package generator

import (
	"embed"
	"fmt"
	"text/template"

	"go.uber.org/zap"
	"google.golang.org/protobuf/compiler/protogen"
)

//go:embed *.gotmpl
var tmplfs embed.FS

// Options for configuring the generator
type Options struct{}

// Generator generates DynamoDB helper functions
type Generator struct {
	opts Options
	logs *zap.Logger
	tmpl *template.Template
}

// New inits the generator
func New(logs *zap.Logger, opts Options) (g *Generator, err error) {
	g = &Generator{
		logs: logs.Named("generator"),
		opts: opts,
	}

	g.tmpl, err = template.ParseFS(tmplfs, "*.gotmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse filesystem: %w", err)
	}

	return g, nil
}

// CreateTarget inits a target for a generator
func (g Generator) CreateTarget(pf *protogen.File) *Target {
	return &Target{
		src:  pf,
		tmpl: g.tmpl,
		logs: g.logs.Named(fmt.Sprintf("target[%s]", *pf.Proto.Name)),
	}
}
