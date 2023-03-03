package generator

import (
	"fmt"
	"io"
	"text/template"

	"go.uber.org/zap"
	"google.golang.org/protobuf/compiler/protogen"
)

// TargetDesc describes generation of a single protobuf file
type TargetDesc struct {
	*protogen.File
	Items []*ItemDesc
}

// ItemFieldDesc describes the field of a dynamodb item
type ItemFieldDesc struct {
	DynamoName string            // name in Dynamo table
	GoName     string            // name in Go code
	Message    *protogen.Message // if the fields type is another message
}

// ItemDesc describes the DynamoDB item and keys for generation
type ItemDesc struct {
	GoIdent protogen.GoIdent // type identifier in Go code
	Fields  []*ItemFieldDesc // protobuf defined fields
}

// Target facilitates generation from a single protobuf file
type Target struct {
	src  *protogen.File
	logs *zap.Logger
	tmpl *template.Template
}

// generateItemField generates an item descriptor from a proto message
func (tg *Target) generateItemField(pgf *protogen.Field) (*ItemFieldDesc, error) {
	desc := &ItemFieldDesc{
		DynamoName: fmt.Sprintf("%d", pgf.Desc.Number()),
		GoName:     pgf.GoName,
		Message:    pgf.Message,
	}

	return desc, nil
}

// generateItem generates an item descriptor from a proto message
func (tg *Target) generateItem(pgm *protogen.Message) (*ItemDesc, error) {
	desc := &ItemDesc{GoIdent: pgm.GoIdent}
	for _, field := range pgm.Fields {
		itemf, err := tg.generateItemField(field)
		if err != nil {
			return nil, fmt.Errorf("failed to generate item field: %w", err)
		}
		desc.Fields = append(desc.Fields, itemf)
	}

	return desc, nil
}

// Generate peforms the actual code generation
func (tg *Target) Generate(w io.Writer) error {
	tg.logs.Info("generating DynamoDB helpers") // @TODO log src and destination
	desc := &TargetDesc{File: tg.src}
	for _, msg := range tg.src.Messages {
		item, err := tg.generateItem(msg)
		if err != nil {
			return fmt.Errorf("failed to generate item from message: %w", err)
		}
		desc.Items = append(desc.Items, item)
	}

	if err := tg.tmpl.ExecuteTemplate(w, "dynamo.gotmpl", desc); err != nil {
		return fmt.Errorf("failed to generate resolving code: %w", err)
	}

	return nil
}
