package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
			pathfp := filepath.Join(
				filepath.Dir(pf.GeneratedFilenamePrefix),
				fmt.Sprintf("%sddb", string(pf.GoPackageName)),
				fmt.Sprintf("%s.go", filepath.Base(pf.GeneratedFilenamePrefix)),
			)
			pathf := gp.NewGeneratedFile(pathfp, pf.GoImportPath)

			// pathf := gp.NewGeneratedFile(
			// 	filepath.Join(
			// 		filepath.Dir(pf.GeneratedFilenamePrefix),
			// 		fmt.Sprintf("%sattr", string(pf.GoPackageName)),
			// 		fmt.Sprintf("%s.go", filepath.Base(pf.GeneratedFilenamePrefix)),
			// 	),
			// 	pf.GoImportPath,
			// )

			fmt.Fprintf(os.Stderr, "%v %v\n", pathfp, pf.GoImportPath)

			// attrf := gp.NewGeneratedFile(fmt.Sprintf("%s/.ddb.go", pf.GeneratedFilenamePrefix))

			tg := gen.CreateTarget(pf)
			if err := tg.GenerateMessageLogic(ddbf); err != nil {
				return fmt.Errorf("failed to generate message logic for '%s': %w", *pf.Proto.Name, err)
			}

			if err := tg.GeneratePathBuilding(pathf); err != nil {
				return fmt.Errorf("failed to generate path building code for '%s': %w", *pf.Proto.Name, err)
			}

		}

		return nil
	})
}
