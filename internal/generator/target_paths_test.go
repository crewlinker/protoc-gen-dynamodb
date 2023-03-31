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

	// check that the resulting expression(names) are as expected
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

	// well-known: anypb
	Entry("any field",
		messagev1ddbpath.Kitchen().SomeAny().TypeURL(),
		"#0.#1",
		map[string]string{"#0": "21", "#1": "1"}),
	Entry("any field",
		messagev1ddbpath.Kitchen().SomeAny().Value(),
		"#0.#1",
		map[string]string{"#0": "21", "#1": "2"}),
	Entry("list of anypb",
		messagev1ddbpath.Kitchen().RepeatedAny().Index(13).TypeURL(),
		"#0[13].#1",
		map[string]string{"#0": "31", "#1": "1"}),
	Entry("map of anypb",
		messagev1ddbpath.Kitchen().MappedAny().Key("koo").TypeURL(),
		"#0.#1.#2",
		map[string]string{"#0": "32", "#1": "koo", "#2": "1"}),

	// well-known paths
	// Entry("durationpb", messagev1ddbpath.Kitchen().Timer())
	// case *durationpb.Duration, *timestamppb.Timestamp: AttributeValueMemberS
	// case *fieldmaskpb.FieldMask: AttributeValueMemberSS
	// case *structpb.Value: <anything>
	// case *wrapperpb.<Kind>Value: just the basic path
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

	// well-known: anypb
	Entry("anypb", messagev1ddbpath.Kitchen(), []string{"21.1"}, ``),
	Entry("anypb", messagev1ddbpath.Kitchen(), []string{"21.2"}, ``),
	Entry("anypb", messagev1ddbpath.Kitchen(), []string{"31[999].1"}, ``),
	Entry("anypb", messagev1ddbpath.Kitchen(), []string{"32.foo.1"}, ``),
	// @TODO test path into 21.1
)
