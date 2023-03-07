package generator

import (
	"fmt"
	"io"
	"strings"

	//lint:ignore ST1001 we want expressive meta code
	. "github.com/dave/jennifer/jen"
	"go.uber.org/zap"
	"google.golang.org/protobuf/compiler/protogen"
)

const (
	attributevalues = "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	dynamodbtypes   = "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Target facilitates generation from a single protobuf file
type Target struct {
	src    *protogen.File
	logs   *zap.Logger
	idents struct {
		marshal   string
		unmarshal string
	}
}

// Generate peforms the actual code generation
func (tg *Target) Generate(w io.Writer) error {
	f := NewFile(string(tg.src.GoPackageName))

	tg.idents.marshal, tg.idents.unmarshal =
		fmt.Sprintf("%s_marshal_dynamo_item", strings.ToLower(tg.src.GoDescriptorIdent.GoName)),
		fmt.Sprintf("%s_unmarshal_dynamo_item", strings.ToLower(tg.src.GoDescriptorIdent.GoName))

		// generate a single marshal function for messages. This way we can handle externally included messages
	// and locally generated message in the same way.
	if err := tg.genCentralMarshal(f); err != nil {
		return fmt.Errorf("failed to generate central message marshal: %w", err)
	}

	// generate a single unmarshal function per proto file so we can unmarshal external messages and
	// messages local to the package in the same way.
	if err := tg.genCentralUnmarshal(f); err != nil {
		return fmt.Errorf("failed to generate central message unmarshal: %w", err)
	}

	// generate per message marshal/unmarshal code
	for _, m := range tg.src.Messages {
		if err := tg.genMessageMarshal(f, m); err != nil {
			return fmt.Errorf("failed to generate marshal: %w", err)
		}
		if err := tg.genMessageUnmarshal(f, m); err != nil {
			return fmt.Errorf("failed to generate unmarshal: %w", err)
		}
	}

	return f.Render(w)
}
