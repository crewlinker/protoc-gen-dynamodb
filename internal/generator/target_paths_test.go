package generator_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/crewlinker/protoc-gen-dynamodb/ddb/ddbpath"
	messagev1ddbpath "github.com/crewlinker/protoc-gen-dynamodb/proto/example/message/v1/ddbpath"
)

// test the building of paths
var _ = DescribeTable("path building", func(s expression.NameBuilder, expCondition string, expNames map[string]string) {

	// check that expression builder accepts the paths
	expr, err := expression.NewBuilder().WithCondition(s.AttributeExists()).Build()
	Expect(err).ToNot(HaveOccurred())

	Expect(*expr.Condition()).To(Equal(fmt.Sprintf(`attribute_exists (%s)`, expCondition)))
	Expect(expr.Names()).To(Equal(expNames))

},
	Entry("basic type field",
		messagev1ddbpath.Kitchen().Brand(),
		"#0",
		map[string]string{"#0": "1"}),
	Entry("nested field",
		messagev1ddbpath.Kitchen().ExtraKitchen().Brand(),
		"#0.#1",
		map[string]string{"#0": "16", "#1": "1"}),
	Entry("extra nested field",
		messagev1ddbpath.Kitchen().ExtraKitchen().ExtraKitchen().Brand(),
		"#0.#0.#1",
		map[string]string{"#0": "16", "#1": "1"}),

	Entry("basic type list",
		messagev1ddbpath.Kitchen().OtherBrands().Index(10),
		"#0[10]",
		map[string]string{"#0": "20"}),
	Entry("message list",
		messagev1ddbpath.Kitchen().ApplianceEngines().Index(3).Brand(),
		"#0[3].#1",
		map[string]string{"#0": "19", "#1": "1"}),

	Entry("basic type map",
		messagev1ddbpath.Kitchen().Calendar().Key("bar"),
		"#0.#1",
		map[string]string{"#0": "14", "#1": "bar"}),
	Entry("message map",
		messagev1ddbpath.Kitchen().Furniture().Key("dar").Brand(),
		"#0.#1.#2",
		map[string]string{"#0": "13", "#1": "dar", "#2": "1"}),
)

// test path validation with generated logic
var _ = DescribeTable("path validation", func(nb interface {
	AppendName(field expression.NameBuilder) expression.NameBuilder
}, paths []string, expErr string) {
	err := ddbpath.Validate(nb, paths...)
	if expErr == "" {
		Expect(err).To(BeNil())
	} else {
		Expect(err.Error()).To(MatchRegexp(expErr))
	}
},
	Entry("should validate named attr", messagev1ddbpath.FieldPresencePath{}, []string{"msg.1"}, ``),
	Entry("omitted field should be invalid", messagev1ddbpath.IgnoredPath{}, []string{"1"}, ` non-existing field '1' on: messagev1ddbpath.IgnoredPath`),
)
