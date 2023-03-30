package main

import (
	"flag"
	"fmt"
	"path"
	"path/filepath"
	"strconv"

	"github.com/crewlinker/protoc-gen-dynamodb/internal/generator"
	"go.uber.org/zap"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	flag.Parse()
	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
	}.Run(func(gp *protogen.Plugin) error {
		gp.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		logs, err := zap.NewDevelopment()
		if err != nil {
			return fmt.Errorf("failed to setup logging: %w", err)
		}

		opts := generator.Config{}

		gen, err := generator.NewGenerator(logs, opts)
		if err != nil {
			return fmt.Errorf("failed to initialize generator: %w", err)
		}

		for _, name := range gp.Request.FileToGenerate {
			pf := gp.FilesByPath[name]
			if len(pf.Messages) < 1 {
				logs.Info("no messages in file, skipping", zap.String("file", name))
				continue // without services there is nothing to build a graphql schema for
			}

			logs.Info("found file with messages", zap.Int("num_messages", len(pf.Messages)))
			ddbf := gp.NewGeneratedFile(fmt.Sprintf("%s.ddb.go", pf.GeneratedFilenamePrefix), pf.GoImportPath)

			// generated file for typed document path in a sub directory for more expressiveness
			ddbPkgName := fmt.Sprintf("%sddb", string(pf.GoPackageName))
			ddbFp := filepath.Join(
				filepath.Dir(pf.GeneratedFilenamePrefix),
				ddbPkgName,
				fmt.Sprintf("%s.go", filepath.Base(pf.GeneratedFilenamePrefix)),
			)

			ddbImpName, _ := strconv.Unquote(pf.GoImportPath.String())
			ddbImpName = path.Join(ddbImpName, ddbPkgName)

			tg := gen.CreateTarget(pf, ddbImpName)
			if err := tg.GenerateMessageLogic(ddbf); err != nil {
				return fmt.Errorf("failed to generate message logic for '%s': %w", *pf.Proto.Name, err)
			}

			if err := tg.GeneratePathBuilding(gp.NewGeneratedFile(ddbFp, pf.GoImportPath)); err != nil {
				return fmt.Errorf("failed to generate path building code for '%s': %w", *pf.Proto.Name, err)
			}

		}

		return nil
	})
}
