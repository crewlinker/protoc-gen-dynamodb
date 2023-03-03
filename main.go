package main

import (
	"flag"
	"fmt"

	"github.com/crewlinker/protoc-gen-dynamodb/internal/generator"
	"go.uber.org/zap"

	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	flag.Parse()
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gp *protogen.Plugin) error {
		logs, err := zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("failed to setup logging: %w", err)
		}

		opts := generator.Options{}

		gen, err := generator.New(logs, opts)
		if err != nil {
			return fmt.Errorf("failed to initialize generator: %w", err)
		}

		for _, name := range gp.Request.FileToGenerate {
			pf := gp.FilesByPath[name]
			if len(pf.Messages) < 1 {
				logs.Info("no messages in file, skipping", zap.String("file", name))
				continue // without services there is nothing to build a graphql schema for
			}

			logs.Info("found file with services", zap.Int("num_services", len(pf.Services)))
			ddbf := gp.NewGeneratedFile(fmt.Sprintf("%s.ddb.go", pf.GeneratedFilenamePrefix), pf.GoImportPath)

			if err := gen.CreateTarget(pf).Generate(ddbf); err != nil {
				return fmt.Errorf("failed to generate for '%s': %w", *pf.Proto.Name, err)
			}
		}

		return nil
	})
}
